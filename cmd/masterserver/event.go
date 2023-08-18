package main

import (
	"github.com/ubis/Freya/share/event"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/rpc"
)

// Registers RPC Server events
func RegisterEvents() {
	log.Info("Registering events...")
	event.Register(event.SyncConnectEvent, event.Handler(OnSyncConnect))
	event.Register(event.SyncDisconnectEvent, event.Handler(OnSyncDisconnect))
}

// OnSyncConnect event informs server about new connection
func OnSyncConnect(event *event.Event) {
	rawClient, ok := event.Get()
	if !ok {
		return
	}

	if c, ok := rawClient.(*rpc.Client); !ok {
		log.Error("Cannot parse onSyncConnect event!")
	} else {
		log.Infof("Client %s connected to the Master Server", c.GetEndPnt())
	}
}

// OnSyncDisconnect event informs server about connection close
func OnSyncDisconnect(event *event.Event) {
	rawClient, ok := event.Get()
	if !ok {
		return
	}

	if c, ok := rawClient.(*rpc.Client); !ok {
		log.Error("Cannot parse onSyncDisconnect event!")
	} else {
		g_ServerManager.RemoveServer(c.GetEndPnt())
		log.Infof("Client %s disconnected from the Master Server", c.GetEndPnt())
	}
}
