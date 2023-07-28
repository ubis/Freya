package skills

import (
	"bytes"
	"encoding/binary"
	"sort"
)

type Links struct {
	List map[int]Link
}

// Initializes Skill link list
func (e *Links) Init() {
	e.List = make(map[int]Link)
}

// Sets skill link to the list by slot
func (e *Links) Set(slot uint16, link Link) {
	e.List[int(slot)] = link
}

// Returns skill link from the list by slot
func (e *Links) Get(slot uint16) Link {
	if value, ok := e.List[int(slot)]; ok {
		return value
	}

	return Link{}
}

// Removes skill link from the list by slot
func (e *Links) Remove(slot uint16) bool {
	if _, ok := e.List[int(slot)]; ok {
		delete(e.List, int(slot))
		return true
	}

	return false
}

// Serializes skill link list into byte array
func (e *Links) Serialize() ([]byte, int) {
	// collect keys for sorted iteration
	var keys []int
	for k := range e.List {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	var length = 0

	var equip bytes.Buffer
	for _, value := range keys {
		binary.Write(&equip, binary.LittleEndian, e.List[value])
		length++
	}

	return equip.Bytes(), length
}
