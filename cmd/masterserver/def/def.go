package def

import (
	"github.com/ubis/Freya/cmd/masterserver/data"
	"github.com/ubis/Freya/cmd/masterserver/database"
	"github.com/ubis/Freya/cmd/masterserver/server"
	"github.com/ubis/Freya/share/rpc"

	"github.com/jmoiron/sqlx"
)

var ServerConfig = &Config{}
var ServerSettings = &Settings{}
var RPCHandler = &rpc.Server{}
var LoginDatabase = &sqlx.DB{}
var ServerManager = &server.ServerManager{}
var DatabaseManager = &database.DatabaseManager{}
var DataLoader = &data.Loader{}
