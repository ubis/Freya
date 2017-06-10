package event

import (
    "sync"
)

// Event type
type Type string

// Event interface, which is later parsed into some struct
type Event interface {

}

// Event handler func
type Handler func(Event)

var handlers = map[Type][]Handler{}
var lock sync.Mutex

/*
    Registers server event
    @param  t   Event Type, which is string, defined in const file
    @param  h   Event Handler, an func which will be called
 */
func Register(t Type, h Handler) {
    lock.Lock()
    handlers[t] = append(handlers[t], h)
    lock.Unlock()
}

/*
    Triggers server event
    @param  t   Event Type, which is string, defined in const file
    @param  e   Event Interface, which is later parsed into some struct
 */
func Trigger(t Type, e Event) {
    lock.Lock()
    hs := handlers[t]
    lock.Unlock()

    for _, h := range hs {
        go h(e)
    }
}