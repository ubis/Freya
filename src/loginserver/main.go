package main

import (
	"share/logger"
	"share/network"
)

var log = logger.Init("loginserver")

var g_ServerConfig = Config{}

func main() {
	log.Info("LoginServer init")
	g_ServerConfig.Read()

	RegisterEvents()

	network.Init(g_ServerConfig.Port)
}
