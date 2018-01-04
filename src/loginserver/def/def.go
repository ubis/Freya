package def

import (
	"share/network"
	"share/rpc"
)

var ServerConfig = &Config{}
var ServerSettings = &Settings{}
var NetworkManager = &network.Network{}
var PacketHandler = &network.PacketHandler{}
var RPCHandler = &rpc.Client{}
