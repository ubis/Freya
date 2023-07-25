package def

import (
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

var ServerConfig = &Config{}
var ServerSettings = &Settings{}
var NetworkManager = &network.Network{}
var PacketHandler = &network.PacketHandler{}
var RPCHandler = &rpc.Client{}
