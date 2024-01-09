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
func PreServerEnvRequest(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified) {
		return
	}

	// username := reader.ReadBytes(129)

	pkt := network.NewWriter(CSCPreServerEnvRequest)
	pkt.WriteBytes(make([]byte, 4108)) // auth captcha image data

	session.Send(pkt)
}

// URLToClient Packet which is NFY
func URLToClient(session *Session) {
	conf := session.ServerConfig

	cash_url := conf.CashWeb_URL
	cash_odc_url := conf.CashWeb_Odc_URL
	cash_charge_url := conf.CashWeb_Charge_URL
	guildweb_url := conf.GuildWeb_URL
	sns_url := conf.Sns_URL

	dataLen := len(cash_url) + 4
	dataLen += len(cash_odc_url) + 4
	dataLen += len(cash_charge_url) + 4
	dataLen += len(guildweb_url) + 4
	dataLen += len(sns_url) + 4

	pkt := network.NewWriter(NFYUrlToClient)
	pkt.WriteInt16(dataLen + 2)
	pkt.WriteInt16(dataLen)
	pkt.WriteInt32(len(cash_url))
	pkt.WriteString(cash_url)
	pkt.WriteInt32(len(cash_odc_url))
	pkt.WriteString(cash_odc_url)
	pkt.WriteInt32(len(cash_charge_url))
	pkt.WriteString(cash_charge_url)
	pkt.WriteInt32(len(guildweb_url))
	pkt.WriteString(guildweb_url)
	pkt.WriteInt32(len(sns_url))
	pkt.WriteString(sns_url)
	pkt.WriteInt32(0)
	pkt.WriteByte(0)

	session.Send(pkt)
}

// SystemMessage Packet which is NFY
func SystemMessage(message byte, length uint16) *network.Writer {
	pkt := network.NewWriter(NFYSystemMessage)
	pkt.WriteByte(message)
	pkt.WriteUint16(length)

	return pkt
}

// SystemMessageEx Packet which is NFY
func SystemMessageEx(msg string) *network.Writer {
	pkt := network.NewWriter(NFYSystemMessage)
	pkt.WriteByte(message.Normal)
	pkt.WriteUint16(len(msg) + 2)
	pkt.WriteString("``") // thanks to Iris
	pkt.WriteString(msg)

	return pkt
}

// ServerState Packet which is NFY
func ServerSate(session *Session) *network.Writer {
	// request server list
	req := server.ListReq{}
	res := server.ListRes{}
	session.RPC.Call(rpc.ServerList, &req, &res)

	s := res.List

	pkt := network.NewWriter(NFYServerState)
	pkt.WriteByte(len(s))

	for i := 0; i < len(s); i++ {
		pkt.WriteByte(s[i].Id)
		pkt.WriteByte(s[i].Hot) // 0x10 = HOT! Flag; or bit_set(5)
		pkt.WriteInt32(0x00)
		pkt.WriteByte(len(s[i].List))

		for j := 0; j < len(s[i].List); j++ {
			c := s[i].List[j]
			pkt.WriteByte(c.Id)
			pkt.WriteUint16(c.CurrentUsers)
			pkt.WriteUint16(0x00)
			pkt.WriteUint16(0xFFFF)
			pkt.WriteUint16(0x00)
			pkt.WriteUint16(0x00)
			pkt.WriteUint32(0x00)
			pkt.WriteUint16(0x00)
			pkt.WriteUint16(0x00)
			pkt.WriteUint16(0x00)
			pkt.WriteByte(0x00)
			pkt.WriteByte(0x00)
			pkt.WriteByte(0x00)
			pkt.WriteByte(0xFF)
			pkt.WriteUint16(c.MaxUsers)

			// if session is local, provide local IPs...
			// this helps during development when you have local & remote clients
			// however, here we assume that locally all servers will run on the
			// same IP
			if session.IsLocal() && c.UseLocalIp {
				ip := net.ParseIP(session.GetLocalEndPntIp())[12:16]
				pkt.WriteUint32(binary.LittleEndian.Uint32(ip))
			} else {
				pkt.WriteUint32(c.Ip)
			}

			pkt.WriteUint16(c.Port)
			pkt.WriteUint32(c.Type)
		}
	}

	return pkt
}

// VerifyLinks Packet
func VerifyLinks(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified) {
		return
	}

	authKey := reader.ReadUint32()
	userIdx := reader.ReadUint16()
	channel := reader.ReadByte()
	server := reader.ReadByte()
	magickey := reader.ReadInt32()

	conf := session.ServerConfig

	state := true

	if !conf.IgnoreVersionCheck && magickey != int32(conf.MagicKey) {
		log.Errorf("Invalid MagicKey (Required: %d, detected: %d, id: %d, src: %s",
			conf.MagicKey, magickey, session.Account, session.GetEndPnt(),
		)

		state = false
	}

	req := account.VerifyReq{
		AuthKey:   authKey,
		UserIdx:   userIdx,
		ServerId:  server,
		ChannelId: channel,
		IP:        session.GetIp(),
		DBIdx:     session.Account,
	}
	res := account.VerifyRes{}
	session.RPC.Call(rpc.UserVerify, &req, &res)

	if state {
		state = res.Verified
	}

	pkt := network.NewWriter(CSCVerifyLinks)
	pkt.WriteByte(channel)
	pkt.WriteByte(server)
	pkt.WriteBool(state)

	session.Send(pkt)
}
