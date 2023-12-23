package database

import (
	"github.com/ubis/Freya/cmd/masterserver/server"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/rpc"

	"github.com/jmoiron/sqlx"
)

type DatabaseManager struct {
	DBList map[int]*Database
	Server *server.ServerManager
}

// Initializes Database Manager which will attempt to connect
func (dm *DatabaseManager) Init(svr *server.ServerManager, db map[int]*Database) {
	dm.DBList = db
	dm.Server = svr

	for _, value := range dm.DBList {
		dm.connect(value)
	}
}

// Returns game database instance of index
func (dm *DatabaseManager) Get(index byte) *sqlx.DB {
	if dm.DBList[int(index)] == nil {
		log.Error("Game #%d database doesn't exist!", index)
		return nil
	}

	return dm.DBList[int(index)].DB
}

func (dm *DatabaseManager) Find(c *rpc.Client) *sqlx.DB {
	server := dm.Server.FindServer(c)

	if server == nil {
		log.Error("Unable to find GameWorld database!")
		return nil
	}

	if dm.DBList[int(server.ServerId)] == nil {
		log.Errorf("Game #%d database doesn't exist!", server.ServerId)
		return nil
	}

	return dm.DBList[int(server.ServerId)].DB
}

// Attempts to connect to specified database
func (dm *DatabaseManager) connect(dba *Database) {
	log.Infof("Attempting to connect to the #%d Game database...", dba.Index)
	if db, err := sqlx.Connect("mysql", dba.Config); err != nil {
		log.Fatalf("[DATABASE] %s", err.Error())
	} else {
		log.Infof("Successfully connected to the #%d Game database!", dba.Index)
		dba.DB = db

		var version []string
		db.Select(&version, "SELECT VERSION()")
		log.Debugf("[DATABASE] Version: %s", version[0])
	}
}
