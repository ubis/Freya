package def

import (
    "share/logger"
    "share/rpc"
)

var log = logger.Init("loginserver")

var ServerConfig   = &Config{}
var ServerSettings = &Settings{}
var RPCHandler     = &rpc.Client{}