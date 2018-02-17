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

	packet := network.NewWriter(Connect2Svr)
	packet.WriteUint32(s.Encryption.Key.Seed2nd)
	packet.WriteUint32(s.AuthKey)
	packet.WriteUint16(s.UserIdx)
	packet.WriteUint16(s.Encryption.RecvXorKeyIdx)

	s.Send(packet)
}

// CheckVersion Packet
func (p *Packet) CheckVersion(s *network.Session, r *network.Reader) {
	v1 := r.ReadInt32()

	verified := &s.Data.Verified
	*verified = true

	if v1 != int32(p.Version) {
		log.Warningf("Version mismatch (client: %d, server: %d) %s",
			v1, p.Version, s.Info())
		*verified = false
	}

	packet := network.NewWriter(CheckVersion)
	packet.WriteInt32(p.Version)
	packet.WriteInt32(0x00) // debug
	packet.WriteInt32(0x00) // reserved
	packet.WriteInt32(0x00) // reserved

	s.Send(packet)
}

// FDisconnect Packet
func (p *Packet) FDisconnect(s *network.Session, r *network.Reader) {
	idx := r.ReadInt32()

	packet := network.NewWriter(FDisconnect)
	if idx != s.Data.AccountId {
		// wooops invalid account id
		packet.WriteByte(0x00) // failed
		s.Send(packet)
		return
	}

	req := account.OnlineReq{Account: idx, Kick: true}
	res := account.OnlineRes{}
	p.RPC.Call(rpc.ForceDisconnect, req, &res)

	packet.WriteByte(res.Result)

	s.Send(packet)
}

// VerifyLinks Packet
func (p *Packet) VerifyLinks(s *network.Session, r *network.Reader) {
	timestamp := r.ReadUint32()
	count := r.ReadUint16()
	channel := r.ReadByte()
	server := r.ReadByte()
	magickey := r.ReadInt32()

	if magickey != int32(p.MagicKey) {
		log.Errorf("Invalid MagicKey (client: %d, server: %d) %s",
			magickey, p.MagicKey, s.Info())
		return
	}

	send := account.VerifyReq{
		AuthKey:   timestamp,
		UserIdx:   count,
		ServerId:  server,
		ChannelId: channel,
		IP:        s.GetIp(),
		DBIdx:     s.Data.AccountId,
	}
	recv := account.VerifyRes{}
	p.RPC.Call(rpc.UserVerify, send, &recv)

	packet := network.NewWriter(VerifyLinks)
	packet.WriteByte(channel)
	packet.WriteByte(server)

	if recv.Verified {
		packet.WriteByte(0x01)
	} else {
		packet.WriteByte(0x00)
	}

	s.Send(packet)
}
