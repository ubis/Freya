package rpc

import (
    "share/logger"
    "share/rpc"
    "masterserver/def"
)

var log = logger.Instance()

var g_RPCHandler    = def.RPCHandler
var g_LoginDatabase = def.LoginDatabase
var g_ServerManager = def.ServerManager

func RegisterPackets() {
    log.Info("Registering RPC packets...")
    g_RPCHandler.Register(rpc.ServerRegister, ServerRegister)
    g_RPCHandler.Register(rpc.AuthCheck, AuthCheck)
}