package packet

import (
	"time"

	"github.com/ubis/Freya/cmd/gameserver/net"
	"github.com/ubis/Freya/share/network"
)

// Connect2Svr Packet
func Connect2Svr(session *network.Session, reader *network.Reader) {
	session.AuthKey = uint32(time.Now().Unix())

	var packet = network.NewWriter(net.CONNECT2SVR)
	packet.WriteUint32(session.Encryption.Key.Seed2nd)
	packet.WriteUint32(session.AuthKey)
	packet.WriteUint16(session.UserIdx)
	packet.WriteUint16(session.Encryption.RecvXorKeyIdx)

	session.Send(packet)
}
