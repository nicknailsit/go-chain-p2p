package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multiaddr"
)

func NewHost() {

	ctx:= context.Background()




	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		panic(err)
	}

	cfg := defaultConfig()

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cfg.listenHost, cfg.listenPort))

	host, err := libp2p.New(ctx,libp2p.Identity(priv), libp2p.ListenAddrs(sourceMultiAddr))
	if err != nil {
		panic(err)
	}


	//sync blockchain #1

	host.SetStreamHandler(protocol.ID(SyncProtocolID), handleSyncStream)

	fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s\n", cfg.listenHost, cfg.listenPort, host.ID().Pretty())

	peerChan := initMDNS(ctx, host, cfg.RendezvousString)
	peer := <-peerChan
	fmt.Println("Found peer:", peer, ", connecting")

	if err := host.Connect(ctx, peer); err != nil {
		fmt.Println("Connection failed:", err)
	}

	stream, err := host.NewStream(ctx, peer.ID, protocol.ID(SyncProtocolID))

	if err != nil {
		fmt.Println("fail opening stream")
	} else {

		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
		go writeData(rw)
		go readData(rw)
		fmt.Println("Connected to:", peer)

	}


	select{}


}

func defaultConfig() *config {

	cfg := &config{
		listenHost:"127.0.0.1",
		listenPort:"9000",
		RendezvousString:"swaggchain",
		ProtocolID: nil,

	}

return cfg

}