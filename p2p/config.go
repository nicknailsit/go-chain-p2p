package p2p

import "time"

type config struct {
	RendezvousString string
	ProtocolID string
	listenHost string
	listenPort string
}

type ClientConfig struct {
	AnalyticsInterval           time.Duration
	AnalyticsURL                string
	AnalyticsUserData           string
	ArtifactCacheSize           int
	ArtifactChunkSize           uint32
	ArtifactMaxBufferSize       uint32
	ArtifactQueueSize           int
	ChallengeMaxBufferSize      uint32
	ClusterID                   int
	CommitmentMaxBufferSize     uint32
	DisableAnalytics            bool
	DisableBroadcast            bool
	DisableNATPortMap           bool
	DisablePeerDiscovery        bool
	DisableStreamDiscovery      bool
	IP                          string
	KBucketSize                 int
	LatencyTolerance            time.Duration
	NATMonitorInterval          time.Duration
	NATMonitorTimeout           time.Duration
	Network                     string
	PingBufferSize              uint32
	Port                        uint16
	ProcessID                   int
	ProofMaxBufferSize          uint32
	RandomSeed                  string
	SampleMaxBufferSize         uint32
	SampleSize                  int
	SeedNodes                   []string
	SpammerCacheSize            int
	StreamstoreInboundCapacity  int
	StreamstoreOutboundCapacity int
	StreamstoreQueueSize        int
	Timeout                     time.Duration
	Version                     string
	WitnessCacheSize            int
}

// DefaultConfig -- Get the default configuration parameters.
func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		AnalyticsInterval:       time.Minute,
		AnalyticsURL:            "",
		AnalyticsUserData:       "",
		ArtifactCacheSize:       65536,
		ArtifactChunkSize:       65536,
		ArtifactMaxBufferSize:   8388608,
		ArtifactQueueSize:       8,
		ChallengeMaxBufferSize:  32,
		ClusterID:               0,
		CommitmentMaxBufferSize: 32,
		DisableAnalytics:        false,
		DisableBroadcast:        false,
		DisableNATPortMap:       false,
		DisablePeerDiscovery:    false,
		DisableStreamDiscovery:  false,
		IP:                          "0.0.0.0",
		KBucketSize:                 16,
		LatencyTolerance:            time.Minute,
		NATMonitorInterval:          time.Second,
		NATMonitorTimeout:           time.Minute,
		Network:                     "swaggchain",
		PingBufferSize:              32,
		Port:                        0,
		ProcessID:                   0,
		ProofMaxBufferSize:          0,
		RandomSeed:                  "",
		SampleMaxBufferSize:         8192,
		SampleSize:                  16,
		SeedNodes:                   nil,
		SpammerCacheSize:            16384,
		StreamstoreInboundCapacity:  48,
		StreamstoreOutboundCapacity: 16,
		StreamstoreQueueSize:        8192,
		Timeout:                     10 * time.Second,
		Version:                     "1.0.0",
		WitnessCacheSize:            65536,
	}
}