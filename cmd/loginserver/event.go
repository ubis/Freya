package main

import (
	"github.com/ubis/Freya/share/event"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/server"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

// Registers server events
func RegisterEvents() {
	log.Info("Registering events...")
	event.Register(event.ClientConnectEvent, event.Handler(OnClientConnect))
	event.Register(event.ClientDisconnectEvent, event.Handler(OnClientDisconnect))
	event.Register(event.PacketReceiveEvent, event.Handler(OnPacketReceive))
	event.Register(event.PacketSendEvent, event.Handler(OnPacketSend))
	event.Register(event.SyncConnectEvent, event.Handler(OnSyncConnect))
	event.Register(event.SyncDisconnectEvent, event.Handler(OnSyncDisconnect))
}

// OnClientConnect event informs server about new connected client
func OnClientConnect(e event.Event) {
	if s, ok := e.(*network.Session); !ok {
		log.Error("Cannot parse onClientConnect event!")
	} else {
		log.Infof("Client `%s` connected to the LoginServer", s.GetEndPnt())
	}
}

// OnClientDisconnect event informs server about disconnected client
func OnClientDisconnect(e event.Event) {
	if s, ok := e.(*network.Session); !ok {
		log.Error("Cannot parse onClientDisconnect event!")
	} else {
		log.Infof("Client `%s` disconnected from the LoginServer", s.GetEndPnt())
	}
}

// OnPacketReceive event informs server about received packet
func OnPacketReceive(e event.Event) {
	if a, ok := e.(*network.PacketArgs); !ok {
		log.Error("Cannot parse onPacketReceive event!")
	} else {
		log.Debugf("Received Packet `%s` (Len: %d, type: %d, src: %s)",
			g_PacketHandler.Name(a.Type), a.Length, a.Type, a.Session.GetEndPnt(),
		)

		// let it handle
		g_PacketHandler.Handle(a)
	}
}

// OnPacketSend event informs server about sent packet
func OnPacketSend(e event.Event) {
	if a, ok := e.(*network.PacketArgs); !ok {
		log.Error("Cannot parse onPacketSent event!")
	} else {
		log.Debugf("Sent Packet `%s` (Len: %d, type: %d, src: %s)",
			g_PacketHandler.Name(a.Type), a.Length, a.Type, a.Session.GetEndPnt(),
		)
	}
}

// OnSyncConnect event informs server about successful connection with the Master Server
func OnSyncConnect(e event.Event) {
	log.Info("Established connection with the Master Server!")

	// register this server
	var q = server.RegisterReq{Type: server.LOGIN_SERVER}
	g_RPCHandler.Call(rpc.ServerRegister, q, &server.RegisterRes{})
}

// OnSyncDisconnect event informs server about lost connection with the Master Server
func OnSyncDisconnect(e event.Event) {
	log.Info("Lost connection with the Master Server! Reconnecting...")
}
