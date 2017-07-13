package rpc

import (
    "share/logger"
    "share/rpc"
    "masterserver/def"
)

var log = logger.Instance()

var g_ServerConfig   = def.ServerConfig
var g_ServerSettings = def.ServerSettings
var g_RPCHandler     = def.RPCHandler

func RegisterPackets() {
    g_RPCHandler.Register(rpc.ServerRegister, ServerRegister)
}