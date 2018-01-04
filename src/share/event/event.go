package event

import (
	"share/logger"
	"sync"
)

// Event interface, which is later parsed into some struct
type Event interface {
}

// Event handler func
type Handler func(Event)

var log = logger.Instance()
var handlers = map[string][]Handler{}
var lock sync.Mutex

// Registers a new server event
func Register(t string, h Handler) {
	log.Debugf("Registered `%s` event", t)
	lock.Lock()
	handlers[t] = append(handlers[t], h)
	lock.Unlock()
}

// Triggers server event in a goroutine
func Trigger(t string, e Event) {
	lock.Lock()
	hs := handlers[t]
	lock.Unlock()

	for _, h := range hs {
		go h(e)
	}
}
