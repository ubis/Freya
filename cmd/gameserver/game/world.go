package game

import (
	"errors"

	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/cmd/gameserver/packet/notify"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/server"
	"github.com/ubis/Freya/share/network"
)

// Constants defining the dimensions of the world map grid.
const (
	worldMapCellColumn = 16 // Number of columns in the world map grid.
	worldMapCellRow    = 16 // Number of rows in the world map grid.
)

// World represents the game world and its grid of cells.
type World struct {
	Id   byte
	Grid [worldMapCellColumn][worldMapCellRow]*Cell

	Warps []context.Warp
}

// isCellValid checks if the cell coordinates are within valid bounds.
func isCellValid(c, r byte) bool {
	return c < worldMapCellColumn && r < worldMapCellRow
}

// getWorldCell calculates cell column and row by position values.
func getWorldCell(x byte, y byte, column *byte, row *byte) error {
	c := x >> 4
	r := y >> 4

	if !isCellValid(c, r) {
		return errors.New("invalid world cell")
	}

	*column = c
	*row = r

	return nil
}

// getCurrentCell retrieves the cell based on the current character position.
func (w *World) getCurrentCell(session *network.Session) *Cell {
	var x, y, column, row byte

	if err := context.GetCharPosition(session, &x, &y); err != nil {
		log.Error("Failed to get character position:", err)
		return nil
	}

	if err := getWorldCell(x, y, &column, &row); err != nil {
		log.Error("Invalid world cell:", err)
		return nil
	}

	return w.Grid[column][row]
}

// getNearbyCells returns a slice of nearby cells based on radius.
func (w *World) getNearbyCells(column, row byte, radius int) []*Cell {
	var nearbyCells []*Cell

	for i := -radius; i <= radius; i++ {
		c := int(column) + i

		if c < 0 || c >= worldMapCellColumn {
			continue
		}

		for j := -radius; j <= radius; j++ {
			r := int(row) + j

			if r < 0 || r >= worldMapCellRow {
				continue
			}

			nearbyCells = append(nearbyCells, w.Grid[c][r])
		}
	}

	return nearbyCells
}

// sendToNearbyCells sends a network packet to nearby cells.
func (w *World) sendToNearbyCells(pkt *network.Writer, column, row byte, radius int) {
	cells := w.getNearbyCells(column, row, radius)

	if len(cells) == 0 {
		return
	}

	for _, v := range cells {
		v.Send(pkt)
	}
}

// Initialize initializes the World grid with Cell instances.
func (w *World) Initialize(wm context.WorldManagerHandler) {
	for i := byte(0); i < worldMapCellColumn; i++ {
		for j := byte(0); j < worldMapCellRow; j++ {
			w.Grid[i][j] = &Cell{}
			w.Grid[i][j].Initialize(i, j)
		}
	}

	// assign data
	w.Warps = wm.GetWarps(w.Id)

	log.Debugf("Loaded %d warps in %d world", len(w.Warps), w.Id)
}

// EnterWorld adds a player session to the current world cell.
func (w *World) EnterWorld(session *network.Session) {
	cell := w.getCurrentCell(session)
	if cell == nil {
		log.Error("Unable to get current cell")
		return
	}

	ctx, err := context.Parse(session)
	if err != nil {
		log.Error("Unable to parse session context:", err.Error())
		return
	}

	column, row := cell.GetId()

	// notify other nearby cells about new player with radius of -2/+2
	pkt := notify.NewUserSingle(session, server.NewUserInit)
	w.sendToNearbyCells(pkt, column, row, 2)

	// notify player about nearby cell states
	cells := w.getNearbyCells(column, row, 2)
	for _, v := range cells {
		v.SendState(session)
	}

	// add player to the cell
	cell.AddPlayer(session)

	// update player context to set current world and cell
	ctx.Mutex.Lock()
	ctx.World = w
	ctx.Cell = cell
	ctx.Mutex.Unlock()
}

// ExitWorld removes a player session from the current world cell.
func (w *World) ExitWorld(session *network.Session, reason server.DelUserType) {
	ctx, err := context.Parse(session)
	if err != nil {
		log.Error("Unable to parse session context:", err.Error())
		return
	}

	ctx.Mutex.RLock()
	cell := ctx.Cell
	ctx.Mutex.RUnlock()

	if cell == nil {
		// player was not in the world
		return
	}

	// remove current world and cell from player's context
	ctx.Mutex.Lock()
	ctx.World = nil
	ctx.Cell = nil
	ctx.Mutex.Unlock()

	// remove player from the cell
	cell.RemovePlayer(session)

	// notify other players about leaving player
	pkt := notify.DelUserList(session, reason)
	if pkt == nil {
		return
	}

	column, row := cell.GetId()
	w.sendToNearbyCells(pkt, column, row, 2)
}

// AdjustCell adjusts the player's cell based on position changes.
func (w *World) AdjustCell(session *network.Session) {
	ctx, err := context.Parse(session)
	if err != nil {
		log.Error("Unable to parse session context:", err.Error())
		return
	}

	ctx.Mutex.RLock()
	cell := ctx.Cell
	ctx.Mutex.RUnlock()

	if cell == nil {
		// player is not in the world
		return
	}

	// get new cell from current position
	newCell := w.getCurrentCell(session)
	if newCell == nil {
		log.Error("Unable to retrieve new cell!")
		return
	}

	// remove player from the cell
	cell.RemovePlayer(session)

	// update current cell
	ctx.Mutex.Lock()
	ctx.Cell = newCell
	ctx.Mutex.Unlock()

	// get column and row from new cell
	nc, nr := newCell.GetId()

	pkt := notify.NewUserSingle(session, server.NewUserMove)
	if pkt == nil {
		log.Error("Failed to create NewUserSingle packet!")
		return
	}

	w.sendToNearbyCells(pkt, nc, nr, 2)

	// notify player about nearby cell states
	cells := w.getNearbyCells(nc, nr, 2)
	for _, v := range cells {
		v.SendState(session)
	}

	// add player to the new cell
	newCell.AddPlayer(session)
}

// BroadcastPacket broadcasts a packet to nearby cells.
func (w *World) BroadcastPacket(column, row byte, pkt *network.Writer) {
	w.sendToNearbyCells(pkt, column, row, 2)
}

// BroadcastSessionPacket sends a packet to nearby cells centered around
// the cell that the given session's character currently occupies.
func (w *World) BroadcastSessionPacket(session *network.Session, pkt *network.Writer) {
	cell := context.GetWorldCell(session)
	if cell == nil {
		log.Error("Unable to get current cell!")
		return
	}

	column, row := cell.GetId()
	w.BroadcastPacket(column, row, pkt)
}

// FindWarp finds a specific warp based on its ID.
func (w *World) FindWarp(warp byte) *context.Warp {
	for _, v := range w.Warps {
		if v.Id != warp {
			continue
		}

		return &v
	}

	return nil
}
