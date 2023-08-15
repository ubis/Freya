package context

import (
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/server"
	"github.com/ubis/Freya/share/network"
)

// CellHandler defines the interface for interacting with a world map cell.
type CellHandler interface {
	GetId() (byte, byte)
	AddPlayer(session *network.Session)
	RemovePlayer(session *network.Session)
	Send(pkt *network.Writer)
}

// WorldHandler defines the interface for interacting with a game world.
type WorldHandler interface {
	EnterWorld(session *network.Session)
	ExitWorld(session *network.Session, reason server.DelUserType)
	AdjustCell(session *network.Session)
	BroadcastSessionPacket(session *network.Session, pkt *network.Writer)
	FindWarp(warp byte) *Warp
	IsMovable(x, y int) bool
}

// WorldManagerHandler defines the interface for interacting with a world manager.
type WorldManagerHandler interface {
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
