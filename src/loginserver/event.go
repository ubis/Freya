package main

import (
	"share/event"
	"share/network"
)

// Registers server events
func RegisterEvents() {
	event.Register(event.ClientConnectEvent, event.Handler(OnClientConnect))
	event.Register(event.ClientDisconnectEvent, event.Handler(OnClientDisconnect))
	event.Register(event.PacketReceivedEvent, event.Handler(OnPacketReceived))
	event.Register(event.PacketSentEvent, event.Handler(OnPacketSent))
}

/*
	OnClientConnect event, informs server about connected client
	@param 	event 	Event interface, which is later parsed into Session struct
 */
func OnClientConnect(event event.Event) {
	var session, err = event.(*network.Session)
	if err != true {
		log.Error("Couldn't parse onClientConnect event!")
		return
	}

	log.Infof("Client %s connected to the LoginServer", session.GetEndPnt())
}

/*
	OnClientDisconnect event, informs server about disconnected client
	@param 	event 	Event interface, which is later parsed into Session struct
 */
func OnClientDisconnect(event event.Event) {
	var session, err = event.(*network.Session)
	if err != true {
		log.Error("Couldn't parse onClientDisconnect event!")
		return
	}

	log.Infof("Client %s disconnected from the LoginServer", session.GetEndPnt())
}

func OnPacketReceived(event event.Event) {
	log.Info("Received packet")
}

func OnPacketSent(event event.Event) {
	log.Info("Sent packet")
}