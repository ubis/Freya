package def

import (
    "share/logger"
    "share/rpc"
    "masterserver/server"
    "masterserver/database"
    "github.com/jmoiron/sqlx"
)

var log = logger.Init("masterserver")

var ServerConfig    = &Config{}
var ServerSettings  = &Settings{}
var RPCHandler      = &rpc.Server{}
var LoginDatabase   = &sqlx.DB{}
var ServerManager   = &server.ServerManager{}
var DatabaseManager = &database.DatabaseManager{}