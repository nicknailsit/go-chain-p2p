package p2p

import (
	net "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"time"
)

func (client *client) auth(peerId peer.ID) (bool, error){

	pid := peerId

	client.logger.Debug("authenticating to the swaggit network", pid)

	stream, err := client.host.NewStream(client.context, pid, client.protocol+"/auth")

	if err != nil {
		addrs := client.peerstore.PeerInfo(pid).Addrs
		client.logger.Debug("cannot connect to", pid, "at", addrs, err)
		client.peerstore.ClearAddrs(pid)
		client.serviceDiscovery.UnregisterNotifee(client.serviceNotifee)
		return false, err
	}

	defer stream.Close()

	commitmentOut := client.requestCommitment()
	err = client.sendCommitment(stream, commitmentOut)
	if err != nil {
		return false, err
	}

	challengeIn, err := client.receiveChallenge(stream)
	if err != nil {
		return false, err
	}

	proofOut := client.requestProof(commitmentOut, challengeIn)
	err = client.sendProof(stream, proofOut)
	if err != nil {
		return false, err
	}

	challengeOut := client.requestChallenge()
	err = client.sendChallenge(stream, challengeOut)


	proofIn, err := client.receiveProof(stream)
	if err != nil {
		return false, err
	}

	success := client.requestVerification(commitmentIn, challengeOut, proofIn)
	return success, nil

}

func (client *client) authHandler(stream net.Stream) {

	defer stream.Close()
	pid := stream.Conn().RemotePeer()
	client.logger.Debug("authenticating", pid)
	client.spammerCacheLock.Lock()

	timestamp, exists := client.spammerCache.Get(pid)
	if exists && time.Since(timestamp.(time.Time)) < 10*time.Minute {
		client.spammerCacheLock.Unlock()
		time.Sleep(client.config.Timeout)
		return
	}
	client.spammerCache.Add(pid, time.Now)
	client.spammerCacheLock.Unlock()

	commitmentIn, err := client.receiveCommitment(stream)
	if err != nil {
		return
	}

	challengeOut := client.requestChallenge()
	err = client.sendChallenge(stream, challengeOut)
	if err != nil {
		return
	}

	proofIn, err := client.receiveProof(stream)
	if err != nil {
		return
	}

	success := client.requestVerification(commitmentIn, challengeOut, proofIn)
	if !success {
		return
	}

	client.logger.Debug("authenticating", pid)

	commitmentOut := client.requestCommitment()
	err = client.sendCommitment(stream, commitmentOut)
	if err != nil {
		return
	}

	challengeIn := client.receiveChallenge(stream)
	if err != nil {
		return
	}

	proofOut := client.requestProof(commitmentOut, challengeIn)
	err = client.sendProof(stream, proofOut)

	if err != nil {
		return
	}
}

func (client *client) registerAuthService() {
	uri := client.protocol + "/auth"
	client.host.SetStreamHandler(uri, client.authHandler)
}