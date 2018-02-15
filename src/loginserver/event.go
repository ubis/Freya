package main

import (
	"share/event"
	"share/log"
	"share/models/server"
	"share/network"
	"share/rpc"
)

// EventManager structure
type EventManager struct {
	rpc *rpc.Client
	net *network.Manager
}

// Register attempts to register server events
func (em *EventManager) Register() {
	log.Info("Registering events...")

	event.Register(event.ClientConnect, event.Handler(em.OnClientConnect))
	event.Register(event.ClientDisconnect, event.Handler(em.OnClientDisconnect))
	event.Register(event.PacketReceive, event.Handler(em.OnPacketReceive))
	event.Register(event.PacketSend, event.Handler(em.OnPacketSend))
	event.Register(event.SyncConnect, event.Handler(em.OnSyncConnect))
	event.Register(event.SyncDisconnect, event.Handler(em.OnSyncDisconnect))
}

// OnClientConnect event informs server about new connected client
func (em *EventManager) OnClientConnect(e event.Event) {
	if s, ok := e.(*network.Session); !ok {
		log.Error("Cannot parse onClientConnect event!")
	} else {
		log.Infof("Client `%s` connected to the LoginServer", s.GetEndPnt())
	}
}

// OnClientDisconnect event informs server about disconnected client
func (em *EventManager) OnClientDisconnect(e event.Event) {
	if s, ok := e.(*network.Session); !ok {
		log.Error("Cannot parse onClientDisconnect event!")
	} else {
		log.Infof("Client `%s` disconnected from the LoginServer", s.GetEndPnt())
	}
}

// OnPacketReceive event informs server about received packet
func (em *EventManager) OnPacketReceive(e event.Event) {
	if a, ok := e.(*network.PacketArgs); !ok {
		log.Error("Cannot parse onPacketReceive event!")
	} else {
		log.Debugf("Received Packet `%s` (Len: %d, type: %d, src: %s)",
			em.net.GetPacketName(a.Type), a.Length, a.Type, a.Session.GetEndPnt())

		// let it handle
		em.net.HandlePacket(a)
	}
}

// OnPacketSend event informs server about sent packet
func (em *EventManager) OnPacketSend(e event.Event) {
	if a, ok := e.(*network.PacketArgs); !ok {
		log.Error("Cannot parse onPacketSent event!")
	} else {
		log.Debugf("Sent Packet `%s` (Len: %d, type: %d, src: %s)",
			em.net.GetPacketName(a.Type), a.Length, a.Type, a.Session.GetEndPnt())
	}
}

// OnSyncConnect event informs server about successful connection with
// the Master Server
func (em *EventManager) OnSyncConnect(e event.Event) {
	log.Info("Established connection with the Master Server!")

	// register this server
	var q = server.RegisterReq{Type: server.LOGIN_SERVER}
	em.rpc.Call(rpc.ServerRegister, q, &server.RegisterRes{})
}

// OnSyncDisconnect event informs server about lost connection with
// the Master Server
func (em *EventManager) OnSyncDisconnect(e event.Event) {
	log.Info("Lost connection with the Master Server! Reconnecting...")
}
