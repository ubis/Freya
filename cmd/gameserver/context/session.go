package context

import (
	"errors"
	"sync"

	"github.com/ubis/Freya/share/models/character"
	"github.com/ubis/Freya/share/network"
)

// Context holds information related to the player's current context within the game.
type Context struct {
	Mutex sync.RWMutex
	Char  *character.Character

	Cell         CellHandler
	World        WorldHandler
	WorldManager WorldManagerHandler
}

// Init initializes a new context for the session.
func Init(session *network.Session) {
	session.DataEx = &Context{}
}

// PreParse retrieves the context from the session before character data is set.
func PreParse(session *network.Session) (*Context, error) {
	err := errors.New("unable to parse session context")

	if session.DataEx == nil {
		// we have invalid session, ignore
		return nil, err
	}

	ctx, ok := session.DataEx.(*Context)
	if !ok {
		// we have invalid session, ignore
		return nil, err
	}

	return ctx, nil
}

// Parse retrieves the context from the session.
// It ensures that the context is valid, the character is set.
func Parse(session *network.Session) (*Context, error) {
	err := errors.New("unable to parse session context")

	if session.DataEx == nil {
		// we have invalid session, ignore
		return nil, err
	}

	ctx, ok := session.DataEx.(*Context)
	if !ok {
		// we have invalid session, ignore
		return nil, err
	}

	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()

	if ctx.Char == nil {
		// session is in the lobby, we cannot receive such messages
		return nil, err
	}

	return ctx, nil
}
