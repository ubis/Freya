package packet

import (
	"share/models/account"
	"share/network"
	"share/rpc"
	"time"
)

// GetSvrTime Packet
func GetSvrTime(session *network.Session, reader *network.Reader) {
	var now = time.Now()
	var _, z = time.Now().Zone()

	z = z / 60 // to hours
	z = z * -1 // add reverse sign

	var packet = network.NewWriter(GETSVRTIME)
	packet.WriteInt64(now.Unix()) // utc time
	packet.WriteInt16(z)          // timezone

	session.Send(packet)
}

// ServerEnv Packet
func ServerEnv(session *network.Session, reader *network.Reader) {
	var packet = network.NewWriter(SERVERENV)
	packet.WriteUint16(0x00BE)      // MaxLevel
	packet.WriteByte(0x01)          // UseDummy
	packet.WriteByte(0x01)          // Allow CashShop
	packet.WriteByte(0x00)          // Allow NetCafePoint
	packet.WriteUint16(0x0A)        // MaxRank
	packet.WriteUint16(0x1E)        // Limit Loud Character Lv
	packet.WriteUint16(0x04)        // Limit Loud Mastery Lv
	packet.WriteInt64(0x0BA43B7400) // Limit Inventory Alz Save
	packet.WriteInt64(0x0BA43B7400) // Limit Warehouse Alz Save
	packet.WriteInt64(0x0BA43B7400) // Limit Trade Alz
	packet.WriteByte(0x00)          // Allow Duplicated PCBang Premium
	packet.WriteByte(0x00)          // Allow GuildBoard
	packet.WriteByte(0x00)          // PCBang Premium Prior Type
	packet.WriteInt32(0x00)         // Use Trade Channel Restriction
	packet.WriteInt32(0x01)         // Use AgentShop
	packet.WriteInt16(0x01)         // Use Lord BroadCast CoolTime Sec
	packet.WriteByte(0x10)          // Dummy Limit
	packet.WriteUint16(0x00)        // AgentShop Restriction Lv
	packet.WriteUint16(0x00)        // PersonalShop Restriction Lv
	packet.WriteByte(0x01)          // Use TPoint
	packet.WriteByte(0x01)          // Use Guild Expansion
	packet.WriteByte(0x00)          // Ignore Party Invite Distance
	packet.WriteByte(0x01)          // Limited BroadCast By Lord
	packet.WriteByte(0x00)          // Limit Normal Chat Lv
	packet.WriteByte(0x00)          // Limit Trade Chat Lv
	packet.WriteInt32(0x64)         // Max DP Limit
	packet.WriteInt32(0x00)         // unk1
	packet.WriteInt16(0x07)         // unk2

	session.Send(packet)
}

// VerifyLinks
func VerifyLinks(session *network.Session, reader *network.Reader) {
	var timestamp = reader.ReadUint32()
	var count = reader.ReadUint16()
	var channel = reader.ReadByte()
	var server = reader.ReadByte()

	var send = account.VerifyReq{
		timestamp, count, server, channel, session.GetIp(), session.Data.AccountId}
	var recv = account.VerifyRes{}
	g_RPCHandler.Call(rpc.UserVerify, send, &recv)

	var packet = network.NewWriter(VERIFYLINKS)
	packet.WriteByte(channel)
	packet.WriteByte(server)

	if recv.Verified {
		packet.WriteByte(0x01)
	} else {
		packet.WriteByte(0x00)
	}

	session.Send(packet)
}

// SystemMessg Packet which is NFY
func SystemMessg(message byte, length uint16) *network.Writer {
	var packet = network.NewWriter(SYSTEMMESSG)
	packet.WriteByte(message)
	packet.WriteUint16(length)

	return packet
}
