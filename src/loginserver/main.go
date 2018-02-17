package main

import (
	"loginserver/internal"
	"loginserver/net"
	"share/log"
	"share/network"
	"share/rpc"
)

func main() {
	// init
	log.Init("LoginServer")

	conf := &Config{}
	rpc := &rpc.Client{}
	network := &network.Server{}
	packets := &net.Packet{RPC: rpc}
	events := &events{rpc: rpc, lst: packets}
	internal := &internal.Comm{Net: network, Lst: packets}

	conf.Read()          // read config
	conf.Assign(packets) // assign config for Packet structure

	rpc.Init(conf.MasterIP, conf.MasterPort) // init RPC client handler
	network.Init(conf.Port)                  // init network server

	packets.Register()     // register server packets
	events.Register()      // register server events
	internal.Register(rpc) // register RPC calls

	rpc.Start()   // start RPC client
	network.Run() // run network server
}
