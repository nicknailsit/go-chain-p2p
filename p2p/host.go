package p2p

import (
	"context"
	"fmt"
	"github.com/ipsn/go-ipfs/repo"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
)

func NewPeerHost(port int, repo *repo.Repo) (host.Host, error) {
	privKey = repo.PrivKey()

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)),
		libp2p.Identity(privKey),
	}

	Host, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	return Host, nil
}
