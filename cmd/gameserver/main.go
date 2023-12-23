package main

import (
	"os"
	"strconv"

	"github.com/ubis/Freya/cmd/gameserver/game"
	"github.com/ubis/Freya/cmd/gameserver/packet"
	"github.com/ubis/Freya/cmd/gameserver/rpc"
	"github.com/ubis/Freya/cmd/gameserver/server"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/script"
)

func readServerParam(i *server.Instance) {
	i.ServerId = 1
	i.ChannelId = 1

	if len(os.Args) < 2 {
		return
	}

	if id, err := strconv.Atoi(os.Args[1]); err == nil {
		i.ServerId = id
	}

	if id, err := strconv.Atoi(os.Args[2]); err == nil {
		i.ChannelId = id
	}
}

func main() {
	inst := server.NewInstance()
	conf := inst.Config

	readServerParam(inst)

	log.Init(inst.GetName())

	// read config
	conf.Read(inst.GetName())

	// set up world manager
	game := &game.WorldManager{}
	game.Initialize()

	// setup server instance
	inst.XorKeyTable.Init() // initialize xor key table

	RegisterEvents(inst, game) // register events

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
