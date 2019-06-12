package swaggchain

import (
	"context"
	"crypto/sha256"

	"github.com/gogo/protobuf/proto"
	"github.com/ipfs/go-cid"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	fs "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multihash"

	"github.com/op/go-logging"
	"io"
	"swaggp2p/pb"
	"swaggp2p/repo"
	"swaggp2p/services"
	"sync"
	"time"
)

var (
	log = logging.MustGetLogger("cmd");
	Topic cid.Cid
)

const (
	ReSubscribeInterval = time.Hour
	ReconnectInterval = time.Minute
	MinConnectedSubscribers = 2

)

func init() {
	h := sha256.Sum256([]byte("floodsub:Chain"))
	enc, err := multihash.Encode(h[:], multihash.SHA2_256)
	if err != nil {
		panic(err)
	}
	mh, err := multihash.Cast(enc)
	if err != nil {
		panic(err)
	}
	Topic = cid.NewCidV1(cid.Raw, mh)
}

type addPeer struct {
	peerID peer.ID
}

type removePeer struct {
	peerID peer.ID
}

type newBlockchain struct {
	serializedMessage []byte
	mine bool
}

type getBlockchain struct {
	serializedMessage []byte
	mine bool
}

type getHeaders struct {
	serializedMessage []byte
	mine bool
}

type getCurrentTime struct {
	serializedMessage []byte
	mine bool
}

type login struct {
	serializedMessage []byte
	mine bool
}

type logout struct {
	serializedMessage []byte
	mine bool
}

type Chain struct {
	chainID string
	blocks Blocks
	serializedData []byte
}

type Block struct {
	*pb.Block
}

type Blocks []*pb.Block;

type BlockHeaders []*pb.BlockHeader;

type SwaggNode struct {
	repo *repo.Repo
	peerHost host.Host
	routing *dht.IpfsDHT
	floodsub *fs.PubSub
	msgChan chan interface{}
	connectedSubs map[peer.ID]bool
	orderBook interface{}

	wireService *services.WireService
	timeService *services.UTCTimeService

}

func NewSwaggNode(repo *repo.Repo, peerHost host.Host, routing *dht.IpfsDHT, floodsub *fs.PubSub) *SwaggNode {
	return &SwaggNode{
		repo: repo,
		peerHost: peerHost,
		routing: routing,
		floodsub: floodsub,
		msgChan: make(chan interface{}),
		connectedSubs: make(map[peer.ID]bool),
		//orderBook: interface{},

	}
}

func (n *SwaggNode) MsgChan() chan interface{} {
	return n.msgChan
}

func (n *SwaggNode) SetWireService(ws *services.WireService) {
	n.wireService = ws
}
/*
func (n *SwaggNode) SetChainService(cs *services.ChainService) {
	n.blockchain = cs
}*/

func (n *SwaggNode) SetUTCTimeService(utc *services.UTCTimeService) {
	n.timeService = utc
}


func (n *SwaggNode) StartOnlineServices() {
	go n.subscribeTopic()
	go n.connectToSubscribers()
	go n.messageHandler()
}

func (n *SwaggNode) messageHandler() {

	for {
		select {
		case m := <-n.msgChan:
			switch msg := m.(type) {
			case addPeer:
				n.connectedSubs[msg.peerID] = true
			case removePeer:
				if _, ok := n.connectedSubs[msg.peerID]; ok {
					log.Infof("Lost subscriber peer %s", msg.peerID.Pretty())
					delete(n.connectedSubs, msg.peerID)
				}
			case newBlockchain:


			case getBlockchain:


			case getCurrentTime:


			}
		}
	}

}

func (n *SwaggNode) subscribeTopic() {
	go n.setSelfAsSubscriber()

	sub, err := n.floodsub.Subscribe("swaggchain")
	if err != nil {
		log.Error(err)
		return
	}

	for {
		msg, err := sub.Next(context.Background())
		if err == io.EOF || err == context.Canceled {
			return
		} else if err != nil {
			log.Error(err)
			return
		}

		mpb := new(pb.Message)
		err = proto.Unmarshal(msg.Data, mpb)
		if err != nil {
			log.Error(err)
			return
		}

		switch mpb.MessageType {

		case pb.Message_SendFullChain:
			n.msgChan <- getBlockchain{serializedMessage: mpb.Payload}

		case pb.Message_AuthLogin:
			n.msgChan <- login{serializedMessage: mpb.Payload}

		case pb.Message_AuthLogout:
			n.msgChan <- logout{serializedMessage: mpb.Payload}

		}
	}
}

func (n *SwaggNode) setSelfAsSubscriber() {
	subscribe := func() {
		err := n.routing.Provide(context.Background(), Topic, true)
		if err != nil {
			log.Error(err)
		}
	}
	subscribe()

	ticker := time.NewTicker(ReSubscribeInterval)
	for range ticker.C {
		subscribe()
	}
}

func (n *SwaggNode) connectToSubscribers() {
	n.connectionRound()
	ticker := time.NewTicker(ReconnectInterval)
	for range ticker.C {
		for peer := range n.connectedSubs {
			conns := n.peerHost.Network().ConnsToPeer(peer)
			if len(conns) == 0 {
				n.msgChan <- removePeer{peer}
			}
		}
		if len(n.connectedSubs) < MinConnectedSubscribers {
			n.connectionRound()
		}
	}
}


func (n *SwaggNode) connectionRound() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	provs := n.routing.FindProvidersAsync(ctx, Topic, 10)
	wg := &sync.WaitGroup{}

	for p := range provs {
		wg.Add(1)
		go func(pi peer.AddrInfo) {

			defer wg.Done()
			ctx, cancel := context.WithTimeout(ctx, time.Second*10)

			defer cancel()

			err := n.peerHost.Connect(ctx, pi)
			if err != nil {
				log.Debug("pubsub discover: ", err)
				return
			}

			log.Debug("connected to pubsub peer:", pi.ID)
			n.msgChan <- addPeer{pi.ID}

			if n.wireService != nil {
				m := &pb.Message{
					MessageType: pb.Message_GetChainInfo,
				}
				n.wireService.SendMessage(pi.ID, m)
			}


		}(p)
	}
	wg.Wait()
}

