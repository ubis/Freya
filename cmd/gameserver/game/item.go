package game

import (
	"math/rand"
	"time"

	"github.com/ubis/Freya/share/models/inventory"
)

const (
	itemExpire      = time.Minute * 6
	itemOwnerExpire = time.Second * 20
)

type Item struct {
	Id          int32
	Key         uint16
	Owner       int32
	X           uint16
	Y           uint16
	Item        *inventory.Item
	OwnerExpire bool
	Created     time.Time
	Expire      time.Duration
}

func NewItem(item *inventory.Item, id, owner int32, x, y int, ownerExpire bool) *Item {
	return &Item{
		Id:          id,
		Key:         uint16(rand.Intn(65536)),
		Owner:       owner,
		X:           uint16(x),
		Y:           uint16(y),
		Item:        item,
		OwnerExpire: ownerExpire,
		Created:     time.Now(),
		// after expiration, item should disappear forever
		// however, when player(owner) drops it or on some bosses, owner should
		// have ~20 seconds to pick it up
		// which means that on some occasions we should add +20 to duration and
		// disallow to pick up it up for non-owner players until 10 second gets
		// passed
		Expire: itemExpire,
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

func (i *Item) IsOwnerExpired() bool {
	return time.Now().After(i.Created.Add(itemOwnerExpire))
}
