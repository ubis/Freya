package def

import (
    "share/logger"
    "share/rpc"
    "share/network"
)

var log = logger.Init("loginserver")

var ServerConfig   = &Config{}
var ServerSettings = &Settings{}
var PacketHandler  = &network.PacketHandler{}
var RPCHandler     = &rpc.Client{}