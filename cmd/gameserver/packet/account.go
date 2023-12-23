package packet

import (
	"bytes"

	"github.com/ubis/Freya/share/models/account"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

// ChargeInfo Packet
func ChargeInfo(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	pkt := network.NewWriter(CSCChargeInfo)
	pkt.WriteInt32(0x00)
	pkt.WriteInt32(0x00)  // service kind
	pkt.WriteUint32(0x00) // service expire

	session.Send(pkt)
}

// CheckUserPrivacyData Packet
func CheckUserPrivacyData(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	// skip 4 bytes
	reader.ReadInt32()

	passwd := string(bytes.Trim(reader.ReadBytes(32), "\x00"))

	req := account.AuthCheckReq{Id: session.Account, Password: passwd}
	res := account.AuthCheckRes{}
	session.RPC.Call(rpc.PasswdCheck, &req, &res)

	session.PasswordVerified = res.Result

	pkt := network.NewWriter(CSCCheckUserPrivacyData)
	pkt.WriteByte(res.Result)

	session.Send(pkt)
}
