package main

import (
	"encoding/binary"
	"net"

	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/cmd/gameserver/game"
	"github.com/ubis/Freya/share/event"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/server"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

// Registers server events
func RegisterEvents(wm *game.WorldManager) {
	log.Info("Registering events...")

	event.Register(event.ClientConnectEvent, event.Handler(func(e *event.Event) {
		OnClientConnect(e, wm)
	}))

	event.Register(event.ClientDisconnectEvent, event.Handler(OnClientDisconnect))
	event.Register(event.PacketReceiveEvent, event.Handler(OnPacketReceive))
	event.Register(event.PacketSendEvent, event.Handler(OnPacketSend))
	event.Register(event.SyncConnectEvent, event.Handler(OnSyncConnect))
	event.Register(event.SyncDisconnectEvent, event.Handler(OnSyncDisconnect))
}

// OnClientConnect event informs server about new connected client
func OnClientConnect(e *event.Event, wm *game.WorldManager) {
	rawSession, ok := e.Get()
	if !ok {
		return
	}

	s, ok := rawSession.(*network.Session)
	if !ok {
		return
	}

	s.DataEx = &context.Context{WorldManager: wm}
	log.Infof("Client `%s` connected to the GameServer", s.GetEndPnt())
}

// OnClientDisconnect event informs server about disconnected client
func OnClientDisconnect(e *event.Event) {
	rawSession, ok := e.Get()
	if !ok {
		return
	}

	s, ok := rawSession.(*network.Session)
	if !ok {
		return
	}

	// in case client was in the world, notify other players
	world := context.GetWorld(s)
	if world != nil {
		world.ExitWorld(s, server.DelUserLogout)
	}

	log.Infof("Client `%s` disconnected from the GameServer", s.GetEndPnt())
}

// OnPacketReceive event informs server about received packet
func OnPacketReceive(e *event.Event) {
	rawPacket, ok := e.Get()
	if !ok {
		return
	}

	a, ok := rawPacket.(*network.PacketArgs)
	if !ok {
		return
	}

	log.Debugf("Received Packet `%s` (Len: %d, type: %d, src: %s)",
		g_PacketHandler.Name(a.Type), a.Length, a.Type, a.Session.GetEndPnt(),
	)

	// let it handle
	g_PacketHandler.Handle(a)
}

// OnPacketSend event informs server about sent packet
func OnPacketSend(e *event.Event) {
	rawPacket, ok := e.Get()
	if !ok {
		return
	}

	a, ok := rawPacket.(*network.PacketArgs)
	if !ok {
		return
	}

	log.Debugf("Sent Packet `%s` (Len: %d, type: %d, src: %s)",
		g_PacketHandler.Name(a.Type), a.Length, a.Type, a.Session.GetEndPnt(),
	)
}

// OnSyncConnect event informs server about successful connection with the Master Server
func OnSyncConnect(event *event.Event) {
	log.Info("Established connection with the Master Server!")

	// register this server
	req := server.RegisterReq{
		Type:         server.GAME_SERVER,
		ServerType:   byte(g_ServerConfig.ServerType),
		ServerId:     byte(g_ServerSettings.ServerId),
		ChannelId:    byte(g_ServerSettings.ChannelId),
		PublicIp:     binary.LittleEndian.Uint32(net.ParseIP(g_ServerConfig.PublicIp)[12:16]),
		PublicPort:   uint16(g_ServerConfig.Port),
		CurrentUsers: uint16(g_NetworkManager.GetOnlineUsers()),
		MaxUsers:     uint16(g_ServerConfig.MaxUsers),
	}
	res := server.RegisterRes{}

	g_RPCHandler.Call(rpc.ServerRegister, &req, &res)
}

// OnSyncDisconnect event informs server about lost connection with the Master Server
func OnSyncDisconnect(event *event.Event) {
	log.Info("Lost connection with the Master Server! Reconnecting...")
}
