package packet

import (
	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/inventory"
	"github.com/ubis/Freya/share/network"
)

type StorageType int

const (
	Inventory StorageType = iota
	Equipment
)

func notifyStorageExchange(session *network.Session, result bool) {
	packet := network.NewWriter(STORAGE_EXCHANGE_MOVE)
	packet.WriteBool(result)
	packet.WriteInt32(0)

	session.Send(packet)
}

func StorageExchangeMove(session *network.Session, reader *network.Reader) {
	isEquip := reader.ReadUint32() == 1
	deleteSlot := uint16(reader.ReadUint32())
	isInventory := reader.ReadUint32() == 1
	createSlot := uint16(reader.ReadUint32())

	var id int32

	ctx, err := context.Parse(session)
	if err != nil {
		log.Error(err.Error())
		return
	}

	id, err = context.GetCharId(session)
	if err != nil {
		log.Error(err.Error())
		return
	}

	switch {
	case isEquip && !isInventory:
		// from equipment to inventory
		ctx.Mutex.Lock()
		ok, err := ctx.Char.Equipment.UnEquipItem(deleteSlot, createSlot, ctx.Char.Inventory)
		ctx.Mutex.Unlock()

		notifyStorageExchange(session, ok)

		if err != nil {
			log.Error(err.Error())
			return
		}

		pkt := network.NewWriter(NFY_ITEM_UNEQUIP)
		pkt.WriteInt32(id)
		pkt.WriteInt16(deleteSlot)

		ctx.World.BroadcastSessionPacket(session, pkt)

		// with one-handed dual weapons we need to move from left hand to
		// the right, if right-hand weapon was removed
		// todo: check for dual-handed weapons and ignore
		if deleteSlot == inventory.RightHand {
			// switch weapon
			ctx.Mutex.Lock()
			ctx.Char.Equipment.MoveItem(inventory.LeftHand, inventory.RightHand)
			ctx.Mutex.Unlock()
		}
	case isInventory && !isEquip:
		// from inventory to equipment
		ctx.Mutex.Lock()
		ok, err := ctx.Char.Equipment.EquipItem(deleteSlot, createSlot, ctx.Char.Inventory)
		item := ctx.Char.Equipment.Get(createSlot)
		ctx.Mutex.Unlock()

		notifyStorageExchange(session, ok)

		if err != nil {
			log.Error(err.Error())
			return
		}

		pkt := network.NewWriter(NFY_ITEM_EQUIP)
		pkt.WriteInt32(id)
		pkt.WriteInt32(item.Kind)
		pkt.WriteInt16(item.Slot)
		pkt.WriteInt32(0)
		pkt.WriteByte(0)

		ctx.World.BroadcastSessionPacket(session, pkt)
	case isEquip && isInventory:
		// exchanging equipment items? rings? because on weaps it doesn't work
		ctx.Mutex.Lock()
		ok, err := ctx.Char.Equipment.MoveItem(deleteSlot, createSlot)
		item := ctx.Char.Equipment.Get(createSlot)
		ctx.Mutex.Unlock()

		notifyStorageExchange(session, ok)

		if err != nil {
			log.Error(err.Error())
			return
		}

		pkt := network.NewWriter(NFY_ITEM_UNEQUIP)
		pkt.WriteInt32(id)
		pkt.WriteInt16(deleteSlot)
		ctx.World.BroadcastSessionPacket(session, pkt)

		pkt = network.NewWriter(NFY_ITEM_EQUIP)
		pkt.WriteInt32(id)
		pkt.WriteInt32(item.Kind)
		pkt.WriteInt16(item.Slot)
		pkt.WriteInt32(0)
		pkt.WriteByte(0)
		ctx.World.BroadcastSessionPacket(session, pkt)
	case !isEquip && !isInventory:
		// moving item in inventory
		ctx.Mutex.Lock()
		ok, err := ctx.Char.Inventory.Move(deleteSlot, createSlot)
		ctx.Mutex.Unlock()

		notifyStorageExchange(session, ok)

		if err != nil {
			log.Error(err.Error())
			return
		}
	default:
		notifyStorageExchange(session, false)
		return
	}
}

func StorageItemSwap(session *network.Session, reader *network.Reader) {
	src := StorageType(reader.ReadInt32())
	srcSlot := uint16(reader.ReadInt32())
	dst := StorageType(reader.ReadInt32())
	dstSlot := uint16(reader.ReadInt32())

	ctx, err := context.Parse(session)
	if err != nil {
		log.Error(err.Error())
		return
	}

	state, err := false, nil

	ctx.Mutex.Lock()
	inv := ctx.Char.Inventory
	eq := &ctx.Char.Equipment

	switch src {
	case Inventory:
		switch dst {
		case Inventory:
			state, err = inv.Swap(srcSlot, dstSlot)
		case Equipment:
			state, err = eq.SwapEquipItem(srcSlot, dstSlot, inv)
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
	ctx.Mutex.Unlock()

	if err != nil {
		log.Error(err.Error())
	}

	pkt := network.NewWriter(STORAGE_ITEM_SWAP)
	pkt.WriteBool(state)

	session.Send(pkt)
}

func StorageItemDrop(session *network.Session, reader *network.Reader) {
	_ = reader.ReadInt32() // unk
	slot := reader.ReadUint16()

	ctx, err := context.Parse(session)
	if err != nil {
		log.Error(err.Error())
		return
	}

	ctx.Mutex.RLock()
	charId := ctx.Char.Id
	item := ctx.Char.Inventory.Get(slot)
	world := ctx.World
	x := int(ctx.Char.X)
	y := int(ctx.Char.Y)
	ctx.Mutex.RUnlock()

	if world == nil {
		log.Error("Unable to get current world!")
		return
	}

	ctx.Mutex.Lock()
	state, err := ctx.Char.Inventory.Remove(slot)
	ctx.Mutex.Unlock()

	if err != nil {
		log.Error(err.Error())
	}

	pkt := network.NewWriter(STORAGE_ITEM_DROP)
	pkt.WriteBool(state)

	session.Send(pkt)

	if state {
		world.DropItem(&item, charId, x, y)
	}
}

func AccessoryEquip(session *network.Session, reader *network.Reader) {
	type AccessoryType int

	const (
		Earring AccessoryType = iota + 1
		Bracelet
		Ring
	)

	slot := uint16(reader.ReadUint32())
	reader.ReadInt32() // seem to be identical to slot
	accyType := AccessoryType(reader.ReadInt32())

	ctx, err := context.Parse(session)
	if err != nil {
		log.Error(err.Error())
		return
	}

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
		log.Error("Unknown accessory type:", accyType)
		return
	}

	ctx.Mutex.Lock()
	ok, err := ctx.Char.Equipment.EquipAccessory(slot, slots, ctx.Char.Inventory)
	ctx.Mutex.Unlock()

	if err != nil {
		log.Error(err.Error())
	}

	pkt := network.NewWriter(ACCESSORY_EQUIP)
	pkt.WriteBool(ok)

	session.Send(pkt)
}
