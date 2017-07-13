package main

import (
    "share/event"
    "share/rpc"
)

// Registers RPC Server events
func RegisterEvents() {
    log.Info("Registering events...")
    event.Register(event.SyncConnectEvent, event.Handler(OnSyncConnect))
    event.Register(event.SyncDisconnectEvent, event.Handler(OnSyncDisconnect))
}

/*
    OnSyncConnect event, informs server about new connection
    @param  event   Event interface which is later parsed into RPC Client
 */
func OnSyncConnect(event event.Event) {
    var c, err = event.(*rpc.Client)

    if err != true {
        log.Error("Cannot parse onSyncConnect event!")
        return
    }

    log.Infof("Client %s connected to the Master Server", c.GetEndPnt())
}

/*
    OnSyncDisconnect event, informs server about connection close
    @param  event   Event interface which is later parsed into RPC Client
 */
func OnSyncDisconnect(event event.Event) {
    var c, err = event.(*rpc.Client)

    if err != true {
        log.Error("Cannot parse onSyncDisconnect event!")
        return
    }

    log.Infof("Client %s disconnected from the Master Server", c.GetEndPnt())
}