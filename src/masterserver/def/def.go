package def

import (
    "share/logger"
    "share/rpc"
)

var log = logger.Init("masterserver")

var ServerConfig   = &Config{}
var ServerSettings = &Settings{}
var RPCHandler     = &rpc.Server{}