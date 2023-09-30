package packet

import (
	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/network"
)

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
	currentItem := ctx.Char.Inventory.Get(slot)
	x, y := ctx.Char.X, ctx.Char.Y
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

	// check looting position
	itemX, itemY := item.GetPosition()
	distX, distY := int(x)-int(itemX), int(y)-int(itemY)
	if distX < -10 || distX > 10 || distY < -10 || distY > 10 {
		pkt := network.NewWriter(ITEMLOOTING)
		pkt.WriteByte(statusOutOfRange)
		pkt.WriteUint32(0)
		pkt.WriteInt32(0)
		pkt.WriteUint16(0)
		pkt.WriteUint32(0)
		session.Send(pkt)
		return
	}

	// owner pick-up duration was not expired yet
	if !item.IsOwnerExpired() && item.GetOwner() != charId {
		pkt := network.NewWriter(ITEMLOOTING)
		pkt.WriteByte(statusOwnerFail)
		pkt.WriteUint32(0)
		pkt.WriteInt32(0)
		pkt.WriteUint16(0)
		pkt.WriteUint32(0)
		session.Send(pkt)
		return
	}

	// check inventory slot is already in use
	// note: this disables stacking items (cannot pick up potions etc)
	// todo: add item stacking
	if currentItem.Kind != 0 {
		pkt := network.NewWriter(ITEMLOOTING)
		pkt.WriteByte(statusAlreadyUseSlot)
		pkt.WriteUint32(0)
		pkt.WriteInt32(0)
		pkt.WriteUint16(0)
		pkt.WriteUint32(0)
		session.Send(pkt)
		return
	}

	invItem := world.PickItem(id)
	if invItem == nil {
		// such item does not exist
		return
	}

	// update slot
	invItem.Slot = slot

	ctx.Mutex.Lock()
	state, err := ctx.Char.Inventory.Set(slot, *invItem)
	ctx.Mutex.Unlock()

	if err != nil {
		log.Error(err.Error())
	}

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
