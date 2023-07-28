package inventory

import (
	"bytes"
	"encoding/binary"
	"sort"
)

type Inventory struct {
	Inv map[int]Item
}

// Initializes Inventory
func (e *Inventory) Init() {
	e.Inv = make(map[int]Item)
}

// Sets inventory item by slot
func (e *Inventory) Set(slot uint16, item Item) {
	e.Inv[int(slot)] = item
}

// Returns inventory item by slot
func (e *Inventory) Get(slot uint16) Item {
	if value, ok := e.Inv[int(slot)]; ok {
		return value
	}

	return Item{}
}

// Removes inventory item by slot
func (e *Inventory) Remove(slot uint16) bool {
	if _, ok := e.Inv[int(slot)]; ok {
		delete(e.Inv, int(slot))
		return true
	}

	return false
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
