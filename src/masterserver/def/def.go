package def

import (
	"masterserver/data"
	"masterserver/database"
	"masterserver/server"
	"share/rpc"

	"github.com/jmoiron/sqlx"
)

var ServerConfig = &Config{}
var ServerSettings = &Settings{}
var RPCHandler = &rpc.Server{}
var LoginDatabase = &sqlx.DB{}
var ServerManager = &server.ServerManager{}
var DatabaseManager = &database.DatabaseManager{}
var DataLoader = &data.Loader{}
