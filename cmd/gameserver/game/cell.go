package game

import (
	"sync"

	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/cmd/gameserver/packet"
	"github.com/ubis/Freya/cmd/gameserver/packet/notify"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/server"
	"github.com/ubis/Freya/share/network"
)

// Cell represents a cell in the world grid.
type Cell struct {
	column byte
	row    byte

	pmutex  sync.RWMutex
	players map[uint16]*network.Session

	mobs   map[int]*Mob
	mmutex sync.RWMutex

	attribute [worldMapCellColumn * worldMapCellRow]uint32
}

// sendPlayers sends information about all players to a specified session.
func (c *Cell) sendPlayers(session *network.Session) {
	c.pmutex.RLock()
	defer c.pmutex.RUnlock()

	if len(c.players) == 0 {
		return
	}

	pkt := notify.NewUserList(c.players, server.NewUserNone)
	if pkt == nil {
		return
	}

	session.Send(pkt)
}

// sendMobs sends a packet containing information about all mobs in the cell
// to a specified session.
func (c *Cell) sendMobs(session *network.Session) {
	c.mmutex.RLock()
	defer c.mmutex.RUnlock()

	if len(c.mobs) == 0 {
		return
	}

	var mobs []context.MobHandler

	for _, v := range c.mobs {
		mobs = append(mobs, v)
	}

	pkt := packet.NewMobsList(mobs)
	if pkt != nil {
		session.Send(pkt)
	}
}

// Initialize initializes a Cell with its column and row coordinates.
func (c *Cell) Initialize() {
	c.players = make(map[uint16]*network.Session)
	c.mobs = make(map[int]*Mob)
}

// GetId returns the column and row values of the cell.
func (c *Cell) GetId() (byte, byte) {
	return c.column, c.row
}

// AddPlayer adds a player's session to the cell.
func (c *Cell) AddPlayer(session *network.Session) {
	id, err := context.GetCharId(session)
	if err != nil {
		log.Error("Failed to retrieve character id:", err.Error())
	}

	log.Debugf("Adding %d player to %d:%d cell", id, c.column, c.row)

	c.pmutex.Lock()
	defer c.pmutex.Unlock()

	index := session.UserIdx
	c.players[index] = session
}

// RemovePlayer removes a player's session from the cell.
func (c *Cell) RemovePlayer(session *network.Session) {
	id, err := context.GetCharId(session)
	if err != nil {
		log.Error("Failed to retrieve character id:", err.Error())
	}

	log.Debugf("Removing %d player from %d:%d cell", id, c.column, c.row)

	c.pmutex.Lock()
	defer c.pmutex.Unlock()

	index := session.UserIdx
	delete(c.players, index)
}

// AddMob adds a mob to the cell.
func (c *Cell) AddMob(mob *Mob) {
	c.mmutex.Lock()
	c.mobs[mob.Id] = mob
	c.mmutex.Unlock()
}

// RemoveMob removes a mob from the cell.
func (c *Cell) RemoveMob(mob *Mob) {
	c.mmutex.Lock()
	delete(c.mobs, mob.Id)
	c.mmutex.Unlock()
}

// SendState sends the state of the cell to a specified session.
func (c *Cell) SendState(session *network.Session) {
	c.sendPlayers(session)
	c.sendMobs(session)
}

// Send sends a network packet to all players in the cell.
func (c *Cell) Send(pkt *network.Writer) {
	c.pmutex.RLock()
	defer c.pmutex.RUnlock()

	if len(c.players) == 0 {
		return
	}

	for _, v := range c.players {
		v.Send(pkt)
	}
}

// IsMovable checks if a specific position within the cell is walkable/movable.
func (c *Cell) IsMovable(x, y int) bool {
	// Map
	map_movable := uint32(0x00000000)
	// map_unmovable := uint32(0x01000010)

	// Town
	// town_movable := uint32(0x06000020)
	// town_unmovable := uint32(0x07000030)

	index := ((x & 0x0F) | ((y & 0x0F) << 4))

	return c.attribute[index] == map_movable
}
