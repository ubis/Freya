package event

// Server events for client communication
const (
    ClientConnectEvent      = "onClientConnect"
    ClientDisconnectEvent   = "onClientDisconnect"
)

// Server events for network
const (
    PacketReceivedEvent     = "onPacketReceive"
    PacketSentEvent         = "onPacketSent"
)