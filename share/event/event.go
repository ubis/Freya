package event

import (
	"sync"

	"github.com/ubis/Freya/share/log"
)

// Event interface, which is later parsed into some struct
type Event struct {
	data  []any
	index int
}

// Event handler func
type Handler func(*Event)
type HandlerId uint64

var handlers = map[string]map[HandlerId]Handler{}
var nextHandlerId HandlerId = 1
var lock sync.RWMutex

func (e *Event) Gett() int {
	return e.index
}

func (e *Event) Get() (interface{}, bool) {
	if e.index >= len(e.data) {
		return nil, false
	}

	result := e.data[e.index]
	e.index++
	return result, true
}

// Registers a new server event and returns HandlerId.
// It is possible to unregister event via Unregisted with such HandlerId.
func Register(t string, h Handler) HandlerId {
	log.Debugf("Registering `%s` event", t)

	lock.Lock()
	defer lock.Unlock()

	if handlers[t] == nil {
		handlers[t] = make(map[HandlerId]Handler)
	}

	id := nextHandlerId
	handlers[t][id] = h
	nextHandlerId++

	return id
}

// Unregister removes a specific handler from the given event type.
func Unregister(t string, id HandlerId) {
	log.Debugf("Unregistering `%s` event", t)

	lock.Lock()
	defer lock.Unlock()

	if hs, ok := handlers[t]; ok {
		delete(hs, id)
		if len(hs) == 0 {
			delete(handlers, t)
		}
	}
}

// Triggers server event in a goroutine
func Trigger(t string, data ...any) {
	lock.RLock()
	defer lock.RUnlock()

	hs, ok := handlers[t]
	if !ok {
		return
	}

	for _, h := range hs {
		e := &Event{data: append([]any(nil), data...)}
		go h(e)
	}
}
