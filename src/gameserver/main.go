package main

import (
    "share/logger"
    "gameserver/def"
    "gameserver/packet"
    "gameserver/rpc"
)

var log = logger.Instance()

var g_ServerConfig   = def.ServerConfig
var g_ServerSettings = def.ServerSettings
var g_NetworkManager = def.NetworkManager
var g_PacketHandler  = def.PacketHandler
var g_RPCHandler     = def.RPCHandler

func main() {
    log.Info("GameServer", def.GetName(), " init")

    // read config
    g_ServerConfig.Read()

    // set server settings
    g_ServerSettings.XorKeyTable.Init()

    // register events
    RegisterEvents()

    // init packet handler
    g_PacketHandler.Init()

    // register packets
    packet.RegisterPackets()

    // init RPC handler
    g_RPCHandler.Init()
    g_RPCHandler.IpAddress = g_ServerConfig.MasterIp
    g_RPCHandler.Port      = g_ServerConfig.MasterPort

    // register RPC calls
    rpc.RegisterCalls()

    // start RPC handler
    g_RPCHandler.Start()

    // create network and start listening for connections
    g_NetworkManager.Init(&g_ServerSettings.Settings)
    g_NetworkManager.Start(g_ServerConfig.Port)
}