package packet

import (
	"github.com/ubis/Freya/share/models/skills"
	"github.com/ubis/Freya/share/network"
)

// QuickLinkSet Packet
func QuickLinkSet(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	slot := reader.ReadUint16()
	skill := reader.ReadUint16()

	state := false
	var err error

	// removing quick link
	if skill == 0xFFFF {
		state, err = session.Character.Links.Remove(slot)
	} else {
		state, err = session.Character.Links.Set(slot, skills.Link{Skill: skill})
	}

	if err != nil {
		session.LogErrorf("An error occurred: %s for character: %d ",
			err.Error(), session.Character.Id)
	}

	pkt := network.NewWriter(CSCQuickLinkSet)
	pkt.WriteBool(state)

	session.Send(pkt)
}

// QuickLinkSwap Packet
func QuickLinkSwap(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	old := reader.ReadUint16()
	new := reader.ReadUint16()

	state, err := session.Character.Links.Swap(old, new)
	if err != nil {
		session.LogErrorf("An error occurred: %s for character: %d ",
			err.Error(), session.Character.Id)
	}

	pkt := network.NewWriter(CSCQuickLinkSwap)
	pkt.WriteBool(state)

	session.Send(pkt)
}
