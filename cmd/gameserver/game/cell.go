package game

import (
	"sync"
	"time"

	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/cmd/gameserver/packet"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/server"
	"github.com/ubis/Freya/share/network"
)

// Cell represents a cell in the world grid.
type Cell struct {
	world  *World
	column byte
	row    byte

	pmutex  sync.RWMutex
	players map[uint16]network.SessionHandler

	mobs   map[int]*Mob
	mmutex sync.RWMutex

	items  map[int32]*Item
	imutex sync.RWMutex

	attribute [worldMapCellColumn * worldMapCellRow]uint32
}

// sendPlayers sends information about all players to a specified session.
func (c *Cell) sendPlayers(session network.SessionHandler) {
	c.pmutex.RLock()
	defer c.pmutex.RUnlock()

	if len(c.players) == 0 {
		return
	}

	pkt := packet.NewUserList(c.players, server.NewUserNone)
	if pkt == nil {
		return
	}

	session.Send(pkt)
}

// sendMobs sends a packet containing information about all mobs in the cell
// to a specified session.
func (c *Cell) sendMobs(session network.SessionHandler) {
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

// sendItems sends a packet containing information about all items in the cell
// to a specified session.
func (c *Cell) sendItems(session network.SessionHandler) {
	c.imutex.RLock()
	defer c.imutex.RUnlock()

	if len(c.items) == 0 {
		return
	}

	var items []context.ItemHandler

	for _, v := range c.items {
		items = append(items, v)
	}

	pkt := packet.NewItemList(items)
	if pkt != nil {
		session.Send(pkt)
	}
}

// Initialize initializes a Cell with its column and row coordinates.
func (c *Cell) Initialize(world *World) {
	c.world = world
	c.players = make(map[uint16]network.SessionHandler)
	c.mobs = make(map[int]*Mob)
	c.items = make(map[int32]*Item)
}

// GetId returns the column and row values of the cell.
func (c *Cell) GetId() (byte, byte) {
	return c.column, c.row
}

func (c *Cell) GetItemCount() int32 {
	c.imutex.RLock()
	defer c.imutex.RUnlock()

	return int32(len(c.items))
}

// AddPlayer adds a player's session to the cell.
func (c *Cell) AddPlayer(session network.SessionHandler) {
	log.Debugf("Adding %d player to %d:%d cell",
		session.GetUserIdx(), c.column, c.row)

	c.pmutex.Lock()
	defer c.pmutex.Unlock()

	index := session.GetUserIdx()
	c.players[index] = session
}

// RemovePlayer removes a player's session from the cell.
func (c *Cell) RemovePlayer(session network.SessionHandler) {
	log.Debugf("Removing %d player from %d:%d cell",
		session.GetUserIdx(), c.column, c.row)

	c.pmutex.Lock()
	defer c.pmutex.Unlock()

	index := session.GetUserIdx()
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

// AddItem adds an item to the cell.
func (c *Cell) AddItem(item *Item) {
	c.imutex.Lock()
	c.items[item.Id] = item
	c.imutex.Unlock()
}

// RemoveItem removes an item from the cell.
func (c *Cell) RemoveItem(item *Item) {
	c.imutex.Lock()
	delete(c.items, item.Id)
	c.imutex.Unlock()
}

func (c *Cell) FindItem(id int32) *Item {
	c.imutex.RLock()
	defer c.imutex.RUnlock()

	return c.items[id]
}

// SendState sends the state of the cell to a specified session.
func (c *Cell) SendState(session network.SessionHandler) {
	c.sendPlayers(session)
	c.sendMobs(session)
	c.sendItems(session)
}

func (c *Cell) Schedule() {
	c.imutex.Lock()
	defer c.imutex.Unlock()

	now := time.Now()

	// run through all items to check expiration
	for k, v := range c.items {
		expire := v.Created.Add(v.Expire)

		if !now.After(expire) {
			continue
		}

		// item has expired
		delete(c.items, k)

		// broadcast info
		pkt := packet.DelItemList(v.Id, 0x14)
		if pkt != nil {
			c.world.sendToNearbyCells(pkt, c.column, c.row, 2)
		}
	}
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
