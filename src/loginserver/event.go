package main

import (
    "share/event"
    "share/network"
)

// Registers server events
func RegisterEvents() {
    event.Register(event.ClientConnectEvent, event.Handler(OnClientConnect))
    event.Register(event.ClientDisconnectEvent, event.Handler(OnClientDisconnect))
    event.Register(event.PacketReceiveEvent, event.Handler(OnPacketReceive))
    event.Register(event.PacketSendEvent, event.Handler(OnPacketSend))
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