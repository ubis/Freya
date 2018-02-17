package internal

import (
	"loginserver/net"
	"share/network"
	"share/rpc"
)

type Comm struct {
	Net *network.Server
	Lst *net.Packet
}

// Register RPC calls
func (cm *Comm) Register(r *rpc.Client) {
	r.Register(rpc.UserVerify, cm.UserVerify)
	r.Register(rpc.OnlineCheck, cm.OnlineCheck)
}
