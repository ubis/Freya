package inventory

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sort"

	"github.com/ubis/Freya/share/rpc"
)

type Inventory struct {
	Inv map[int]Item

	rpcHandler *rpc.Client
	character  int32
	serverId   byte
}

// Initializes Inventory
func (e *Inventory) Init() {
	e.Inv = make(map[int]Item)
}

func (e *Inventory) sync(cmd string, old *Item, new *Item) (bool, error) {
	// being initialized
	if e.character == 0 && e.serverId == 0 {
		return true, nil
	}

	if e.rpcHandler == nil {
		return false, errors.New("rpc handler is not ready")
	}

	req := ItemRequest{
		Server:  e.serverId,
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

func (e *Inventory) Setup(rpc *rpc.Client, id int32, server byte) {
	e.rpcHandler = rpc
	e.character = id
	e.serverId = server
}

// Sets inventory item by slot
func (e *Inventory) Set(slot uint16, item Item) (bool, error) {
	ok, err := e.sync(rpc.AddItem, &item, nil)
	if err == nil {
		e.Inv[int(slot)] = item
	}

	return ok, err
}

// Stack inventory item
func (e *Inventory) Stack(slot uint16, total int32) (bool, error) {
	// update amount
	item := e.Get(slot)
	item.Option = total

	ok, err := e.sync(rpc.StackItem, &item, nil)
	if err == nil {
		delete(e.Inv, int(slot))
		e.Inv[int(slot)] = item
	}

	return ok, err
}

// Returns inventory item by slot
func (e *Inventory) Get(slot uint16) Item {
	if value, ok := e.Inv[int(slot)]; ok {
		return value
	}

	return Item{}
}

// Removes inventory item by slot
func (e *Inventory) Remove(slot uint16) (bool, error) {
	item, ok := e.Inv[int(slot)]
	if !ok {
		return ok, errors.New("such item does not exist in the inventory")
	}

	ok, err := e.sync(rpc.RemoveItem, &item, nil)
	if err == nil {
		delete(e.Inv, int(slot))
	}

	return ok, err
}

// Removes inventory item by slot
func (e *Inventory) Swap(old, new uint16) (bool, error) {
	oldItem, ok := e.Inv[int(old)]
	if !ok {
		return ok, errors.New("such item does not exist in the inventory")
	}

	newItem, ok := e.Inv[int(new)]
	if !ok {
		return ok, errors.New("such item does not exist in the inventory")
	}

	// swap slots
	oldItem.Slot = new
	newItem.Slot = old

	ok, err := e.sync(rpc.SwapItem, &oldItem, &newItem)
	if err == nil {
		delete(e.Inv, int(old))
		delete(e.Inv, int(new))
		e.Inv[int(oldItem.Slot)] = oldItem
		e.Inv[int(newItem.Slot)] = newItem
	}

	return ok, err
}

// Move item to a new slot in inventory
func (e *Inventory) Move(old, new uint16) (bool, error) {
	oldItem, ok := e.Inv[int(old)]
	if !ok {
		return ok, errors.New("such item does not exist in the inventory")
	}

	if _, ok = e.Inv[int(new)]; ok {
		return ok, errors.New("such item already exists in the inventory")
	}

	// make a copy and swap slot
	newItem := oldItem
	newItem.Slot = new

	ok, err := e.sync(rpc.MoveItem, &oldItem, &newItem)
	if err == nil {
		delete(e.Inv, int(old))
		e.Inv[int(newItem.Slot)] = newItem
	}

	return ok, err
}

// Serializes inventory into byte array
func (e *Inventory) Serialize() ([]byte, int) {
	// collect keys for sorted iteration
	var keys []int
	for k := range e.Inv {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	var length = 0

	var equip bytes.Buffer
	for _, value := range keys {
		if e.Inv[value].Kind > 0 {
			binary.Write(&equip, binary.LittleEndian, e.Inv[value])
			length++
		}
	}

	return equip.Bytes(), length
}
