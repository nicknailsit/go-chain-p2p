package p2p

import (
	"bufio"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"sync"
)

type StreamHandlerInterface interface {
	Reader(rw *bufio.ReadWriter)
	Writer(rw *bufio.ReadWriter)
	Handler(stream network.Stream)
	Endpoint() string
}

type TypeHandler struct {
	sync.Mutex
	StreamType string
	DataWriter bufio.Writer
	DataReader bufio.Reader
	endpoint string
}



var SyncTypes  = []string{"SYNCFULL", "SYNCFROM", "SYNCONE", "SYNCLITE"}
var ContentTypes = []string{"SIMAGES","SVIDEOS","SPOSTS","SMESSAGES","SIMAGE","SVIDEO","SPOST","SMESSAGE"}

func (t TypeHandler) Reader(rw *bufio.ReadWriter) {


}

func (t TypeHandler) Writer(rw *bufio.ReadWriter) {

}

func (t TypeHandler) Endpoint() string {
	return t.endpoint
}

func (t TypeHandler) Handler(stream network.Stream)  {

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	t.Reader(rw)
	t.Writer(rw)

}

func (t TypeHandler) create(streamType string,  endpoint string) *TypeHandler {

	T := &TypeHandler{
		StreamType: streamType,
		endpoint:endpoint,
	}

	return T

}


type SyncChainHandler struct {
	th *TypeHandler
	ChainService *ChainService
	state map[string][]byte
	HeartBeat *ping.PingService
}

func (sc SyncChainHandler) FullChain(peer host.Host, stream network.Stream) *SyncChainHandler {

	T := sc.th.create("SYNCFULL", "/chain/sync/full")

	Syncer := &SyncChainHandler{
		th: T,
		ChainService: NewChainService(peer, uint8(0x77)),
		state: make(map[string][]byte),
		HeartBeat: ping.NewPingService(peer),

	}

	return Syncer


}


