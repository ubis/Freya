package packet

import (
	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/cmd/gameserver/net"
	"github.com/ubis/Freya/share/network"
)

func NewMobsList(mobs []context.MobHandler) *network.Writer {
	pkt := network.NewWriter(net.NFY_NEWMOBSLIST)

	pkt.WriteByte(len(mobs)) // count

	for _, v := range mobs {
		id := v.GetId()
		species := v.GetSpecies()
		pos := v.GetPosition()

		pkt.WriteInt32(id)
		pkt.WriteInt16(pos.InitialX)
		pkt.WriteInt16(pos.InitialY)
		pkt.WriteInt16(pos.FinalX)
		pkt.WriteInt16(pos.FinalY)
		pkt.WriteInt16(species)
		pkt.WriteInt32(0xF3) //max hp
		pkt.WriteInt32(0xF3) // current hp
		pkt.WriteByte(0)
		pkt.WriteInt32(9) // level
		pkt.WriteInt32(0)
	}

	return pkt
}

func MobMoveBegin(mob context.MobHandler) *network.Writer {
	pkt := network.NewWriter(net.NFY_MOBSMOVEBGN)

	id := mob.GetId()
	pos := mob.GetPosition()

	pkt.WriteInt32(id)
	pkt.WriteUint32(pos.MoveBegin)
	pkt.WriteInt16(pos.InitialX)
	pkt.WriteInt16(pos.InitialY)
	pkt.WriteInt16(pos.FinalX)
	pkt.WriteInt16(pos.FinalY)

	return pkt
}

func MobMoveEnd(mob context.MobHandler) *network.Writer {
	pkt := network.NewWriter(net.NFY_MOBSMOVEEND)

	id := mob.GetId()
	pos := mob.GetPosition()

	pkt.WriteInt32(id)
	pkt.WriteInt16(pos.FinalX)
	pkt.WriteInt16(pos.FinalY)

	return pkt
}
