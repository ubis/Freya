package server

import (
	"fmt"

	"github.com/ubis/Freya/share/encryption"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

type Instance struct {
	ServerId  int
	ChannelId int

	Config        *Config
	Server        *network.Network
	PacketHandler *network.PacketHandler
	XorKeyTable   *encryption.XorKeyTable
	RPC           *rpc.Client
}

func NewInstance() *Instance {
	return &Instance{
		Config:        &Config{},
		Server:        &network.Network{},
		PacketHandler: &network.PacketHandler{},
		XorKeyTable:   &encryption.XorKeyTable{},
		RPC:           &rpc.Client{},
	}
}

func (i *Instance) GetName() string {
	return fmt.Sprintf("GameServer_%d_%d", i.ServerId, i.ChannelId)
}
