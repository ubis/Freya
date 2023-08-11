package packet

import (
	"bytes"

	"github.com/ubis/Freya/cmd/gameserver/net"
	"github.com/ubis/Freya/share/models/account"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

// ChargeInfo Packet
func ChargeInfo(session *network.Session, reader *network.Reader) {
	var packet = network.NewWriter(net.CHARGEINFO)
	packet.WriteInt32(0x00)
	packet.WriteInt32(0x00)  // service kind
	packet.WriteUint32(0x00) // service expire

	session.Send(packet)
}

// CheckUserPrivacyData Packet
func CheckUserPrivacyData(session *network.Session, reader *network.Reader) {
	// skip 4 bytes
	reader.ReadInt32()

	var passwd = string(bytes.Trim(reader.ReadBytes(32), "\x00"))

	var req = account.AuthCheckReq{session.Data.AccountId, passwd}
	var res = account.AuthCheckRes{}
	g_RPCHandler.Call(rpc.PasswdCheck, req, &res)

	var packet = network.NewWriter(net.CHECK_USR_PDATA)

	if res.Result {
		// password verified
		packet.WriteByte(0x01)
		session.Data.CharVerified = true
	} else {
		packet.WriteByte(0x00)
	}

	session.Send(packet)
}
