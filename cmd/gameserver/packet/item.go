package packet

import (
	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/character"
	"github.com/ubis/Freya/share/models/inventory"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

func syncInventory(id int32, cmd string, item *inventory.Item) (bool, error) {

	req := character.ItemPickRequest{
		Server: byte(g_ServerSettings.ServerId),
		Id:     id,
		Item:   *item,
	}

	res := character.ItemPickResponse{}
	if err := g_RPCHandler.Call(cmd, &req, &res); err != nil {
		return false, err
	}

	return res.Result, nil
}

func ItemLooting(session *network.Session, reader *network.Reader) {
	id := reader.ReadInt32()
	key := reader.ReadUint16()
	_ = reader.ReadUint32() // kind idx
	slot := reader.ReadUint16()

	const (
		statusOk = 0x60 + iota
		statusOwnerFail
		statusAfterImage
		statusAlreadyUseSlot
		statusAntiOnlineGame
		statusOutOfRange
	)

	ctx, err := context.Parse(session)
	if err != nil {
		log.Error(err.Error())
		return
	}

	ctx.Mutex.RLock()
	charId := ctx.Char.Id
	world := ctx.World
	ctx.Mutex.RUnlock()

	if world == nil {
		log.Error("Unable to get current world!")
		return
	}

	item := world.PeekItem(id, key)
	if item == nil {
		// such item does not exist
		return
	}

	// todo: verify looting range
	// todo: verify ownership
	// todo: verify inventory space

	invItem := world.PickItem(id)
	if invItem == nil {
		// such item does not exist
		return
	}

	// update slot
	invItem.Slot = slot

	state, err := syncInventory(charId, rpc.PickItem, invItem)
	if err != nil {
		log.Error(err.Error())
		return
	}

	ctx.Mutex.Lock()
	ctx.Char.Inventory.Set(slot, *invItem)
	ctx.Mutex.Unlock()

	pkt := network.NewWriter(ITEMLOOTING)

	if state {
		pkt.WriteByte(statusOk)
	} else {
		pkt.WriteByte(statusOwnerFail)
	}

	pkt.WriteUint32(invItem.Kind)
	pkt.WriteInt32(invItem.Option)
	pkt.WriteUint16(slot)
	pkt.WriteUint32(invItem.Expire)

	session.Send(pkt)
}

func NewItemSingle(item context.ItemHandler, dropped bool) *network.Writer {
	pkt := network.NewWriter(NFY_NEWITEMLIST)
	pkt.WriteByte(1)

	pkt.WriteInt32(item.GetId())
	pkt.WriteInt32(item.GetOption())

	// we need to send owner id for the first time, because game client will
	// show "drop" animation
	// however, if we would send owner id everytime - the same drop animation
	// would appear in-game, and to not have this, id of -1 is sent instead
	if dropped {
		pkt.WriteInt32(item.GetOwner())
	} else {
		pkt.WriteInt32(-1)
	}

	pkt.WriteUint32(item.GetKind())

	x, y := item.GetPosition()
	pkt.WriteUint16(x)
	pkt.WriteUint16(y)
	pkt.WriteUint16(item.GetKey())
	pkt.WriteByte(0x02) // type
	pkt.WriteByte(0x06) // unk

	return pkt
}

func NewItemList(items []context.ItemHandler) *network.Writer {
	count := len(items)

	pkt := network.NewWriter(NFY_NEWITEMLIST)
	pkt.WriteByte(count)

	for _, i := range items {
		pkt.WriteInt32(i.GetId())
		pkt.WriteInt32(i.GetOption())
		pkt.WriteInt32(0xFFFFFFFF) // owner id, see NewItemSingle
		pkt.WriteUint32(i.GetKind())

		x, y := i.GetPosition()
		pkt.WriteUint16(x)
		pkt.WriteUint16(y)
		pkt.WriteUint16(i.GetKey())
		pkt.WriteByte(0x02) // type
		pkt.WriteByte(0x06) // unk

	}

	return pkt
}

func DelItemList(id int32, reason byte) *network.Writer {

	pkt := network.NewWriter(NFY_DELITEMLIST)
	pkt.WriteInt32(id)
	pkt.WriteByte(reason)

	return pkt
}
