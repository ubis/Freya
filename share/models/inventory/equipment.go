package inventory

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/ubis/Freya/share/rpc"
)

type Equipment struct {
	Equip    map[int]Item
	mutex    sync.RWMutex
	mutexOut sync.Mutex // used for multi-transaction operations

	rpcHandler *rpc.Client
	character  int32
}

// Initializes Equipment
func (e *Equipment) Init() {
	e.Equip = make(map[int]Item)
}

func (e *Equipment) sync(cmd string, old *Item, new *Item) (bool, error) {
	// being initialized
	if e.character == 0 {
		return true, nil
	}

	if e.rpcHandler == nil {
		return false, errors.New("rpc handler is not ready")
	}

	req := ItemRequest{
		Id:      e.character,
		Command: cmd,
		Item:    *old,
		NewItem: new,
	}

	res := ItemResponse{}
	if err := e.rpcHandler.Call(cmd, &req, &res); err != nil {
		return false, err
	}

	return res.Result, nil
}

func (e *Equipment) Setup(rpc *rpc.Client, id int32) {
	e.rpcHandler = rpc
	e.character = id
}

// Sets equipment item by slot
func (e *Equipment) Set(slot uint16, item Item) (bool, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	ok, err := e.sync(rpc.EquipItem, &item, nil)
	if err == nil {
		e.Equip[int(slot)] = item
	}

	return ok, err
}

// Returns equipment item by slot
func (e *Equipment) Get(slot uint16) Item {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	if value, ok := e.Equip[int(slot)]; ok {
		return value
	}

	return Item{}
}

// Removes equipment item by slot
func (e *Equipment) Remove(slot uint16) (bool, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	item, ok := e.Equip[int(slot)]
	if !ok {
		return ok, errors.New("such item does not exist in the equipment")
	}

	ok, err := e.sync(rpc.UnEquipItem, &item, nil)
	if err == nil {
		delete(e.Equip, int(slot))
	}

	return ok, err
}

func (e *Equipment) EquipItem(old, new uint16, i *Inventory) (bool, error) {
	e.mutexOut.Lock()
	defer e.mutexOut.Unlock()

	// take from inventory (old) and move into equipment(new)
	item := i.Get(old)

	if item := e.Get(new); item.Kind != 0 {
		return false, errors.New("such item already exists in the equipment")
	}

	if ok, err := i.Remove(old); !ok {
		return ok, err
	}

	// update slot
	item.Slot = new

	ok, err := e.Set(new, item)
	if ok {
		return ok, nil
	}

	// attempt to rollback
	// todo: actually we should do this in a single transaction, like before
	item.Slot = old
	if ok, err2 := i.Set(old, item); !ok {
		return ok, fmt.Errorf("unable to rollback: %s and %s",
			err.Error(), err2.Error())
	}

	return ok, err
}

func (e *Equipment) SwapEquipItem(old, new uint16, i *Inventory) (bool, error) {
	e.mutexOut.Lock()
	defer e.mutexOut.Unlock()

	// swap equipment item from the inventory
	oldItem := e.Get(new)
	newItem := i.Get(old)

	if oldItem.Kind == 0 {
		return false, errors.New("such item does not exist in the equipment")
	}

	if newItem.Kind == 0 {
		return false, errors.New("such item does not exist in the inventory")
	}

	// remove new item from the inventory
	if ok, err := i.Remove(old); !ok {
		return ok, err
	}

	if ok, err := e.UnEquipItem(new, old, i); !ok {
		// todo: attempt to rollback
		return ok, err
	}

	newItem.Slot = new
	ok, err := e.Set(new, newItem)
	if !ok {
		// todo: attempt to rollback
		return ok, err
	}

	return ok, err
}

func (e *Equipment) UnEquipItem(old, new uint16, i *Inventory) (bool, error) {
	e.mutexOut.Lock()
	defer e.mutexOut.Unlock()

	// take from equipment (old) and move into inventory(new)
	item, ok := e.Equip[int(old)]
	if !ok {
		return ok, errors.New("such item does not exist in the equipment")
	}

	if ok, err := e.Remove(old); !ok {
		return ok, err
	}

	// update slot
	item.Slot = new

	ok, err := i.Set(new, item)
	if ok {
		return ok, nil
	}

	// attempt to rollback
	// todo: actually we should do this in a single transaction, like before
	item.Slot = old
	if ok, err2 := i.Set(old, item); !ok {
		return ok, fmt.Errorf("unable to rollback: %s and %s",
			err.Error(), err2.Error())
	}

	return ok, err
}

