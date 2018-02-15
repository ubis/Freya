package internal

import (
	"loginserver/net"
	"share/network"
	"share/rpc"
)

type Manager struct {
	Network *network.Manager
	Packets *net.Packet
}

// Register RPC calls
func (m *Manager) Register(r *rpc.Client) {
	r.Register(rpc.UserVerify, m.UserVerify)
	r.Register(rpc.OnlineCheck, m.OnlineCheck)
}
