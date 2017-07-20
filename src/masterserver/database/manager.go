package database

import (
    "share/logger"
    "github.com/jmoiron/sqlx"
)

var log = logger.Instance()

type DatabaseManager struct {
    DBList map[int]*Database
}

/*
    Initializes Database Manager which will attempt to connect
     to all GameDatabases
    @param  db  map of game databases
 */
func (dm *DatabaseManager) Init(db map[int]*Database) {
    dm.DBList = db

    for i := 1; i < len(dm.DBList); i ++ {
        if dm.DBList[i] != nil {
            dm.connect(dm.DBList[i])
        }
    }
}

/*
    Returns game database instance of index
    @param  index   game database index
    @return pointer to sqlx.DB or nil if not found
 */
func (dm *DatabaseManager) Get(index int) *sqlx.DB {
    if dm.DBList[index] == nil {
        log.Error("Game #%d database doesn't exist!", index)
        return nil
    }

    return dm.DBList[index].DB
}

/*
    Attempts to connect to specified database
    @param  dba pointer to database structure
 */
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