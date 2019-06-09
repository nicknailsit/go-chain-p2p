package p2p

import "github.com/libp2p/go-libp2p-core/protocol"

const (
	SyncProtocolID = protocol.ID("/sync/1.0.0")
	AuthProtocolID = protocol.ID("/auth/1.0.0")
	WalletProtocolID = protocol.ID("/auth/1.0.0")
	TransactionProtocolID = protocol.ID("/tx/1.0.0")
	ContentProtocolID = protocol.ID("/content/1.0.0")
	FloodSubProtocolID = protocol.ID("/floodsub/1.0.0")
	MDNSProtocolID = protocol.ID("/mdns/1.0.0")
)
