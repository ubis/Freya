package main

import (
	"gameserver/internal"
	"gameserver/net"
	"os"
	"share/log"
	"share/network"
	"share/rpc"
	"strconv"
)

var serverID = 1
var groupID = 1

// init function, which runs before main()
func init() {
	if len(os.Args) > 2 {
		if id, err := strconv.Atoi(os.Args[1]); err == nil {
			serverID = id
		}

		if id, err := strconv.Atoi(os.Args[2]); err == nil {
			groupID = id
		}
	}
}

func main() {
	// init
	conf := &Config{ServerID: serverID, GroupID: groupID}
	rpc := &rpc.Client{}
	network := &network.Server{}
	packets := &net.Packet{RPC: rpc}
	events := &events{rpc: rpc, lst: packets, svr: network, cfg: conf}
	internal := &internal.Comm{Net: network, Lst: packets}

	log.Init(conf.GetName()) // log init

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
