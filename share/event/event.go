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

var handlers = map[string][]Handler{}
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

// Registers a new server event
func Register(t string, h Handler) {
	log.Debugf("Registered `%s` event", t)
	lock.Lock()
	handlers[t] = append(handlers[t], h)
	lock.Unlock()
}

// Triggers server event in a goroutine
func Trigger(t string, data ...any) {
	lock.RLock()
	hs := handlers[t]
	lock.RUnlock()

	for _, h := range hs {
		go h(&Event{data: data})
	}
}
