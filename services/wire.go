package services

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/op/go-logging"
	"io"
	"swaggp2p/chain/pb"
	ggio "github.com/gogo/protobuf/io"
	ctxio "github.com/jbenet/go-context/io"
	inet "github.com/libp2p/go-libp2p-core/network"
	"time"
)

const WireProtocol protocol.ID = "/Wire/1.0.0"

var logservice = logging.MustGetLogger("service")

type WireService struct {
	msgChan WireChannel
	chainSync *ChainService
	timeSync *UTCTimeService
	orderBook *interface{}
	peerHost host.Host
}

type WireChannel chan interface{}

func NewWireService(msgChan chan interface{}, wireChannel WireChannel, peerHost host.Host) *WireService {

	ws := &WireService {
		msgChan:msgChan,
	}
	ws.peerHost.SetStreamHandler(WireProtocol, ws.handleNewStream)
	return ws


}

func (ws *WireService) SendMessage(peer peer.ID, pmes *pb.Message) error {

	s, err := ws.peerHost.NewStream(context.Background(), peer, WireProtocol)
	if err != nil {
		return err
	}
	writer := ggio.NewDelimitedWriter(s)
	return writer.WriteMsg(pmes)

}

func (ws *WireService) SendRequest(peer peer.ID, pmes *pb.Message) (*pb.Message, error) {
	s, err := ws.peerHost.NewStream(context.Background(), peer, WireProtocol)
	if err != nil {
		return nil, err
	}
	writer := ggio.NewDelimitedWriter(s)
	err = writer.WriteMsg(pmes)
	if err != nil {
		return nil, err
	}

	cr := ctxio.NewReader(context.Background(), s)
	r := ggio.NewDelimitedReader(cr, inet.MessageSizeMax)
	rmes := new(pb.Message)
	if err := r.ReadMsg(rmes); err != nil {
		s.Reset()
		if err == io.EOF {
			log.Debugf("Disconnected from peer %s", peer.Pretty())
		}
		return nil, err
	}

	return rmes, nil
}

func (ws *WireService) handleNewStream(s inet.Stream) {
	go ws.handleNewMessage(s)
}

func (ws *WireService) handleNewMessage(s inet.Stream) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	timer := time.NewTimer(time.Second *TTL)
	defer timer.Stop()
	go func(){
		select {
		case <-timer.C:
		case <-ctx.Done():
		}
		s.Close()
	}()

	for {
		cr := ctxio.NewReader(context.Background(), s)
		r := ggio.NewDelimitedReader(cr, inet.MessageSizeMax)
		mPeer := s.Conn().RemotePeer()
		pmes := new(pb.Message)
		if err := r.ReadMsg(pmes); err != nil {
			s.Reset()
			if err == io.EOF {
				log.Debugf("Disconnected from peer %s", mPeer.Pretty())

			}
			return
		}
		handler := ws.handlerForMsgType(pmes.MessageType)
		rmes, err := handler(mPeer, pmes)
		if err != nil {
			log.Error(err)
			return
		}
		if rmes != nil {
			writer := ggio.NewDelimitedWriter(s)
			err = writer.WriteMsg(rmes)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (ws *WireService) handlerForMsgType(t pb.Message_MessageType) func(peer.ID, *pb.Message) (*pb.Message, error) {
	switch t {
	case pb.Message_SendFullChain:
		return ws.handleSendFullChain

	case pb.Message_AuthLogin:
		return ws.handleAuthLogin

	case pb.Message_AuthLogout:
		return ws.handleAuthLogout

	case pb.Message_LimitOrder:
		return ws.handleLimitOrder

	case pb.Message_GetOrderBook:
		return ws.handleGetOrderBook

	default:
		return nil
	}
}

func (ws *WireService) handleSendFullChain(p peer.ID, msg *pb.Message) (*pb.Message, error) {

	return nil,nil

}

func (ws *WireService) handleAuthLogin(p peer.ID, msg *pb.Message) (*pb.Message, error) {

	return nil,nil

}

func (ws *WireService) handleAuthLogout(p peer.ID, msg *pb.Message) (*pb.Message, error) {

	return nil,nil

}

func (ws *WireService) handleLimitOrder(p peer.ID, msg *pb.Message) (*pb.Message, error) {

	return nil,nil

}

func (ws *WireService) handleGetOrderBook(p peer.ID, msg *pb.Message) (*pb.Message, error) {

	return nil,nil

}



