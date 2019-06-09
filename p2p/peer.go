package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	ma "github.com/multiformats/go-multiaddr"

)

type SwaggPeer struct {
	ctx context.Context
	P2PHost host.Host
	Identity libp2p.Option
	ListenAddr ma.Multiaddr
}

func (p *SwaggPeer) NewHost(host, port string) *SwaggPeer {

	r := rand.Reader
	ctx := context.Background()
	privKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)

	if err != nil {
		panic(err)
	}

	srcMa, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", host, port))

	sp := &SwaggPeer{}

	sp.ctx = ctx
	sp.ListenAddr = srcMa
	sp.Identity = libp2p.Identity(privKey)
	sp.P2PHost, err = libp2p.New(sp.ctx, libp2p.ListenAddrs(sp.ListenAddr), sp.Identity)
	if err != nil {
		panic(err)
	}

	return sp

}

func StartNewStream(hostname, port string, streamType string, streamHandler network.StreamHandler,  protocolID int) {


	S := &SwaggPeer{}
	sp := S.NewHost(hostname, port)
	sp.P2PHost.SetStreamHandler(protocol.ID(protocolID), streamHandler)

	peerChan := initMDNS(sp.ctx, hostname, streamType)
	peer := <-peerChan
	fmt.Println("Found peer:", peer, ", connecting")

	if err := sp.P2PHost.Connect(sp.ctx, peer); err != nil {
		fmt.Println("Connection failed: ", err)
	}

	stream, err := sp.P2PHost.NewStream(sp.ctx, peer.ID, protocol.ID(protocolID))

	if err != nil {
		fmt.Printf("%s stream failed %v ", streamType, err)
	} else {
		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
		go writeData(rw)
		go readData(rw)
		fmt.Println("Connected to: " ,peer)
	}

	select {}

}


