package event

// Server events for client communication
const (
	ClientConnectEvent    = "onClientConnect"
	ClientDisconnectEvent = "onClientDisconnect"
)

// Server events for network
const (
	PacketReceiveEvent = "onPacketReceive"
	PacketSendEvent    = "onPacketSent"
)

// Server events for RPC
const (
	SyncConnectEvent    = "onSyncConnect"
	SyncDisconnectEvent = "onSyncDisconnect"
)
