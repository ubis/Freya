package server

import (
	"github.com/ubis/Freya/share/encryption"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

type Instance struct {
	Config        *Config
	Server        *network.Network
	PacketHandler *network.PacketHandler
	XorKeyTable   *encryption.XorKeyTable
	RSA           *RSA
	RPC           *rpc.Client
}

func NewInstance() *Instance {
	return &Instance{
		Config:        &Config{},
		Server:        &network.Network{},
		PacketHandler: &network.PacketHandler{},
		XorKeyTable:   &encryption.XorKeyTable{},
		RSA:           &RSA{},
		RPC:           &rpc.Client{},
	}
}
