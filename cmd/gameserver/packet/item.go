package packet

import (
	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/share/models/inventory"
	"github.com/ubis/Freya/share/network"
)

func StackItem(dst *inventory.Item, src context.ItemHandler) bool {
	// should not be stacked: no items with this kind index
	if dst.Kind == 0 {
		return false
	}

	// should not be stacked: item kinds are different
	if dst.Kind != src.GetKind() {
		return false
	}

	// todo: check stackable option
	// todo: check quest items

	amount := src.GetOption()
	total := dst.Option + amount

	// todo: check total amount based on item type

	dst.Option = total

	return true
}

// ItemLooting Packet
func ItemLooting(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	id := reader.ReadInt32()
	key := reader.ReadUint16()
	_ = reader.ReadUint32() // kind idx
	slot := reader.ReadUint16()

	world := GetCurrentWorld(session)
	if world == nil {
		session.LogErrorf("Unable to find current world for character: %d ",
			session.Character.Id)
		return
	}

	const (
		statusOk = 0x60 + iota
		statusOwnerFail
		statusAfterImage
		statusAlreadyUseSlot
		statusAntiOnlineGame
		statusOutOfRange
	)

	item := world.PeekItem(id, key)
	if item == nil {
		// such item does not exist
		return
	}

	// check looting position
	x, y := session.Character.GetPosition()
	itemX, itemY := item.GetPosition()
	distX, distY := int(x)-int(itemX), int(y)-int(itemY)
	if distX < -10 || distX > 10 || distY < -10 || distY > 10 {
		pkt := network.NewWriter(CSCItemLooting)
		pkt.WriteByte(statusOutOfRange)
		pkt.WriteUint32(0)
		pkt.WriteInt32(0)
		pkt.WriteUint16(0)
		pkt.WriteUint32(0)
		session.Send(pkt)
		return
	}

	// owner pick-up duration was not expired yet
	if !item.IsOwnerExpired() && item.GetOwner() != session.Character.Id {
		pkt := network.NewWriter(CSCItemLooting)
		pkt.WriteByte(statusOwnerFail)
		pkt.WriteUint32(0)
		pkt.WriteInt32(0)
		pkt.WriteUint16(0)
		pkt.WriteUint32(0)
		session.Send(pkt)
		return
	}

	// check inventory slot is already in use
	// also check if it should be stacked
	currentItem := session.Character.Inventory.Get(slot)
	isStacked := StackItem(&currentItem, item)
	if !isStacked && currentItem.Kind != 0 {
		pkt := network.NewWriter(CSCItemLooting)
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

	if isStacked {
		invItem.Option = currentItem.Option
	}

	state := false
	var err error

	if isStacked {
		state, err = session.Character.Inventory.Stack(slot, currentItem.Option)
	} else {
		state, err = session.Character.Inventory.Set(slot, *invItem)
	}

	if err != nil {
		session.LogErrorf("An error occurred: %s for character: %d ",
			err.Error(), session.Character.Id)
	}

	pkt := network.NewWriter(CSCItemLooting)

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
	pkt := network.NewWriter(NFYNewItemList)
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

	pkt := network.NewWriter(NFYNewItemList)
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

	pkt := network.NewWriter(NFYDelItemList)
	pkt.WriteInt32(id)
	pkt.WriteByte(reason)

	return pkt
}
