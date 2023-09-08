package packet

import (
	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/character"
	"github.com/ubis/Freya/share/models/inventory"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

func updatePlayerStorage(id int32, cmd string, old uint16, new uint16) (*inventory.Item, error) {

	req := character.ItemMoveReq{
		Server:     byte(g_ServerSettings.ServerId),
		Id:         id,
		DeleteSlot: old,
		CreateSlot: new,
	}
	res := character.ItemMoveRes{}
	if err := g_RPCHandler.Call(cmd, &req, &res); err != nil {
		return nil, err
	}

	if res.Result != nil {
		return nil, res.Result
	}

	return &res.Item, nil
}

func handleItemEquip(id int32, ctx *context.Context, old, new uint16) *inventory.Item {
	item, err := updatePlayerStorage(id, rpc.EquipItem, old, new)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	ctx.Mutex.Lock()
	ctx.Char.Inventory.Remove(old)
	ctx.Char.Equipment.Set(new, *item)
	ctx.Mutex.Unlock()

	return item
}

func handleItemUnequip(id int32, ctx *context.Context, old, new uint16) *inventory.Item {
	item, err := updatePlayerStorage(id, rpc.UnEquipItem, old, new)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	ctx.Mutex.Lock()
	ctx.Char.Equipment.Remove(old)
	ctx.Char.Inventory.Set(new, *item)
	ctx.Mutex.Unlock()

	return item
}

func handleEquipMove(id int32, ctx *context.Context, old, new uint16) *inventory.Item {
	item, err := updatePlayerStorage(id, rpc.ChangeInventoryItemSlot, old, new)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	ctx.Mutex.Lock()
	ctx.Char.Equipment.Remove(old)
	ctx.Char.Equipment.Set(new, *item)
	ctx.Mutex.Unlock()

	return item
}

func handleItemMove(id int32, ctx *context.Context, old, new uint16) *inventory.Item {
	item, err := updatePlayerStorage(id, rpc.ChangeInventoryItemSlot, old, new)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	ctx.Mutex.Lock()
	ctx.Char.Inventory.Remove(old)
	ctx.Char.Inventory.Set(new, *item)
	ctx.Mutex.Unlock()

	return item
}

func notifyStorageExchange(session *network.Session, result byte) {
	packet := network.NewWriter(STORAGE_EXCHANGE_MOVE)
	packet.WriteByte(result)
	packet.WriteInt32(0)

	session.Send(packet)
}

func notifyItemEquip(id int32, item *inventory.Item) *network.Writer {
	pkt := network.NewWriter(NFY_ITEM_EQUIP)
	pkt.WriteInt32(id)
	pkt.WriteInt32(item.Kind)
	pkt.WriteInt16(item.Slot)
	pkt.WriteInt32(0)
	pkt.WriteByte(0)

	return pkt
}

func notifyItemUnequip(id int32, slot uint16) *network.Writer {
	pkt := network.NewWriter(NFY_ITEM_UNEQUIP)
	pkt.WriteInt32(id)
	pkt.WriteInt16(slot)

	return pkt
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
		item := handleItemUnequip(id, ctx, deleteSlot, createSlot)
		if item == nil {
			notifyStorageExchange(session, 0)
			return
		}

		pkt := notifyItemUnequip(id, deleteSlot)
		ctx.World.BroadcastSessionPacket(session, pkt)
	case isInventory && !isEquip:
		// from inventory to equipment
		item := handleItemEquip(id, ctx, deleteSlot, createSlot)
		if item == nil {
			notifyStorageExchange(session, 0)
			return
		}

		pkt := notifyItemEquip(id, item)
		ctx.World.BroadcastSessionPacket(session, pkt)
	case isEquip && isInventory:
		// exchanging equipment items? rings? because on weaps it doesn't work
		item := handleEquipMove(id, ctx, deleteSlot, createSlot)
		if item == nil {
			notifyStorageExchange(session, 0)
			return
		}

		// not sure about this

		pkt := notifyItemUnequip(id, deleteSlot)
		ctx.World.BroadcastSessionPacket(session, pkt)

		pkt = notifyItemEquip(id, item)
		ctx.World.BroadcastSessionPacket(session, pkt)
	case !isEquip && !isInventory:
		// moving item in inventory
		item := handleItemMove(id, ctx, deleteSlot, createSlot)
		if item == nil {
			notifyStorageExchange(session, 0)
			return
		}
	default:
		notifyStorageExchange(session, 0)
		return
	}

	notifyStorageExchange(session, 1)
}

func StorageItemSwap(session *network.Session, reader *network.Reader) {
	_ = reader.ReadInt32() // unk
	oldSlot := reader.ReadInt32()
	_ = reader.ReadInt32() // unk
	newSlot := reader.ReadInt32()

	ctx, err := context.Parse(session)
	if err != nil {
		log.Error(err.Error())
		return
	}

	ctx.Mutex.RLock()
	charId := ctx.Char.Id
	itemOld := ctx.Char.Inventory.Get(uint16(oldSlot))
	itemNew := ctx.Char.Inventory.Get(uint16(newSlot))
	ctx.Mutex.RUnlock()

	if itemOld.Kind == 0 || itemNew.Kind == 0 {
		// not found
		return
	}

	// update slot
	itemOld.Slot = uint16(newSlot)
	itemNew.Slot = uint16(oldSlot)

	state, err := syncItemSwap(charId, rpc.SwapItem, &itemOld, &itemNew)
	if err != nil {
		log.Error(err.Error())
		return
	}

	pkt := network.NewWriter(STORAGE_ITEM_SWAP)
	pkt.WriteBool(state)

	session.Send(pkt)

	if state {
		ctx.Mutex.Lock()
		ctx.Char.Inventory.Remove(uint16(oldSlot))
		ctx.Char.Inventory.Remove(uint16(newSlot))
		ctx.Char.Inventory.Set(uint16(newSlot), itemOld)
		ctx.Char.Inventory.Set(uint16(oldSlot), itemNew)
		ctx.Mutex.Unlock()
	}
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

	if item.Kind == 0 {
		// not found
		return
	}

	if world == nil {
		log.Error("Unable to get current world!")
		return
	}

	state, err := syncInventory(charId, rpc.DropItem, &item)
	if err != nil {
		log.Error(err.Error())
		return
	}

	pkt := network.NewWriter(STORAGE_ITEM_DROP)
	pkt.WriteBool(state)

	session.Send(pkt)

	if state {
		ctx.Mutex.Lock()
		ctx.Char.Inventory.Remove(slot)
		ctx.Mutex.Unlock()

		world.DropItem(&item, charId, x, y)
	}
}
