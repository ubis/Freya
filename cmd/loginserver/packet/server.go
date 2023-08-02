package packet

import (
	"encoding/binary"
	"net"

	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/account"
	"github.com/ubis/Freya/share/models/message"
	"github.com/ubis/Freya/share/models/server"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

// PreServerEnvRequest Packet
func PreServerEnvRequest(session *network.Session, reader *network.Reader) {
	var packet = network.NewWriter(PRE_SERVER_ENV_REQUEST)
	packet.WriteBytes(make([]byte, 4113))
	session.Send(packet)
}

// URLToClient Packet which is NFY
func URLToClient(session *network.Session) {
	var cash_url = g_ServerConfig.CashWeb_URL
	var cash_odc_url = g_ServerConfig.CashWeb_Odc_URL
	var cash_charge_url = g_ServerConfig.CashWeb_Charge_URL
	var guildweb_url = g_ServerConfig.GuildWeb_URL
	var sns_url = g_ServerConfig.Sns_URL

	var dataLen = len(cash_url) + 4
	dataLen += len(cash_odc_url) + 4
	dataLen += len(cash_charge_url) + 4
	dataLen += len(guildweb_url) + 4
	dataLen += len(sns_url) + 4

	var packet = network.NewWriter(URLTOCLIENT)
	packet.WriteInt16(dataLen + 2)
	packet.WriteInt16(dataLen)
	packet.WriteInt32(len(cash_url))
	packet.WriteString(cash_url)
	packet.WriteInt32(len(cash_odc_url))
	packet.WriteString(cash_odc_url)
	packet.WriteInt32(len(cash_charge_url))
	packet.WriteString(cash_charge_url)
	packet.WriteInt32(len(guildweb_url))
	packet.WriteString(guildweb_url)
	packet.WriteInt32(len(sns_url))
	packet.WriteString(sns_url)

	session.Send(packet)
}

// SystemMessg Packet which is NFY
func SystemMessg(message byte, length uint16) *network.Writer {
	var packet = network.NewWriter(SYSTEMMESSG)
	packet.WriteByte(message)
	packet.WriteUint16(length)

	return packet
}

// SystemMessgEx Packet which is NFY
func SystemMessgEx(msg string) *network.Writer {
	var packet = network.NewWriter(SYSTEMMESSG)
	packet.WriteByte(message.Normal)
	packet.WriteUint16(len(msg) + 2)
	packet.WriteString("``") // thanks to Iris
	packet.WriteString(msg)

	return packet
}

// ServerState Packet which is NFY
func ServerSate(session *network.Session) *network.Writer {
	// request server list
	var r = server.ListRes{}
	g_RPCHandler.Call(rpc.ServerList, server.ListReq{}, &r)
	var s = r.List

	var packet = network.NewWriter(SERVERSTATE)
	packet.WriteByte(len(s))

	for i := 0; i < len(s); i++ {
		packet.WriteByte(s[i].Id)
		packet.WriteByte(s[i].Hot) // 0x10 = HOT! Flag; or bit_set(5)
		packet.WriteInt32(0x00)
		packet.WriteByte(len(s[i].List))

		for j := 0; j < len(s[i].List); j++ {
			var c = s[i].List[j]
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

			// if session is local, provide local IPs...
			// this helps during development when you have local & remote clients
			// however, here we assume that locally all servers will run on the
			// same IP
			if session.IsLocal() {
				ip := net.ParseIP(session.GetLocalEndPntIp())[12:16]
				packet.WriteUint32(binary.LittleEndian.Uint32(ip))
			} else {
				packet.WriteUint32(c.Ip)
			}

			packet.WriteUint16(c.Port)
			packet.WriteUint32(c.Type)
		}
	}

	return packet
}

// VerifyLinks
func VerifyLinks(session *network.Session, reader *network.Reader) {
	var timestamp = reader.ReadUint32()
	var count = reader.ReadUint16()
	var channel = reader.ReadByte()
	var server = reader.ReadByte()
	var magickey = reader.ReadInt32()

	if magickey != int32(g_ServerConfig.MagicKey) {
		log.Errorf("Invalid MagicKey (Required: %d, detected: %d, id: %d, src: %s",
			g_ServerConfig.MagicKey, magickey, session.Data.AccountId, session.GetEndPnt(),
		)
		return
	}

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
