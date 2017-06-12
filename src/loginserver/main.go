package main

import (
    "share/logger"
    "share/network"
    "loginserver/packet"
)

var log = logger.Init("loginserver")

var g_ServerConfig   = &Config{}
var g_ServerSettings = &Settings{}
var g_PacketHandler  = &packet.PacketHandler{}

func main() {
    log.Info("LoginServer init")

    // read config
    g_ServerConfig.Read()

    // set server settings
    g_ServerSettings.Global.ListenPort = g_ServerConfig.Port
    g_ServerSettings.Global.MaxUsers   = g_ServerConfig.MaxUsers
    g_ServerSettings.Global.XorKeyTable.Init()

    // register events
    RegisterEvents()

    // init packet handler
    g_PacketHandler.Init()

    // create network and start listening for connections
    network.Init(g_ServerSettings.Global)
}
