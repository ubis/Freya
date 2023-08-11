package notify

import (
	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/cmd/gameserver/net"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/server"
	"github.com/ubis/Freya/share/network"
)

func fillPlayerInfo(pkt *network.Writer, session *network.Session) {
	ctx, err := context.Parse(session)
	if err != nil {
		log.Error(err.Error())
		return
	}

	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()

	c := ctx.Char

	if c == nil {
		// client is not ready
		return
	}

	pkt.WriteUint32(c.Id)
	pkt.WriteUint32(session.UserIdx)
	pkt.WriteUint32(c.Level)
	pkt.WriteInt32(0x01C2)    // might be dwMoveBgnTime
	pkt.WriteUint16(c.BeginX) // start
	pkt.WriteUint16(c.BeginY)
	pkt.WriteUint16(c.EndX) // end
	pkt.WriteUint16(c.EndY)
	pkt.WriteByte(0)
	pkt.WriteInt32(0)
	pkt.WriteInt16(0)
	pkt.WriteInt32(c.Style.Get())
	pkt.WriteByte(0) // animation id aka "live style"
	pkt.WriteInt16(0)

	eq, eqlen := c.Equipment.SerializeEx()
	pkt.WriteInt16(eqlen)
	pkt.WriteInt16(0x00)

	for i := 0; i < 21; i++ {
		pkt.WriteByte(0)
	}

	pkt.WriteByte(len(c.Name) + 1)
	pkt.WriteString(c.Name)
	pkt.WriteByte(0) // guild name len
	// pkt.WriteString("guild name")

	pkt.WriteBytes(eq)
}

func NewUserSingle(session *network.Session) *network.Writer {
	pkt := network.NewWriter(net.NEWUSERLIST)
	pkt.WriteUint16(1) // player num

	fillPlayerInfo(pkt, session)

	return pkt
}

func NewUserList(players map[uint16]*network.Session) *network.Writer {
	online := len(players)

	pkt := network.NewWriter(net.NEWUSERLIST)
	pkt.WriteUint16(online)

	for _, v := range players {
		fillPlayerInfo(pkt, v)
	}

	return pkt
}

// DelUserList to all already connected players
func DelUserList(session *network.Session, reason server.DelUserType) *network.Writer {
	charId, err := context.GetCharId(session)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	pkt := network.NewWriter(net.DELUSERLIST)
	pkt.WriteUint32(charId)
	pkt.WriteByte(byte(reason)) // type

	/* types:
	 * dead = 0x10
	 * warp = 0x11
	 * logout = 0x12
	 * retn = 0x13
	 * dissapear = 0x14
	 * nfsdead = 0x15
	 */

	return pkt
}
