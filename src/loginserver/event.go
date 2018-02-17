package main

import (
	"loginserver/net"
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
}

// Register attempts to register server events
func (e *events) Register() {
	log.Info("Registering events...")

	event.Register(event.ClientConnect, event.Handler(e.onClientConnect))
	event.Register(event.ClientDisconnect, event.Handler(e.onClientDisconnect))
	event.Register(event.PacketReceive, event.Handler(e.onPacketReceive))
	event.Register(event.PacketSend, event.Handler(e.onPacketSend))
	event.Register(event.SyncConnect, event.Handler(e.onSyncConnect))
	event.Register(event.SyncDisconnect, event.Handler(e.onSyncDisconnect))
}

// onClientConnect event informs server about new connected client
func (e *events) onClientConnect(ev event.Event) {
	s, ok := ev.(*network.Session)
	if !ok {
		return
	}

	log.Infof("Client `%s` connected to the LoginServer", s.GetEndPnt())
}

// onClientDisconnect event informs server about disconnected client
func (e *events) onClientDisconnect(ev event.Event) {
	s, ok := ev.(*network.Session)
	if !ok {
		return
	}

	log.Infof("Client `%s` disconnected from the LoginServer", s.GetEndPnt())
}

// onPacketReceive event informs server about received packet
func (e *events) onPacketReceive(ev event.Event) {
	a, ok := ev.(*network.PacketArgs)
	if !ok {
		return
	}

	log.Debugf("Received Packet `%s` (Len: %d, type: %d, src: %s)",
		e.lst.GetName(a.Type), a.Length, a.Type, a.Session.GetEndPnt())

	// let it handle
	e.lst.Handle(a)
}

// onPacketSend event informs server about sent packet
func (e *events) onPacketSend(ev event.Event) {
	a, ok := ev.(*network.PacketArgs)
	if !ok {
		return
	}

	log.Debugf("Sent Packet `%s` (Len: %d, type: %d, src: %s)",
		e.lst.GetName(a.Type), a.Length, a.Type, a.Session.GetEndPnt())
}

// onSyncConnect event informs server about successful connection with
// the Master Server
func (e *events) onSyncConnect(ev event.Event) {
	log.Info("Established connection with the Master Server!")

	// register this server
	req := server.RegisterReq{Type: server.LOGIN_SERVER}
	e.rpc.Call(rpc.ServerRegister, req, &server.RegisterRes{})
}

// onSyncDisconnect event informs server about lost connection with
// the Master Server
func (e *events) onSyncDisconnect(ev event.Event) {
	log.Info("Lost connection with the Master Server! Reconnecting...")
}
