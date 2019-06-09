package p2p

import (
	"context"
	"errors"
	"fmt"
	"github.com/ipsn/go-ipfs/thirdparty/math2"
	"github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	pstore "github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-kad-dht"
	"math/rand"
	"sync"
	"time"
)

var ErrNotEnoughBootstrapPeers = errors.New("not enough bootstrap peers")

type BootstrapConfig struct {
	MinPeerTreshold int
	Period time.Duration
	ConnectionTimeout time.Duration
	BootstrapPeers func() []core.Host
}

var DefaultBootstrapConfig = BootstrapConfig{
	MinPeerTreshold:2,
	Period: 30 * time.Second,
	ConnectionTimeout: (30 * time.Second) / 3,
}

func BootstrapConfigWithPeers(pis []core.Host) BootstrapConfig {
	cfg := DefaultBootstrapConfig
	cfg.BootstrapPeers = func() []core.Host{
		return pis
	}
	return cfg
	}

func Bootstrap(DHT dht.IpfsDHT, peerHost host.Host, cfg BootstrapConfig) error {

	periodic := func() {
		if err := bootstrapRound(context.Background(), peerHost, cfg); err != nil {
			log.Debugf("bootstrap error: %s", err)
		}
	}

	ticker := time.NewTicker(cfg.Period)

	go func() {

		for {
			select {
			case <-ticker.C:
				periodic()
			}
		}

	}()

	periodic()

	if err := DHT.BootstrapWithConfig(context.Background(), dht.DefaultBootstrapConfig); err != nil {
		return err
	}

	return nil

}

func bootstrapRound(ctx context.Context, host host.Host, cfg BootstrapConfig) error {

	ctx, cancel := context.WithTimeout(ctx, cfg.ConnectionTimeout)
	defer cancel()
	id := host.ID()

	peers := cfg.BootstrapPeers()

	connected := host.Network().Peers()
	if(len(connected) >= cfg.MinPeerTreshold) {
		log.Debugf("%s core bootstrap skipped --- connected to %d (> %d) nodes", id,  len(connected), cfg.MinPeerTreshold)
		return nil

	}

	numToDial := cfg.MinPeerTreshold - len(connected)

	var notConnected []core.Host

	for _, p := range peers {

			if host.Network().Connectedness(p.ID()) != network.Connected {
				notConnected = append(notConnected, p)
			}
	}

	if(len(notConnected) < 1) {
		log.Debugf("%s no more bootstrap peers to create %d connections", id, numToDial)
		return ErrNotEnoughBootstrapPeers
	}

	randSubset := randomSubsetOfPeers(notConnected, numToDial)
	log.Debugf("%s bootstrapping to %d nodes: %s", id, numToDial, randSubset)
	return bootstrapConnect(ctx, host, randSubset)


}

func bootstrapConnect(ctx context.Context, ph host.Host, peers []core.Host) error {

	if(len(peers) < 1) {
		return ErrNotEnoughBootstrapPeers
	}

	errs := make(chan error, len(peers))
	var wg sync.WaitGroup

	for _, p := range peers {

		wg.Add(1)
		go func(p core.Host){

			defer wg.Done()
			log.Debugf("%s bootstrapping to %s", ph.ID(), p.ID)

			ph.Peerstore().AddAddrs(p.ID(), p.Addrs(), pstore.PermanentAddrTTL)
			if err := ph.Connect(ctx, p.Peerstore().PeerInfo(p.ID())); err != nil {
				log.Debugf("failed to bootstrap with %v: %s", p.ID(), err)
				errs <- err
				return
			}
			log.Infof("bootstrapped with %v", p.ID)

		}(p)

	}
	wg.Wait()

	close(errs)
	count := 0
	var err error

	for err = range errs {
		if err != nil {
			count++
		}
	}

	if count == len(peers) {
		return fmt.Errorf("failed to bootstrap %s", err)
	}

	return nil

}

func randomSubsetOfPeers(in []core.Host, max int) []core.Host {
	n := math2.IntMin(max, len(in))
	var out []core.Host
	for _, val := range rand.Perm(len(in)) {
		out = append(out, in[val])
		if(len(out) >= n) {
			break
		}
	}
	return out
}