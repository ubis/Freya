package main

import (
	"loginserver/internal"
	"loginserver/net"
	"loginserver/rsa"
	"share/log"
	"share/network"
	"share/rpc"
)

func main() {
	// init
	conf := &Config{}
	rsa := &rsa.RSA{}
	sync := &rpc.Client{}
	packets := &net.Packet{RPC: sync, RSA: rsa}
	network := &network.Manager{}
	events := &EventManager{rpc: sync, net: network}

	internal := &internal.Manager{Network: network, Packets: packets}

	log.Init("LoginServer")

	// read config
	conf.Read()

	// assign config for Packet structure
	conf.Assign(packets)

	// initialize encryption
	rsa.Init()

	// init network manager
	network.Init(conf.Port)

	// register packets
	packets.Register(network)

	// register events
	events.Register()

	// init RPC handler
	sync.Init()
	sync.IpAddress = conf.MasterIP
	sync.Port = conf.MasterPort

	// register RPC calls
	internal.Register(sync)

	// start RPC handler
	sync.Start()

	// run network server
	network.Run()
}
