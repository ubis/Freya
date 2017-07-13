package main

import (
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

/*
    OnClientConnect event, informs server about connected client
    @param  event    Event interface, which is later parsed into Session struct
 */
func OnClientConnect(event event.Event) {
    var session, err = event.(*network.Session)

    if err != true {
        log.Error("Cannot parse onClientConnect event!")
        return
    }

    log.Infof("Client `%s` connected to the LoginServer", session.GetEndPnt())
}

/*
    OnClientDisconnect event, informs server about disconnected client
    @param  event   Event interface, which is later parsed into Session struct
 */
func OnClientDisconnect(event event.Event) {
    var session, err = event.(*network.Session)

    if err != true {
        log.Error("Cannot parse onClientDisconnect event!")
        return
    }

    log.Infof("Client `%s` disconnected from the LoginServer", session.GetEndPnt())
}

/*
    OnPacketReceive event, informs server about received packet
    @param  event   Event interface, which is later parsed into PacketArgs struct
 */
func OnPacketReceive(event event.Event) {
    var args, err = event.(*network.PacketArgs)

    if err != true {
        log.Error("Cannot parse onPacketReceive event!")
        return
    }

    log.Debugf("Received Packet `%s` (Len: %d, Type: %d, Src: %s, UserIdx: %d)",
        g_PacketHandler.Name(args.Type),
        args.Length,
        args.Type,
        args.Session.GetEndPnt(),
        args.Session.UserIdx,
    )

    // let it handle
    g_PacketHandler.Handle(args)
}

/*
    OnPacketSend event, informs server about sent packet
    @param  event   Event interface, which is later parsed into PacketArgs struct,
    however Data field is nil, since we don't need packet's data anymore
 */
func OnPacketSend(event event.Event) {
    var args, err = event.(*network.PacketArgs)

    if err != true {
        log.Error("Cannot parse onPacketSent event!")
        return
    }

    log.Debugf("Sent Packet `%s` (Len: %d, Type: %d, Src: %s, UserIdx: %d)",
        g_PacketHandler.Name(args.Type),
        args.Length,
        args.Type,
        args.Session.GetEndPnt(),
        args.Session.UserIdx,
    )
}

/*
    OnSyncConnect event, informs server about succesfull connection with the Master Server
    @param  event   Event interface which is nil
 */
func OnSyncConnect(event event.Event) {
    log.Info("Established connection with the Master Server!")

    // register this server
    var req  = server.RegRequest{Type: server.LOGIN_SERVER_TYPE}
    var resp = server.RegResponse{}

    g_RPCHandler.Call(rpc.ServerRegister, req, &resp)
}

/*
    OnSyncDisconnect event, informs server about lost connection with the Master Server
    @param  event   Event interface which is nil
 */
func OnSyncDisconnect(event event.Event) {
    log.Info("Lost connection with the Master Server! Reconnecting...")
}