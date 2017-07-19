package packet

import (
    "share/logger"
    "gameserver/def"
)

var log = logger.Instance()

var g_ServerConfig   = def.ServerConfig
var g_ServerSettings = def.ServerSettings
var g_PacketHandler  = def.PacketHandler
var g_RPCHandler     = def.RPCHandler

// Registers network packets
func RegisterPackets() {
    log.Info("Registering packets...")
    
    var pk = g_PacketHandler
    pk.Register(GETMYCHARTR, "GetMyChartr", GetMyChartr)
    pk.Register(CONNECT2SVR, "Connect2Svr", Connect2Svr)
    pk.Register(GETSVRTIME, "GetSvrTime", GetSvrTime)
    pk.Register(CHARGEINFO, "ChargeInfo", ChargeInfo)
    pk.Register(SERVERENV, "ServerEnv", ServerEnv)
}