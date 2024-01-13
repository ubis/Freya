package packet

import (
	"github.com/ubis/Freya/share/models/inventory"
	"github.com/ubis/Freya/share/network"
)

type StorageType int

const (
	Inventory StorageType = iota
	Equipment
)

// StorageExchangeMove Packet
func StorageExchangeMove(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	src := StorageType(reader.ReadInt32())
	srcSlot := uint16(reader.ReadInt32())
	dst := StorageType(reader.ReadInt32())
	dstSlot := uint16(reader.ReadInt32())

	inv := session.Character.Inventory
	eq := &session.Character.Equipment

	state := false
	var err error

	switch src {
	case Inventory:
		switch dst {
		case Inventory:
			state, err = inv.Move(srcSlot, dstSlot)
		case Equipment:
			state, err = eq.EquipItem(srcSlot, dstSlot, inv)
			if !state {
				break
			}

			item := eq.Get(dstSlot)

			pkt := network.NewWriter(NFYItemEquip)
			pkt.WriteInt32(session.Character.Id)
			pkt.WriteInt32(item.Kind)
			pkt.WriteInt16(item.Slot)
			pkt.WriteInt32(0)
			pkt.WriteByte(0)

			session.Broadcast(pkt)
		}
	case Equipment:
		switch dst {
		case Inventory:
			state, err = eq.UnEquipItem(srcSlot, dstSlot, inv)
			if !state {
				break
			}

			pkt := network.NewWriter(NFYItemUnEquip)
			pkt.WriteInt32(session.Character.Id)
			pkt.WriteInt16(srcSlot)

			session.Broadcast(pkt)

			// with one-handed dual weapons we need to move from left hand to
			// the right, if right-hand weapon was removed
			// todo: check for dual-handed weapons and ignore
			if srcSlot != inventory.RightHand {
				break
			}

			// switch weapon
			eq.MoveItem(inventory.LeftHand, inventory.RightHand)
		case Equipment:
			state, err = eq.MoveItem(srcSlot, dstSlot)
			if !state {
				break
			}

			item := eq.Get(dstSlot)

			pkt := network.NewWriter(NFYItemUnEquip)
			pkt.WriteInt32(session.Character.Id)
			pkt.WriteInt16(srcSlot)
			session.Broadcast(pkt)

			pkt = network.NewWriter(NFYItemEquip)
			pkt.WriteInt32(session.Character.Id)
			pkt.WriteInt32(item.Kind)
			pkt.WriteInt16(item.Slot)
			pkt.WriteInt32(0)
			pkt.WriteByte(0)
			session.Broadcast(pkt)
		}
	}

	if err != nil {
		session.LogErrorf("An error occurred: %s for character: %d ",
			err.Error(), session.Character.Id)
	}

	pkt := network.NewWriter(CSCStorageExchangeMove)
	pkt.WriteBool(state)
	pkt.WriteInt32(0)

	session.Send(pkt)
}

// StorageItemSwap
func StorageItemSwap(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	src := StorageType(reader.ReadInt32())
	srcSlot := uint16(reader.ReadInt32())
	dst := StorageType(reader.ReadInt32())
	dstSlot := uint16(reader.ReadInt32())

	state := false
	var err error

	inv := session.Character.Inventory
	eq := &session.Character.Equipment

	switch src {
	case Inventory:
		switch dst {
		case Inventory:
			state, err = inv.Swap(srcSlot, dstSlot)
		case Equipment:
			state, err = eq.SwapEquipItem(srcSlot, dstSlot, inv)
			if !state {
				break
			}

			item := eq.Get(dstSlot)

			pkt := network.NewWriter(NFYItemEquip)
			pkt.WriteInt32(session.Character.Id)
			pkt.WriteInt32(item.Kind)
			pkt.WriteInt16(item.Slot)
			pkt.WriteInt32(0)
			pkt.WriteByte(0)
			session.Broadcast(pkt)
		}
	case Equipment:
		switch dst {
		case Inventory:
			/* do nothing */
			return
		case Equipment:
			state, err = eq.Swap(srcSlot, dstSlot)
		}
	}

	if err != nil {
		session.LogErrorf("An error occurred: %s for character: %d ",
			err.Error(), session.Character.Id)
	}

	pkt := network.NewWriter(CSCStorageItemSwap)
	pkt.WriteBool(state)

	session.Send(pkt)
}

// StorageItemDrop Packet
func StorageItemDrop(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	_ = reader.ReadInt32() // unk
	slot := reader.ReadUint16()

	world := GetCurrentWorld(session)
	if world == nil {
		session.LogErrorf("Unable to find current world for character: %d ",
			session.Character.Id)
		return
	}

	item := session.Character.Inventory.Get(slot)

	state, err := session.Character.Inventory.Remove(slot)
	if err != nil {
		session.LogErrorf("An error occurred: %s for character: %d ",
			err.Error(), session.Character.Id)
	}

	pkt := network.NewWriter(CSCStorageItemDrop)
	pkt.WriteBool(state)

	session.Send(pkt)

	if state {
		x, y := session.Character.GetPosition()
		world.DropItem(&item, session.Character.Id, int(x), int(y))
	}
}

// AccessoryEquip
func AccessoryEquip(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	type AccessoryType int

	const (
		Earring AccessoryType = iota + 1
		Bracelet
		Ring
	)

	slot := uint16(reader.ReadUint32())
	reader.ReadInt32() // seem to be identical to slot
	accyType := AccessoryType(reader.ReadInt32())

	// we need to un-equip last type
	// then switch 1st to 2nd, 2nd to 3rd, 3rd to 4th
	// and equip the one from inventory in the 1st type
	var slots []inventory.EquipmentType

	switch accyType {
	case Earring:
		slots = []inventory.EquipmentType{
			inventory.LeftEarring, inventory.RightEarring,
		}
	case Bracelet:
		slots = []inventory.EquipmentType{
			inventory.LeftBracelet, inventory.RightBracelet,
		}

	case Ring:
		slots = []inventory.EquipmentType{
			inventory.Ring1, inventory.Ring2, inventory.Ring3, inventory.Ring4,
		}
	default:
		session.LogErrorf("Unknown accessory type: %d for character: %d ",
			accyType, session.Character.Id)
		return
	}

	inv := session.Character.Inventory
	eq := &session.Character.Equipment

	state, err := eq.EquipAccessory(slot, slots, inv)
	if err != nil {
		session.LogErrorf("An error occurred: %s for character: %d ",
			err.Error(), session.Character.Id)
	}

	pkt := network.NewWriter(CSCAccessoryEquip)
	pkt.WriteBool(state)

	session.Send(pkt)
}
