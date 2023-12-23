package packet

import (
	"encoding/binary"
	nnet "net"
	"time"

	"github.com/ubis/Freya/share/models/account"
	"github.com/ubis/Freya/share/models/server"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

// GetServerTime Packet
func GetServerTime(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	var now = time.Now()
	var _, z = time.Now().Zone()

	z = z / 60 // to hours
	z = z * -1 // add reverse sign

	pkt := network.NewWriter(CSCGetServerTime)
	pkt.WriteInt64(now.Unix()) // utc time
	pkt.WriteInt16(z)          // timezone

	session.Send(pkt)
}

// ServerEnv Packet
func ServerEnv(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	pkt := network.NewWriter(CSCServerEnv)
	pkt.WriteUint16(0x00BE)      // MaxLevel
	pkt.WriteByte(0x01)          // UseDummy
	pkt.WriteByte(0x01)          // Allow CashShop
	pkt.WriteByte(0x00)          // Allow NetCafePoint
	pkt.WriteUint16(0x0A)        // MaxRank
	pkt.WriteUint16(0x1E)        // Limit Loud Character Lv
	pkt.WriteUint16(0x04)        // Limit Loud Mastery Lv
	pkt.WriteInt64(0x0BA43B7400) // Limit Inventory Alz Save
	pkt.WriteInt64(0x0BA43B7400) // Limit Warehouse Alz Save
	pkt.WriteInt64(0x0BA43B7400) // Limit Trade Alz
	pkt.WriteByte(0x00)          // Allow Duplicated PCBang Premium
	pkt.WriteByte(0x00)          // Allow GuildBoard
	pkt.WriteByte(0x00)          // PCBang Premium Prior Type
	pkt.WriteInt32(0x00)         // Use Trade Channel Restriction
	pkt.WriteInt32(0x01)         // Use AgentShop
	pkt.WriteInt16(0x01)         // Use Lord BroadCast CoolTime Sec
	pkt.WriteByte(0x10)          // Dummy Limit
	pkt.WriteUint16(0x00)        // AgentShop Restriction Lv
	pkt.WriteUint16(0x00)        // PersonalShop Restriction Lv
	pkt.WriteByte(0x01)          // Use TPoint
	pkt.WriteByte(0x01)          // Use Guild Expansion
	pkt.WriteByte(0x00)          // Ignore Party Invite Distance
	pkt.WriteByte(0x01)          // Limited BroadCast By Lord
	pkt.WriteByte(0x00)          // Limit Normal Chat Lv
	pkt.WriteByte(0x00)          // Limit Trade Chat Lv
	pkt.WriteInt32(0x64)         // Max DP Limit
	pkt.WriteInt32(0x00)         // unk1
	pkt.WriteInt16(0x07)         // unk2

	session.Send(pkt)
}

// VerifyLinks Packet
func VerifyLinks(session *Session, reader *network.Reader) {
	// if !verifyState(session, StateInGame, reader.Type) {
	// 	return
	// }

	timestamp := reader.ReadUint32()
	count := reader.ReadUint16()
	channel := reader.ReadByte()
	server := reader.ReadByte()

	req := account.VerifyReq{
		AuthKey:   timestamp,
		UserIdx:   count,
		ServerId:  server,
		ChannelId: channel,
		IP:        session.GetIp(),
		DBIdx:     session.Account,
	}
	res := account.VerifyRes{}
	session.RPC.Call(rpc.UserVerify, &req, &res)

	pkt := network.NewWriter(CSCVerifyLinks)
	pkt.WriteByte(channel)
	pkt.WriteByte(server)
	pkt.WriteBool(res.Verified)

	session.Send(pkt)
}

// SystemMessage Packet which is NFY
func SystemMessage(message byte, length uint16) *network.Writer {
	pkt := network.NewWriter(NFYSystemMessage)
	pkt.WriteByte(message)
	pkt.WriteUint16(length)

	return pkt
}

// BackToCharacterLobby Packet
func BackToCharacterLobby(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	pkt := network.NewWriter(CSCBackToCharacterLobby)
	pkt.WriteByte(1)

	session.Send(pkt)
}

// ChannelList Packet
func ChannelList(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	// request server list
	req := server.ListReq{}
	res := server.ListRes{}
	session.RPC.Call(rpc.ServerList, &req, &res)

	var server *server.ServerItem

	for _, v := range res.List {
		if v.Id != byte(session.ServerInstance.ServerId) {
			continue
		}

		server = &v
		break
	}

	pkt := network.NewWriter(CSCChannelList)

	if server == nil {
		pkt.WriteByte(0)
		session.Send(pkt)
		return
	}

	pkt.WriteByte(len(server.List))

	for _, v := range server.List {
		pkt.WriteByte(v.Id)
		pkt.WriteUint16(v.CurrentUsers)
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
		pkt.WriteUint16(v.MaxUsers)

		// if session is local, provide local IPs...
		// this helps during development when you have local & remote clients
		// however, here we assume that locally all servers will run on the
		// same IP
		if session.IsLocal() && v.UseLocalIp {
			ip := nnet.ParseIP(session.GetLocalEndPntIp())[12:16]
			pkt.WriteUint32(binary.LittleEndian.Uint32(ip))
		} else {
			pkt.WriteUint32(v.Ip)
		}

		pkt.WriteUint16(v.Port)
		pkt.WriteUint32(v.Type)
	}

	session.Send(pkt)
}

// ChannelChange Packet
func ChannelChange(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	_ = reader.ReadByte() // channel id

	pkt := network.NewWriter(CSCChannelChange)
	pkt.WriteInt32(1)

	session.Send(pkt)
}
