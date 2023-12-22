package packet

import (
	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/skills"
	"github.com/ubis/Freya/share/network"
)

func QuickLinkSet(session *network.Session, reader *network.Reader) {
	slot := reader.ReadUint16()
	skill := reader.ReadUint16()

	ctx, err := context.Parse(session)
	if err != nil {
		log.Error(err.Error())
		return
	}

	ok, err := false, nil

	// removing quick link
	if skill == 0xFFFF {
		ok, err = ctx.Char.Links.Remove(slot)
	} else {
		ok, err = ctx.Char.Links.Set(slot, skills.Link{Skill: skill})
	}

	if err != nil {
		log.Error(err.Error())
	}

	pkt := network.NewWriter(QUICKLINKSET)
	pkt.WriteBool(ok)

	session.Send(pkt)
}

func QuickLinkSwap(session *network.Session, reader *network.Reader) {
	old := reader.ReadUint16()
	new := reader.ReadUint16()

	ctx, err := context.Parse(session)
	if err != nil {
		log.Error(err.Error())
		return
	}

	ok, err := ctx.Char.Links.Swap(old, new)
	if err != nil {
		log.Error(err.Error())
	}

	pkt := network.NewWriter(QUICKLINKSWAP)
	pkt.WriteBool(ok)

	session.Send(pkt)
}
