package main

import (
	"encoding/binary"
	"gameserver/net"
	netutil "net"
	"share/event"
	"share/log"
	"share/models/server"
	"share/network"
	"share/rpc"
)

// events structure
type events struct {
	rpc *rpc.Client
	lst *net.Packet
	svr *network.Server
	cfg *Config
}

// Register server events
func (e *events) Register() {
	log.Info("Registering events...")

	event.Register(event.ClientConnect, event.Handler(e.OnClientConnect))
	event.Register(event.ClientDisconnect, event.Handler(e.OnClientDisconnect))
	event.Register(event.PacketReceive, event.Handler(e.OnPacketReceive))
	event.Register(event.PacketSend, event.Handler(e.OnPacketSend))
	event.Register(event.SyncConnect, event.Handler(e.OnSyncConnect))
	event.Register(event.SyncDisconnect, event.Handler(e.OnSyncDisconnect))
}

// OnClientConnect event informs server about new connected client
func (e *events) OnClientConnect(ev event.Event) {
	s, ok := ev.(*network.Session)
	if !ok {
		return
	}

	log.Infof("Client `%s` connected to the GameServer", s.GetEndPnt())
}

// OnClientDisconnect event informs server about disconnected client
func (e *events) OnClientDisconnect(ev event.Event) {
	s, ok := ev.(*network.Session)
	if !ok {
		return
	}

	log.Infof("Client `%s` disconnected from the GameServer", s.GetEndPnt())
}

// OnPacketReceive event informs server about received packet
func (e *events) OnPacketReceive(ev event.Event) {
	a, ok := ev.(*network.PacketArgs)
	if !ok {
		return
	}

	log.Debugf("Received Packet `%s` (Len: %d, type: %d, src: %s)",
		e.lst.GetName(a.Type), a.Length, a.Type, a.Session.GetEndPnt())

	// let it handle
	e.lst.Handle(a)
}

// OnPacketSend event informs server about sent packet
func (e *events) OnPacketSend(ev event.Event) {
	a, ok := ev.(*network.PacketArgs)
	if !ok {
		return
	}

	log.Debugf("Sent Packet `%s` (Len: %d, type: %d, src: %s)",
		e.lst.GetName(a.Type), a.Length, a.Type, a.Session.GetEndPnt())
}

// OnSyncConnect event informs server about successful connection with
// the Master Server
func (e *events) OnSyncConnect(ev event.Event) {
	log.Info("Established connection with the Master Server!")

	// register this server
	req := server.RegisterReq{
		server.GAME_SERVER,
		byte(e.cfg.ServerType),
		byte(e.cfg.ServerID),
		byte(e.cfg.GroupID),
		binary.LittleEndian.Uint32(netutil.ParseIP(e.cfg.PublicIp)[12:16]),
		uint16(e.cfg.Port),
		uint16(e.svr.GetSessionCount()),
		uint16(e.cfg.MaxUsers),
	}
	res := server.RegisterRes{}

	e.rpc.Call(rpc.ServerRegister, req, &res)
}

// OnSyncDisconnect event informs server about lost connection with
// the Master Server
func (e *events) OnSyncDisconnect(ev event.Event) {
	log.Info("Lost connection with the Master Server! Reconnecting...")
}
