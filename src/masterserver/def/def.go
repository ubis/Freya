package def

import (
    "share/logger"
    "share/rpc"
    "github.com/jmoiron/sqlx"
)

var log = logger.Init("masterserver")

var ServerConfig   = &Config{}
var ServerSettings = &Settings{}
var RPCHandler     = &rpc.Server{}
var LoginDatabase  = &sqlx.DB{}