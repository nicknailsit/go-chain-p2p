package chain

import (
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peer "github.com/libp2p/go-libp2p-peer"
	fs "github.com/libp2p/go-libp2p-pubsub"
	"github.com/op/go-logging"
	"swaggp2p/repo"
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

/*func makeRandomNode(port int, done chan bool) *Node {
	priv, _, _ := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	listen, _ := ma.NewMultiAddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port))
	host, _ := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(listen),
		libp2p.Identity(priv),
		)

	return NewNode(host, done)
}

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
*/
type SwaggNode struct {
	repo *repo.Repo
	peerHost host.Host
	routing *dht.IpfsDHT
	floodsub *fs.PubSub
	msgChan chan interface{}
	connectedSubs map[peer.ID]bool
	orderBook interface{}
	wireService interface{}


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

func (n *SwaggNode) SetWireService() {

}
