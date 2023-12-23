package skills

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sort"
	"sync"

	"github.com/ubis/Freya/share/rpc"
)

type Links struct {
	List  map[int]Link
	mutex sync.RWMutex

	rpcHandler *rpc.Client
	character  int32
}

func (e *Links) sync(cmd string, old *Link, new *Link) (bool, error) {
	// being initialized
	if e.character == 0 {
		return true, nil
	}

	if e.rpcHandler == nil {
		return false, errors.New("rpc handler is not ready")
	}

	req := QuickLinkRequest{
		Id:      e.character,
		Command: cmd,
		OldLink: old,
		NewLink: new,
	}

	res := QuickLinkResponse{}
	if err := e.rpcHandler.Call(cmd, &req, &res); err != nil {
		return false, err
	}

	return res.Result, nil
}

// Initializes Skill link list
func (e *Links) Init() {
	e.List = make(map[int]Link)
}

func (e *Links) Setup(rpc *rpc.Client, id int32) {
	e.rpcHandler = rpc
	e.character = id
}

// Sets skill link to the list by slot
func (e *Links) Set(slot uint16, link Link) (bool, error) {
	link.Slot = slot

	ok, err := e.sync(rpc.QuickLinkSet, nil, &link)
	if err == nil {
		e.mutex.Lock()
		e.List[int(slot)] = link
		e.mutex.Unlock()
	}

	return ok, err
}

// Swaps link from old to a new slot
func (e *Links) Swap(old, new uint16) (bool, error) {
	oldLink := e.Get(old)
	newLink := e.Get(new)

	if oldLink == nil {
		return false, errors.New("such old link does not exist")
	}

	if newLink == nil {
		return false, errors.New("such new link does not exist")
	}

	// switch slots
	oldLink.Slot = new
	newLink.Slot = old

	e.mutex.Lock()
	defer e.mutex.Unlock()

	ok, err := e.sync(rpc.QuickLinkSwap, oldLink, newLink)
	if err == nil {
		delete(e.List, int(old))
		delete(e.List, int(new))
		e.List[int(new)] = *oldLink
		e.List[int(old)] = *newLink
	}

	return ok, err
}

// Returns skill link from the list by slot
func (e *Links) Get(slot uint16) *Link {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if value, ok := e.List[int(slot)]; ok {
		return &value
	}

	return nil
}

// Removes skill link from the list by slot
func (e *Links) Remove(slot uint16) (bool, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	link, ok := e.List[int(slot)]
	if !ok {
		return ok, errors.New("such link does not exist in the link bar")
	}

	ok, err := e.sync(rpc.QuickLinkRemove, &link, nil)
	if err == nil {
		delete(e.List, int(slot))
	}

	return ok, err
}

// Serializes skill link list into byte array
func (e *Links) Serialize() ([]byte, int) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

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
