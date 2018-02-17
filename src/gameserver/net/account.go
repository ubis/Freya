package net

import (
	"bytes"
	"share/models/account"
	"share/network"
	"share/rpc"
)

// ChargeInfo Packet
func (p *Packet) ChargeInfo(session *network.Session, reader *network.Reader) {
	var packet = network.NewWriter(ChargeInfo)
	packet.WriteInt32(0x00)
	packet.WriteInt32(0x00)  // service kind
	packet.WriteUint32(0x00) // service expire

	session.Send(packet)
}

// CheckUserPrivacyData Packet
func (p *Packet) CheckUserPrivacyData(session *network.Session,
	reader *network.Reader) {
	// skip 4 bytes
	reader.ReadInt32()

	var passwd = string(bytes.Trim(reader.ReadBytes(32), "\x00"))

	var req = account.AuthCheckReq{session.Data.AccountId, passwd}
	var res = account.AuthCheckRes{}
	p.RPC.Call(rpc.PasswdCheck, req, &res)

	var packet = network.NewWriter(CheckUserPrivacyData)

	if res.Result {
		// password verified
		packet.WriteByte(0x01)
		session.Data.CharVerified = true
	} else {
		packet.WriteByte(0x00)
	}

	session.Send(packet)
}