// Swap inventory item by slot
func (e *Equipment) Swap(old, new uint16) (bool, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	oldItem, ok := e.Equip[int(old)]
	if !ok {
		return ok, errors.New("such item does not exist in the equipment")
	}

	newItem, ok := e.Equip[int(new)]
	if !ok {
		return ok, errors.New("such item does not exist in the equipment")
	}

	// swap slots
	oldItem.Slot = new
	newItem.Slot = old

	ok, err := e.sync(rpc.SwapEquipmentItem, &oldItem, &newItem)
	if err == nil {
		delete(e.Equip, int(old))
		delete(e.Equip, int(new))
		e.Equip[int(oldItem.Slot)] = oldItem
		e.Equip[int(newItem.Slot)] = newItem
	}

	return ok, err
}

func (e *Equipment) MoveItem(old, new uint16) (bool, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// take from equipment (old) and move into equipment(new)
	oldItem, ok := e.Equip[int(old)]
	if !ok {
		return ok, errors.New("such item does not exist in the equipment")
	}

	newItem, ok := e.Equip[int(new)]
	if ok {
		return ok, errors.New("such item already exists in the equipment")
	}

	// set up new slot
	newItem.Slot = new

	ok, err := e.sync(rpc.MoveEquipmentItem, &oldItem, &newItem)
	if err == nil {
		delete(e.Equip, int(old))

		// swap slot
		oldItem.Slot = new
		e.Equip[int(oldItem.Slot)] = oldItem
	}

	return ok, err
}

func (e *Equipment) EquipAccessory(old uint16, slots []EquipmentType, i *Inventory) (bool, error) {
	// primary accessory slot is empty
	// no need to do any shifting
	if item := e.Get(uint16(slots[0])); item.Kind == 0 {
		return e.EquipItem(old, uint16(slots[0]), i)
	}

	item := i.Get(old)

	// first remove item we want to equip from inventory
	ok, err := i.Remove(old)
	if !ok {
		return ok, err
	}

	total := len(slots)
	first := uint16(slots[0])
	last := uint16(slots[total-1])

	// if last slot is used, then un-equip and put it in the `old` slot
	if e.Get(last).Kind != 0 {
		ok, err := e.UnEquipItem(last, old, i)
		if !ok {
			// attempt to rollback
			// todo: actually we should do this in a single transaction
			if ok, err2 := i.Set(old, item); !ok {
				return ok, fmt.Errorf("unable to rollback: %s and %s",
					err.Error(), err2.Error())
			}
			return ok, err
		}
	}

	// now attempt to move from left to the right
	for i := total - 2; i >= 0; i-- {
		e.MoveItem(uint16(slots[i]), uint16(slots[i+1])) // fixme
	}

	// finally equip to the primary slot
	item.Slot = first
	ok, err = e.Set(item.Slot, item)
	if !ok {
		// fixme
		// attempt to rollback
		return ok, err
	}

	return ok, err
}

// Serializes equipment into byte array
func (e *Equipment) Serialize() ([]byte, int) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// collect keys for sorted iteration
	var keys []int
	for k := range e.Equip {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	var length = 0

	var equip bytes.Buffer
	for _, value := range keys {
		binary.Write(&equip, binary.LittleEndian, e.Equip[value])
		length++
	}

	return equip.Bytes(), length
}

// Serializes equipment kind_idx into byte array
func (e *Equipment) SerializeKind() []byte {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// collect keys for sorted iteration
	var keys []int
	for k := range eqTypes {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	var equip bytes.Buffer
	for key := range keys {
		item, ok := e.Equip[key]
		if !ok {
			item = Item{}
			item.Slot = uint16(key)
		}

		binary.Write(&equip, binary.LittleEndian, item.Kind)
	}

	return equip.Bytes()
}

// SerializeEx serializes equipment with kind and option into byte array
func (e *Equipment) SerializeEx() ([]byte, int) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// collect keys for sorted iteration
	var keys []int
	for k := range e.Equip {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	var length = 0

	var equip bytes.Buffer
	for _, value := range keys {
		binary.Write(&equip, binary.LittleEndian, byte(e.Equip[value].Slot))
		binary.Write(&equip, binary.LittleEndian, e.Equip[value].Kind)
		binary.Write(&equip, binary.LittleEndian, e.Equip[value].Option)
		length++
	}

	return equip.Bytes(), length
}
