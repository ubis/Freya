package context

import (
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/network"
)

// CellHandler defines the interface for interacting with a world map cell.
type CellHandler interface {
	Initialize(column, row byte)
	GetId() (byte, byte)
	AddPlayer(session *network.Session)
	RemovePlayer(session *network.Session)
	Send(pkt *network.Writer)
}

// WorldHandler defines the interface for interacting with a game world.
type WorldHandler interface {
	Initialize(wm WorldManagerHandler)
	EnterWorld(session *network.Session)
	ExitWorld(session *network.Session)
	AdjustCell(session *network.Session)
	BroadcastPacket(session *network.Session, pkt *network.Writer)
	FindWarp(warp byte) *Warp
}

// WorldManagerHandler defines the interface for interacting with a world manager.
type WorldManagerHandler interface {
	Initialize()
	FindWorld(id byte) WorldHandler
	GetWarps(world byte) []Warp
}

// GetWorldManager retrieves the WorldManagerHandler for the given session.
func GetWorldManager(session *network.Session) WorldManagerHandler {
	ctx, err := Parse(session)
	if err != nil {
		log.Error("failed to parse session context:", err.Error())
		return nil
	}

	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()

	return ctx.WorldManager
}

// GetWorld retrieves the WorldHandler for the given session.
func GetWorld(session *network.Session) WorldHandler {
	ctx, err := Parse(session)
	if err != nil {
		log.Error("failed to parse session context:", err.Error())
		return nil
	}

	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()

	return ctx.World
}

// GetWorldCell retrieves the CellHandler for the given session.
func GetWorldCell(session *network.Session) CellHandler {
	ctx, err := Parse(session)
	if err != nil {
		log.Error("failed to parse session context:", err.Error())
		return nil
	}

	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()

	return ctx.Cell
}
