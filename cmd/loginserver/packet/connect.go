package packet

import (
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/account"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

// Connect2Svr Packet
func Connect2Svr(session *Session, reader *network.Reader) {
	if !verifyState(session, StateUnknown) {
		return
	}

	if err := network.Deserialize(reader, &C2SConnect2Svr{}); err != nil {
		session.SerializationError(reader.Type, err)
		return
	}

	pkt := S2CConnect2Svr{
		S2CHeader: S2CHeader{
			MagicCode: MagicKey,
			Opcode:    CSCConnect2Svr,
		},
		XorSeed:   session.GetSeed(),
		AuthKey:   session.GetAuthKey(),
		UserIdx:   session.GetUserIdx(),
		XorKeyIdx: uint16(session.GetKeyIdx()),
	}

	// mark state as connected and keys exchanged
	session.state = StateConnected

	session.Send(pkt)
}

// CheckVersion Packet
func CheckVersion(session *Session, reader *network.Reader) {
	version1 := reader.ReadInt32()

	conf := session.ServerConfig

	targetVersion := int32(conf.Version)
	if conf.IgnoreVersionCheck {
		targetVersion = version1
	}

	session.state = StateVerified

	if version1 != targetVersion {
		log.Errorf("Client version mismatch (Client: %d, server: %d, src: %s)",
			version1, targetVersion, session.GetEndPnt())

		session.state = StateConnected
	}

	pkt := network.NewWriter(CSCCheckVersion)
	pkt.WriteInt32(targetVersion)
	pkt.WriteInt32(0x00) // debug
	pkt.WriteInt32(0x00) // reserved
	pkt.WriteInt32(0x00) // reserved

	session.Send(pkt)
}

// ForceDisconnect Packet
func ForceDisconnect(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified) {
		return
	}

	idx := reader.ReadInt32()

	packet := network.NewWriter(CSCForceDisconnect)

	if idx != session.Account {
		// wooops invalid account id
		packet.WriteByte(0x00) // failed
		session.Send(packet)
		return
	}

	req := account.OnlineReq{Account: idx, Kick: true}
	res := account.OnlineRes{}
	session.RPC.Call(rpc.ForceDisconnect, &req, &res)

	packet.WriteBool(res.Result)

	session.Send(packet)
}
