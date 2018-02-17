package net

import (
	"share/log"
	"share/models/account"
	"share/network"
	"share/rpc"
	"time"
)

// Connect2Svr Packet
func (p *Packet) Connect2Svr(s *network.Session, r *network.Reader) {
	s.AuthKey = uint32(time.Now().Unix())

	var packet = network.NewWriter(Connect2Svr)
	packet.WriteUint32(s.Encryption.Key.Seed2nd)
	packet.WriteUint32(s.AuthKey)
	packet.WriteUint16(s.UserIdx)
	packet.WriteUint16(s.Encryption.RecvXorKeyIdx)

	s.Send(packet)
}

// CheckVersion Packet
func (p *Packet) CheckVersion(s *network.Session, r *network.Reader) {
	var v1 = r.ReadInt32()

	verified := &s.Data.Verified
	*verified = true

	if v1 != int32(p.Version) {
		log.Errorf("Client version mismatch (Client: %d, server: %d, src: %s)",
			v1, p.Version, s.GetEndPnt())
		*verified = false
	}

	var packet = network.NewWriter(CheckVersion)
	packet.WriteInt32(p.Version)
	packet.WriteInt32(0x00) // debug
	packet.WriteInt32(0x00) // reserved
	packet.WriteInt32(0x00) // reserved

	s.Send(packet)
}

// FDisconnect Packet
func (p *Packet) FDisconnect(s *network.Session, r *network.Reader) {
	var idx = r.ReadInt32()

	var packet = network.NewWriter(FDisconnect)
	if idx != s.Data.AccountId {
		// wooops invalid account id
		packet.WriteByte(0x00) // failed
		s.Send(packet)
		return
	}

	var req = account.OnlineReq{Account: idx, Kick: true}
	var res = account.OnlineRes{}
	p.RPC.Call(rpc.ForceDisconnect, req, &res)

	packet.WriteByte(res.Result)

	s.Send(packet)
}

// VerifyLinks Packet
func (p *Packet) VerifyLinks(s *network.Session, r *network.Reader) {
	var timestamp = r.ReadUint32()
	var count = r.ReadUint16()
	var channel = r.ReadByte()
	var server = r.ReadByte()
	var magickey = r.ReadInt32()

	if magickey != int32(p.MagicKey) {
		log.Errorf("Invalid MagicKey (Client: %d, Server: %d, id: %d, src: %s",
			magickey, p.MagicKey, s.Data.AccountId, s.GetEndPnt())
		return
	}

	var send = account.VerifyReq{
		AuthKey:   timestamp,
		UserIdx:   count,
		ServerId:  server,
		ChannelId: channel,
		IP:        s.GetIp(),
		DBIdx:     s.Data.AccountId,
	}
	var recv = account.VerifyRes{}
	p.RPC.Call(rpc.UserVerify, send, &recv)

	var packet = network.NewWriter(VerifyLinks)
	packet.WriteByte(channel)
	packet.WriteByte(server)

	if recv.Verified {
		packet.WriteByte(0x01)
	} else {
		packet.WriteByte(0x00)
	}

	s.Send(packet)
}
