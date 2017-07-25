package main

import (
    "net"
    "encoding/binary"
    "share/event"
    "share/network"
    "share/rpc"
    "share/models/server"
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
        log.Infof("Client `%s` connected to the GameServer", s.GetEndPnt())
    }
}

// OnClientDisconnect event informs server about disconnected client
func OnClientDisconnect(e event.Event) {
    if s, ok := e.(*network.Session); !ok {
        log.Error("Cannot parse onClientDisconnect event!")
    } else {
        log.Infof("Client `%s` disconnected from the GameServer", s.GetEndPnt())
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
func OnSyncConnect(event event.Event) {
    log.Info("Established connection with the Master Server!")

    // register this server
    var req  = server.RegisterReq{
        server.GAME_SERVER,
        byte(g_ServerConfig.ServerType),
        byte(g_ServerSettings.ServerId),
        byte(g_ServerSettings.ChannelId),
        binary.LittleEndian.Uint32(net.ParseIP(g_ServerConfig.PublicIp)[12:16]),
        uint16(g_ServerConfig.Port),
        uint16(g_NetworkManager.GetOnlineUsers()),
        uint16(g_ServerConfig.MaxUsers),
    }
    var res = server.RegisterRes{}

    g_RPCHandler.Call(rpc.ServerRegister, req, &res)
}

// OnSyncDisconnect event informs server about lost connection with the Master Server
func OnSyncDisconnect(event event.Event) {
    log.Info("Lost connection with the Master Server! Reconnecting...")
}