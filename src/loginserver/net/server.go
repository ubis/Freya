package net

import (
	"share/models/server"
	"share/network"
	"share/rpc"
)

// PreServerEnvRequest Packet
func (p *Packet) PreServerEnvRequest(s *network.Session, r *network.Reader) {
	packet := network.NewWriter(PreServerEnvRequest)
	packet.WriteBytes(make([]byte, 4113))

	s.Send(packet)
}

// URLToClient Packet which is NFY
func (p *Packet) URLToClient(s *network.Session) {
	cash := p.URL[CashWeb]
	odc := p.URL[CashWebCharge]
	charge := p.URL[CashWebCharge]
	guild := p.URL[GuildWeb]
	sns := p.URL[Sns]

	dataLen := len(cash) + 4
	dataLen += len(odc) + 4
	dataLen += len(charge) + 4
	dataLen += len(guild) + 4
	dataLen += len(sns) + 4

	packet := network.NewWriter(URLToClient)
	packet.WriteInt16(dataLen + 2)
	packet.WriteInt16(dataLen)
	packet.WriteInt32(len(cash))
	packet.WriteString(cash)
	packet.WriteInt32(len(odc))
	packet.WriteString(odc)
	packet.WriteInt32(len(charge))
	packet.WriteString(charge)
	packet.WriteInt32(len(guild))
	packet.WriteString(guild)
	packet.WriteInt32(len(sns))
	packet.WriteString(sns)

	s.Send(packet)
}

// SystemMessg Packet which is NFY
func (p *Packet) SystemMessg(message byte, length uint16) *network.Writer {
	packet := network.NewWriter(SystemMessg)
	packet.WriteByte(message)
	packet.WriteUint16(length)

	return packet
}

// ServerState Packet which is NFY
func (p *Packet) ServerState() *network.Writer {
	// request server list
	req := server.ListReq{}
	rsp := server.ListRes{}
	p.RPC.Call(rpc.ServerList, req, &rsp)
	s := rsp.List

	packet := network.NewWriter(ServerState)
	packet.WriteByte(len(s))

	for i := 0; i < len(s); i++ {
		packet.WriteByte(s[i].Id)
		packet.WriteByte(s[i].Hot) // 0x10 = HOT! Flag; or bit_set(5)
		packet.WriteInt32(0x00)
		packet.WriteByte(len(s[i].List))

		for j := 0; j < len(s[i].List); j++ {
			c := s[i].List[j]
			packet.WriteByte(c.Id)
			packet.WriteUint16(c.CurrentUsers)
			packet.WriteUint16(0x00)
			packet.WriteUint16(0xFFFF)
			packet.WriteUint16(0x00)
			packet.WriteUint16(0x00)
			packet.WriteUint32(0x00)
			packet.WriteUint16(0x00)
			packet.WriteUint16(0x00)
			packet.WriteUint16(0x00)
			packet.WriteByte(0x00)
			packet.WriteByte(0x00)
			packet.WriteByte(0x00)
			packet.WriteByte(0xFF)
			packet.WriteUint16(c.MaxUsers)
			packet.WriteUint32(c.Ip)
			packet.WriteUint16(c.Port)
			packet.WriteUint32(c.Type)
		}
	}

	return packet
}
