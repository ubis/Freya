package main

import (
	"github.com/ubis/Freya/cmd/gameserver/game"
	"github.com/ubis/Freya/cmd/gameserver/packet"
	"github.com/ubis/Freya/cmd/gameserver/server"
	"github.com/ubis/Freya/share/event"
	"github.com/ubis/Freya/share/log"
	svr "github.com/ubis/Freya/share/models/server"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

type EventFunc func(*server.Instance, *event.Event)

func register(inst *server.Instance, name string, method EventFunc) {
	event.Register(name, event.Handler(func(e *event.Event) {
		method(inst, e)
	}))
}

// Registers server events
func RegisterEvents(inst *server.Instance, wm *game.WorldManager) {
	log.Info("Registering events...")

	event.Register(event.ClientConnectEvent, event.Handler(func(e *event.Event) {
		OnClientConnect(inst, e, wm)
	}))

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
func OnClientConnect(i *server.Instance, e *event.Event, wm *game.WorldManager) {
	session, ok := parseSession(e)
	if !ok {
		return
	}

	session.Store(packet.NewSession(session, i, wm))

	log.Infof("Client `%s` connected to the GameServer", session.GetEndPnt())
}

// OnClientDisconnect event informs server about disconnected client
func OnClientDisconnect(i *server.Instance, e *event.Event) {
	session, ok := parseSession(e)
	if !ok {
		return
	}

	// in case client was in the world, notify other players
	world := packet.GetCurrentWorld(session)
	if world != nil {
		world.ExitWorld(session, svr.DelUserLogout)
	}

	log.Infof("Client `%s` disconnected from the GameServer", session.GetEndPnt())
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
func OnSyncConnect(i *server.Instance, event *event.Event) {
	log.Info("Established connection with the Master Server!")

	// register this server
	req := svr.RegisterReq{
		Type:         svr.GAME_SERVER,
		ServerType:   byte(i.Config.ServerType),
		ServerId:     byte(i.ServerId),
		ChannelId:    byte(i.ChannelId),
		PublicIp:     i.Config.PublicIp,
		PublicPort:   uint16(i.Config.Port),
		UseLocalIp:   i.Config.UseLocalIp,
		CurrentUsers: uint16(i.Server.GetOnlineUsers()),
		MaxUsers:     uint16(i.Config.MaxUsers),
	}
	res := svr.RegisterRes{}

	i.RPC.Call(rpc.ServerRegister, &req, &res)
}

// OnSyncDisconnect event informs server about lost connection with the Master Server
func OnSyncDisconnect(i *server.Instance, event *event.Event) {
	log.Info("Lost connection with the Master Server! Reconnecting...")
}
