package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	"time"
)

//libp2p MDNS Discoverer


type discoveryChannel struct {
	PeerChan chan peer.AddrInfo
}

func (n *discoveryChannel) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

//init mdns

func initMDNS(ctx context.Context, peerhost host.Host, rendezvous string) chan peer.AddrInfo {
	ser, err := discovery.NewMdnsService(ctx, peerhost, time.Minute, rendezvous)
	if err != nil {
		panic(err)
	}
	n := &discoveryChannel{}
	n.PeerChan = make(chan peer.AddrInfo)
	ser.RegisterNotifee(n)
	return n.PeerChan
}

