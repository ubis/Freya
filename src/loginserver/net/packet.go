package net

import (
	"loginserver/rsa"
	"share/log"
	"share/network"
	"share/rpc"
)

// Packet structure
type Packet struct {
	RPC *rpc.Client
	RSA *rsa.RSA

	Version  int
	MagicKey int
	URL      []string
}

// Register network packets
func (p *Packet) Register(m *network.Manager) {
	log.Info("Registering packets...")

	m.RegisterPacket(Connect2Svr, "Connect2Svr", p.Connect2Svr)
	m.RegisterPacket(VerifyLinks, "VerifyLinks", p.VerifyLinks)
	m.RegisterPacket(AuthAccount, "AuthAccount", p.AuthAccount)
	m.RegisterPacket(FDisconnect, "FDisconnect", p.FDisconnect)
	m.RegisterPacket(SystemMessg, "SystemMessg", nil)
	m.RegisterPacket(ServerState, "ServerState", nil)
	m.RegisterPacket(CheckVersion, "CheckVersion", p.CheckVersion)
	m.RegisterPacket(URLToClient, "URLToClient", nil)
	m.RegisterPacket(PublicKey, "PublicKey", p.PublicKey)
	m.RegisterPacket(PreServerEnvRequest, "PreServerEnvRequest",
		p.PreServerEnvRequest)
}
