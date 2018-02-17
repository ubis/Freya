package net

import (
	"loginserver/rsa"
	"share/log"
	"share/network/packet"
	"share/rpc"
)

// Packet structure
type Packet struct {
	packet.List
	rsa *rsa.Encryption

	RPC *rpc.Client

	Version  int
	MagicKey int
	URL      []string
}

// preInit makes required steps before registering packets
func (p *Packet) preInit() {
	// init RSA encryption
	p.rsa = &rsa.Encryption{}
	p.rsa.Init()
}

// Register network packets
func (p *Packet) Register() {
	// do pre-init work
	p.preInit()

	// register packets
	log.Info("Registering packets...")

	p.Add(Connect2Svr, "Connect2Svr", p.Connect2Svr)
	p.Add(VerifyLinks, "VerifyLinks", p.VerifyLinks)
	p.Add(AuthAccount, "AuthAccount", p.AuthAccount)
	p.Add(FDisconnect, "FDisconnect", p.FDisconnect)
	p.Add(SystemMessg, "SystemMessg", nil)
	p.Add(ServerState, "ServerState", nil)
	p.Add(CheckVersion, "CheckVersion", p.CheckVersion)
	p.Add(URLToClient, "URLToClient", nil)
	p.Add(PublicKey, "PublicKey", p.PublicKey)
	p.Add(PreServerEnvRequest, "PreServerEnvRequest", p.PreServerEnvRequest)
}
