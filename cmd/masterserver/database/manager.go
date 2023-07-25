package database

import (
	"github.com/ubis/Freya/share/log"

	"github.com/jmoiron/sqlx"
)

type DatabaseManager struct {
	DBList map[int]*Database
}

// Initializes Database Manager which will attempt to connect
func (dm *DatabaseManager) Init(db map[int]*Database) {
	dm.DBList = db

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
