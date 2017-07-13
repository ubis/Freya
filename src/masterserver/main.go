package main

import (
    "share/logger"
    "masterserver/def"
    "masterserver/rpc"
)

var log = logger.Instance()

var g_ServerConfig   = def.ServerConfig
var g_ServerSettings = def.ServerSettings
var g_RPCHandler     = def.RPCHandler


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

    // start RPC Server
    g_RPCHandler.Run()
}