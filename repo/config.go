package repo

import (
	iaddr "github.com/ipfs/go-ipfs-addr"
	"github.com/libp2p/go-libp2p-core/peer"
)

var defaultBootstrapPeers = []string{}

func ParseBootstrapPeer(addr string) (iaddr.IPFSAddr, error) {

	ia, err := iaddr.ParseString(addr)
	if err != nil {
		return nil, err
	}

	return ia, err
}

func ParseBootstrapPeers(addrs []string) ([]peer.AddrInfo, error) {

	var peers []peer.AddrInfo
	for _, addr := range addrs {
		ia, err := ParseBootstrapPeer(addr)
		if err != nil {
			return nil, err
		}
		pi, _ := peer.AddrInfoFromP2pAddr(ia.Multiaddr())
		peers = append(peers, *pi)
	}
	return peers, nil

}
