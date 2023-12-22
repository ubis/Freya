package main

import (
	"github.com/ubis/Freya/cmd/loginserver/packet"
	"github.com/ubis/Freya/cmd/loginserver/server"
	"github.com/ubis/Freya/share/event"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/network"
)

type EventFunc func(*server.Instance, *event.Event)

func register(inst *server.Instance, name string, method EventFunc) {
	event.Register(name, event.Handler(func(e *event.Event) {
		method(inst, e)
	}))
}

// Registers server events
func RegisterEvents(inst *server.Instance) {
	log.Info("Registering events...")

	register(inst, event.ClientConnectEvent, OnClientConnect)
	register(inst, event.ClientConnectEvent, OnClientConnect)
	register(inst, event.ClientDisconnectEvent, OnClientDisconnect)
	register(inst, event.PacketReceiveEvent, OnPacketReceive)
	register(inst, event.PacketSendEvent, OnPacketSend)
	register(inst, event.SyncConnectEvent, OnSyncConnect)
	register(inst, event.SyncDisconnectEvent, OnSyncDisconnect)
}

func parseSession(e *event.Event) (*network.Session, bool) {
	rawSession, ok := e.Get()
	if !ok {
		return nil, false
	}

	session, ok := rawSession.(*network.Session)
	return session, ok
}

func parsePacketArgs(e *event.Event) (*network.PacketArgs, bool) {
	rawPacket, ok := e.Get()
	if !ok {
		return nil, false
	}

	session, ok := rawPacket.(*network.PacketArgs)
	return session, ok
}

// OnClientConnect event informs server about new connected client
func OnClientConnect(i *server.Instance, e *event.Event) {
	session, ok := parseSession(e)
	if !ok {
		return
	}

	session.Store(packet.NewSession(session, i))

	log.Infof("Client `%s` connected to the LoginServer", session.GetEndPnt())
}

// OnClientDisconnect event informs server about disconnected client
func OnClientDisconnect(i *server.Instance, e *event.Event) {
	session, ok := parseSession(e)
	if !ok {
		return
	}

	log.Infof("Client `%s` disconnected from the LoginServer", session.GetEndPnt())
}

// OnPacketReceive event informs server about received packet
func OnPacketReceive(i *server.Instance, e *event.Event) {
	handler := i.PacketHandler

	a, ok := parsePacketArgs(e)
	if !ok {
		return
	}

	log.Debugf("Received Packet `%s` (Len: %d, type: %d, src: %s)",
		handler.Name(a.Type), a.Length, a.Type, a.Session.GetEndPnt())

	// let it handle
	handler.Handle(a)
}

// OnPacketSend event informs server about sent packet
func OnPacketSend(i *server.Instance, e *event.Event) {
	handler := i.PacketHandler

	a, ok := parsePacketArgs(e)
	if !ok {
		return
	}

	log.Debugf("Sent Packet `%s` (Len: %d, type: %d, src: %s)",
		handler.Name(a.Type), a.Length, a.Type, a.Session.GetEndPnt())
}

// OnSyncConnect event informs server about successful connection with the Master Server
func OnSyncConnect(i *server.Instance, e *event.Event) {
	log.Info("Established connection with the Master Server!")

	// register this server
	// req := server.RegisterReq{Type: server.LOGIN_SERVER}
	// res := server.RegisterRes{}
	// i.RPC.Call(rpc.ServerRegister, &req, &res)
}

// OnSyncDisconnect event informs server about lost connection with the Master Server
func OnSyncDisconnect(i *server.Instance, e *event.Event) {
	log.Info("Lost connection with the Master Server! Reconnecting...")
}
