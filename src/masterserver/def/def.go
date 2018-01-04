package def

import (
	"github.com/jmoiron/sqlx"
	"masterserver/data"
	"masterserver/database"
	"masterserver/server"
	"share/logger"
	"share/rpc"
)

var log = logger.Init("masterserver")

var ServerConfig = &Config{}
var ServerSettings = &Settings{}
var RPCHandler = &rpc.Server{}
var LoginDatabase = &sqlx.DB{}
var ServerManager = &server.ServerManager{}
var DatabaseManager = &database.DatabaseManager{}
var DataLoader = &data.Loader{}
