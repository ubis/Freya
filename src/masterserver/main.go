package main

import (
    "share/logger"
    "masterserver/def"
    "masterserver/rpc"
    "github.com/jmoiron/sqlx"
    _"github.com/go-sql-driver/mysql"
)

var log = logger.Instance()

// globals
var g_ServerConfig    = def.ServerConfig
var g_RPCHandler      = def.RPCHandler
var g_LoginDatabase   = def.LoginDatabase
var g_ServerManager   = def.ServerManager
var g_DatabaseManager = def.DatabaseManager
var g_DataLoader      = def.DataLoader

func main() {
    log.Info("MasterServer init")

    // read config
    g_ServerConfig.Read()

    // register events
    RegisterEvents()

    // init RPC handler
    g_RPCHandler.Init()
    g_RPCHandler.Port = g_ServerConfig.Port

    // register RPC packets
    rpc.RegisterPackets()

    // init ServerManager
    g_ServerManager.Init()

    // connect to login database
    log.Info("Attempting to connect to the Login database...")
    var cfg = g_ServerConfig.GetDBConfig(g_ServerConfig.LoginDB)
    if db, err := sqlx.Connect("mysql", cfg); err != nil {
        log.Fatalf("[DATABASE] %s", err.Error())
    } else {
        log.Info("Successfully connected to the Login database!")
        *g_LoginDatabase = *db

        var version []string
        db.Select(&version, "SELECT VERSION()")
        log.Debugf("[DATABASE] Version: %s", version[0])
    }

    // init DatabaseManager
    g_DatabaseManager.Init(g_ServerConfig.GameDB)

    // init DataLoader
    g_DataLoader.Init()

    // start RPC Server
    g_RPCHandler.Run()
}