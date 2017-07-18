package packet

import (
    "share/logger"
    "loginserver/def"
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
    pk.Register(CONNECT2SVR, "Connect2Svr", Connect2Svr)
    pk.Register(VERIFYLINKS, "VerifyLinks", VerifyLinks)
    pk.Register(AUTHACCOUNT, "AuthAccount", AuthAccount)
    pk.Register(SYSTEMMESSG, "SystemMessg", nil)
    pk.Register(SERVERSTATE, "ServerState", nil)
    pk.Register(CHECKVERSION, "CheckVersion", CheckVersion)
    pk.Register(URLTOCLIENT, "URLToClient", nil)
    pk.Register(PUBLIC_KEY, "PublicKey", PublicKey)
    pk.Register(PRE_SERVER_ENV_REQUEST, "PreServerEnvRequest", PreServerEnvRequest)
}