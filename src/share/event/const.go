package event

// Server events for client communication
const (
	ClientConnect    = "onClientConnect"
	ClientDisconnect = "onClientDisconnect"
)

// Server events for network
const (
	PacketReceive = "onPacketReceive"
	PacketSend    = "onPacketSent"
)

// Server events for RPC
const (
	SyncConnect    = "onSyncConnect"
	SyncDisconnect = "onSyncDisconnect"
)
