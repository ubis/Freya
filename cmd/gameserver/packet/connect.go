package packet

import (
	"github.com/ubis/Freya/share/network"
)

// Connect2Svr Packet
func Connect2Svr(session *Session, reader *network.Reader) {
	if !verifyState(session, StateUnknown, reader.Type) {
		return
	}

	session.SetState(StateConnected)

	pkt := network.NewWriter(CSCConnect2Svr)
	pkt.WriteUint32(session.GetSeed())
	pkt.WriteUint32(session.GetAuthKey())
	pkt.WriteUint16(session.GetUserIdx())
	pkt.WriteUint16(session.GetKeyIdx())

	session.Send(pkt)
}
