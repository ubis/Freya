package game

import (
	"sync"

	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/cmd/gameserver/packet/notify"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/server"
	"github.com/ubis/Freya/share/network"
)

// Cell represents a cell in the world grid.
type Cell struct {
	column  byte
	row     byte
	mutex   sync.RWMutex
	players map[uint16]*network.Session
}

// sendPlayers sends information about all players to a specified session.
func (c *Cell) sendPlayers(session *network.Session) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if len(c.players) == 0 {
		return
	}

	pkt := notify.NewUserList(c.players, server.NewUserNone)
	if pkt == nil {
		return
	}

	session.Send(pkt)
}

// Initialize initializes a Cell with its column and row coordinates.
func (c *Cell) Initialize(column, row byte) {
	c.column = column
	c.row = row
	c.players = make(map[uint16]*network.Session)
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

	c.mutex.Lock()
	defer c.mutex.Unlock()

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

	c.mutex.Lock()
	defer c.mutex.Unlock()

	index := session.UserIdx
	delete(c.players, index)
}

// SendState sends the state of the cell to a specified session.
func (c *Cell) SendState(session *network.Session) {
	c.sendPlayers(session)
}

// Send sends a network packet to all players in the cell.
func (c *Cell) Send(pkt *network.Writer) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if len(c.players) == 0 {
		return
	}

	for _, v := range c.players {
		v.Send(pkt)
	}
}
