package rpc

import (
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

type RPC struct {
	*rpc.Client

	Server *network.Network
}

func RegisterCalls(r *rpc.Client, net *network.Network) {
	inst := &RPC{Client: r, Server: net}

	r.Register(rpc.UserVerify, inst.UserVerify)
	r.Register(rpc.OnlineCheck, inst.OnlineCheck)
}
