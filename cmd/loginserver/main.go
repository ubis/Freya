package main

import (
	"github.com/ubis/Freya/cmd/loginserver/packet"
	"github.com/ubis/Freya/cmd/loginserver/rpc"
	"github.com/ubis/Freya/cmd/loginserver/server"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/script"
)

func main() {
	inst := server.NewInstance()
	conf := inst.Config

	log.Init("loginserver")

	// read config
	conf.Read()

	// setup server instance
	inst.XorKeyTable.Init() // initialize xor key table
	inst.RSA.Init()         // create RSA keypair

	RegisterEvents(inst) // register events

	inst.PacketHandler.Init() // init packet handler

	// register scripting engine
	script.Initialize(conf.ScriptDirectory)

	// register packets
	packet.RegisterPackets(inst.PacketHandler)

	// register scripting functions
	packet.RegisterFunc()

	// setup RPC subsystem
	inst.RPC.Init()
	inst.RPC.IpAddress = conf.MasterIp
	inst.RPC.Port = conf.MasterPort

	// register RPC calls
	rpc.RegisterCalls(inst.RPC, inst.Server)

	// start RPC handler
	inst.RPC.Start()

	// create network and start listening for connections
	inst.Server.Init(conf.Port, inst.XorKeyTable)
}
