package skills

import (
	"bytes"
	"encoding/binary"
	"sort"
)

type SkillList struct {
	List map[int]Skill
}

// Initializes Skill list
func (e *SkillList) Init() {
	e.List = make(map[int]Skill)
}

// Sets skill to the list by slot
func (e *SkillList) Set(slot uint16, skill Skill) {
	e.List[int(slot)] = skill
}

// Returns skill from the list by slot
func (e *SkillList) Get(slot uint16) Skill {
	if value, ok := e.List[int(slot)]; ok {
		return value
	}

	return Skill{}
}

// Removes skill from the list by slot
func (e *SkillList) Remove(slot uint16) bool {
	if _, ok := e.List[int(slot)]; ok {
		delete(e.List, int(slot))
		return true
	}

	return false
}

// Serializes skill list into byte array
func (e *SkillList) Serialize() ([]byte, int) {
	// collect keys for sorted iteration
	var keys []int
	for k := range e.List {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	var length = 0

	var equip bytes.Buffer
	for _, value := range keys {
		if e.List[value].Id > 0 {
			binary.Write(&equip, binary.LittleEndian, e.List[value])
			length++
		}
	}

	return equip.Bytes(), length
}
