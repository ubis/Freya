package game

import (
	"math/rand"

	"github.com/ubis/Freya/share/models/inventory"
)

type Item struct {
	Id    int32
	Key   uint16
	Owner int32
	X     uint16
	Y     uint16
	Item  *inventory.Item
}

func NewItem(item *inventory.Item, id, owner int32, x, y int) *Item {
	return &Item{
		Id:    id,
		Key:   uint16(rand.Intn(65536)),
		Owner: owner,
		X:     uint16(x),
		Y:     uint16(y),
		Item:  item,
	}
}

func (i *Item) GetId() int32 {
	return i.Id
}

func (i *Item) GetOwner() int32 {
	return i.Owner
}

func (i *Item) GetKind() uint32 {
	return i.Item.Kind
}

func (i *Item) GetOption() int32 {
	return i.Item.Option
}

func (i *Item) GetPosition() (uint16, uint16) {
	return i.X, i.Y
}

func (i *Item) GetKey() uint16 {
	return i.Key
}
