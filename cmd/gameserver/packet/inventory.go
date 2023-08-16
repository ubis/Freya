package packet

import (
	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/character"
	"github.com/ubis/Freya/share/models/inventory"
	"github.com/ubis/Freya/share/models/server"
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
	pkt.WriteByte(c.LiveStyle) // animation id aka "live style"
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

func NewUserSingle(session *network.Session, reason server.NewUserType) *network.Writer {
	pkt := network.NewWriter(NEWUSERLIST)
	pkt.WriteByte(1) // player num
	pkt.WriteByte(byte(reason))

	fillPlayerInfo(pkt, session)

	return pkt
}

func NewUserList(players map[uint16]*network.Session, reason server.NewUserType) *network.Writer {
	online := len(players)

	pkt := network.NewWriter(NEWUSERLIST)
	pkt.WriteByte(online)
	pkt.WriteByte(byte(reason))

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

	pkt := network.NewWriter(DELUSERLIST)
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
