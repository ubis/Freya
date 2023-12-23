package context

import (
	"github.com/ubis/Freya/share/models/inventory"
	"github.com/ubis/Freya/share/models/server"
	"github.com/ubis/Freya/share/network"
)

// CellHandler defines the interface for interacting with a world map cell.
type CellHandler interface {
	GetId() (byte, byte)
	AddPlayer(session network.SessionHandler)
	RemovePlayer(session network.SessionHandler)
	Send(pkt *network.Writer)
}

// WorldHandler defines the interface for interacting with a game world.
type WorldHandler interface {
	EnterWorld(session network.SessionHandler)
	ExitWorld(session network.SessionHandler, reason server.DelUserType)
	AdjustCell(session network.SessionHandler)
	FindCell(x, y int) CellHandler
	BroadcastSessionPacket(session network.SessionHandler, pkt *network.Writer)
	FindWarp(warp byte) *Warp
	IsMovable(x, y int) bool
	DropItem(item *inventory.Item, owner int32, x, y int) bool
	PickItem(id int32) *inventory.Item
	PeekItem(id int32, key uint16) ItemHandler
}

// WorldManagerHandler defines the interface for interacting with a world manager.
type WorldManagerHandler interface {
	FindWorld(id byte) WorldHandler
	GetWarps(world byte) []Warp
}
