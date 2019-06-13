package p2p

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/protocol"
	kbucket "github.com/libp2p/go-libp2p-kbucket"
	"github.com/libp2p/go-libp2p-kbucket/keyspace"
	"github.com/op/go-logging"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	"swaggp2p/swaggchain"
	"sync"
	"time"
)

type Client interface {

	Addresses() []string

	ID() string

	PeerCount() int

	StreamCount() int

	Send(data []byte)

	Receive() (data []byte)

	Request(chksum [32]byte) ([]byte, error)

	SetBlockChainHandler(handler BlockChainHandler)

	SetChallengeHandler(handler ChallengeHandler)

	SetProofHandler(handler ProofHandler)

	SetVerificationHandler(handler VerificationHandler)

	SetTransactionHandler(handler TransactionHandler)

	SetSocialContentHandler(handler SocialContentHandler)

	SetPOWHandler(handler PowHandler)

	SetSearchHandler(handler SearchHandler)

	SetTorNetworkProxyHandler(handler TorNetworkProxyHandler)


}

type Notifee struct {
	PeerChan chan peer.AddrInfo
}

type client struct {

	cache *lru.Cache
	cacheLock *sync.Mutex
	chainRequests chan chainRequest
	blockRequests chan blockRequest
	challengeRequests chan challengeRequest
	commitmentRequests chan commitmentRequest
	config *ClientConfig
	context context.Context
	host host.Host
	id peer.ID
	key keyspace.Key
	logger *logging.Logger
	peerstore peerstore.Peerstore
	proofRequests chan proofRequest
	protocol protocol.ID
	receive chan []byte
	receiveBlocks chan swaggchain.Block
	sendBlocks chan swaggchain.Block
	sendTX chan swaggchain.TX
	receiveTX chan swaggchain.TX
	send chan []byte
	spammerCache *lru.Cache
	spammerCacheLock *sync.Mutex
	//not yet implemented
	streamstore interface{}

	serviceDiscovery discovery.Service
	serviceNotifee *Notifee

	unsetChainHandler func()
	unsetBlockHandler func()
	unsetChallengeHandler func()
	unsetProofHandler func()
	unsetTransactionHandler func()
	unsetSocialContentHandler func()
	unsetPOWHandler func()
	unsetSearchHandler func()
	unsetTorNetworkProxyHandler func()
	unsetServiceDiscoveryHandler func()
	unsetVerificationHandler func()
	unsetHandlerLock *sync.Mutex
	unsetCommitmentHandler func ()

	verificationRequests chan verificationRequest
	witnessCache *lru.Cache
	witnessCacheLock *sync.Mutex

	blocksCache *lru.Cache
	blocksCacheLock *sync.Mutex


}

type BlockchainHandler func(chksum [32]byte, response chan swaggchain.Blockchain)
type BlockHandler func(chksum [32]byte, response chan swaggchain.Block)

type chainRequest struct {
	chksum [32]byte
	response chan swaggchain.Blockchain
}

type blockRequest struct {
	chksum [32]byte
	response chan swaggchain.Block
}


func (client *client) ID() string {
	return client.id.Pretty()
}

func (client *client) PeerCount() int {
	return len(client.peerstore.Peers())
}

func (client *client) Send(data []byte) {
	client.send <- data
}

func (client *client) Receive() []byte {

	return <-client.receive

}

func (client *client) Request(checksum [32]byte) ([]byte, error) {
	return nil, errors.New("TODO: Implement request method.")
}

func (client *client) SetBlockChainHandler(handler BlockchainHandler) {

	notify := make(chan struct{})
	client.unsetHandlerLock.Lock()
	client.unsetChainHandler()
	client.unsetChainHandler = func() {
		close(notify)
	}

	client.unsetHandlerLock.Unlock()

	go func() {

		for {
			select {
			case <-notify:
				return
				case request := <-client.chainRequests:
					handler(request.chksum, request.response)
			}
		}

	}()

}

func (client *client) StreamCount() int {

	return 0
}

func (client *client) Addresses() string {
	return ""
}


func (cconfig *ClientConfig) New() (Client, func(), error) {
	return cconfig.create()
}

func (cconfig *ClientConfig) create() (*client, func(), error) {

	var err error
	client := &client{}
	client.config = &ClientConfig{}
	*client.config = *cconfig

	client.cache, err = lru.New(client.config.ArtifactCacheSize)
	if err != nil {
		return nil,nil,err
	}

	client.cacheLock = &sync.Mutex{}

	client.chainRequests = make(chan chainRequest, client.config.ArtifactQueueSize)
	client.blockRequests = make(chan blockRequest, client.config.ArtifactQueueSize)
	client.challengeRequests = make(chan challengeRequest, 1)
	client.commitmentRequests = make(chan commitmentRequest, 1)

	client.context = context.Background()

	seed := make([]byte, 32)


	priv, pub, err := crypto.GenerateEd25519Key(bytes.NewReader(seed))

	if err != nil {
		return nil,nil,err
	}

	client.id, err = peer.IDFromPublicKey(pub)

	if err != nil {
		return nil,nil,err
	}

	client.key = keyspace.XORKeySpace.Key(kbucket.ConvertPeerID(client.id))

	client.logger = logging.MustGetLogger("p2p")

	_ = client.peerstore.AddPrivKey(client.id, priv)
	_  = client.peerstore.AddPubKey(client.id, pub)

	client.proofRequests = make(chan proofRequest, 1)

	client.protocol = protocol.ID(
		fmt.Sprintf(
			"/%s/%s",
			client.config.Network,
			client.config.Version,
		),
	)

	client.send = make(chan []byte, client.config.ArtifactQueueSize)
	client.receive = make(chan []byte, client.config.ArtifactQueueSize)

	client.spammerCache, err = lru.New(client.config.SpammerCacheSize)
	if err != nil {
		return nil,nil,err
	}

	client.spammerCacheLock = &sync.Mutex{}

	client.serviceDiscovery, err = discovery.NewMdnsService(
		client.context, client.host, time.Second, "swaggchain")

	if err != nil {
		return nil,nil,err
	}

	n := &Notifee{}
	n.PeerChan = make(chan peer.AddrInfo)
	client.serviceNotifee = n
	client.serviceDiscovery.RegisterNotifee(n)

	client.unsetChainHandler = func() {}
	client.unsetBlockHandler = func() {}
	client.unsetCommitmentHandler = func() {}
	client.unsetChallengeHandler = func() {}
	client.unsetProofHandler = func() {}
	client.unsetVerificationHandler = func() {}

	client.verificationRequests = make(chan verificationRequest, 1)

	client.witnessCache, err = lru.New(client.config.WitnessCacheSize)
	if err != nil {
		return nil,nil,err
	}
	client.witnessCacheLock = &sync.Mutex{}

	shutdown, err = client.bootstrap()

	if err != nil {
		return nil, nil, err
	}

	return client, func() {
		client.unsetChainHandler()
		client.unsetCommitmentHandler()
		client.unsetChallengeHandler()
		client.unsetProofHandler()
		client.unsetVerificationHandler()
		client.unsetBlockHandler()
		shutdown()
	}, nil


}

func (n *Notifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}