package game

import (
	"time"

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
	manager *WorldManager

	Id   byte
	Grid [worldMapCellColumn][worldMapCellRow]*Cell

	// mobs
	mobsMove *time.Ticker
	mobs     map[int]*Mob

	// data
	Warps []context.Warp
}

// isCellValid checks if the cell coordinates are within valid bounds.
func isCellValid(c, r int) bool {
	return c >= 0 && c < worldMapCellColumn && r >= 0 && r < worldMapCellRow
}

// getWorldCell calculates cell column and row by position values.
func (w *World) getWorldCell(x, y int) *Cell {
	column := x >> 4
	row := y >> 4

	if !isCellValid(column, row) {
		return nil
	}

	return w.Grid[column][row]
}

// getCurrentCell retrieves the cell based on the current character position.
func (w *World) getCurrentCell(session *network.Session) *Cell {
	var x, y byte

	if err := context.GetCharPosition(session, &x, &y); err != nil {
		log.Error("Failed to get character position:", err)
		return nil
	}

	cell := w.getWorldCell(int(x), int(y))
	if cell == nil {
		log.Error("Invalid world cell")
		return nil
	}

	return cell
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

// initializeCells creates a grid of Cells to represent the game world.
func (w *World) initializeCells() {
	for i := byte(0); i < worldMapCellColumn; i++ {
		for j := byte(0); j < worldMapCellRow; j++ {
			cell := &Cell{column: i, row: j}
			cell.Initialize()

			w.Grid[i][j] = cell
		}
	}
}

// initializeMobs sets up the mobs in the game world using the provided list of mobs.
func (w *World) initializeMobs(mobs []*Mob) {
	w.mobs = make(map[int]*Mob)

	for k, v := range mobs {
		cell := w.getWorldCell(int(v.SpawnX), int(v.SpawnY))
		if cell == nil {
			log.Error("Invalid world cell")
			continue
		}

		mob := mobs[k]
		mob.Id = k + 1
		mob.world = w
		mob.cell = cell

		if parentMob := w.manager.GetMob(mob.Species); parentMob != nil {
			mob.Merge(parentMob)
		}

		w.mobs[mob.Id] = mob
		mob.Initialize()
		cell.AddMob(mob)
	}
}

// initializeTimers sets up and starts timers user in the game world.
func (w *World) initializeTimers() {
	// Initialize a ticker for monster movement
	w.mobsMove = time.NewTicker(time.Millisecond * 150)

	go func() {
		for range w.mobsMove.C {
			for _, mob := range w.mobs {
				mob.Update()
			}
		}
	}()
}

// Initialize initializes the World grid with Cell instances.
func (w *World) Initialize(manager *WorldManager) {
	w.manager = manager

	// load & assign data
	w.Warps = manager.GetWarps(w.Id)
	mobs := w.loadMobs()

	log.Debugf("Loaded %d warps in %d world", len(w.Warps), w.Id)
	log.Debugf("Loaded %d mobs in %d world", len(mobs), w.Id)

	// initialize everything
	w.initializeCells()
	w.loadThreadMap()
	w.initializeMobs(mobs)
	w.initializeTimers()
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

// IsMovable determines if a specific position (x, y) within the world is
// navigable by characters or entities.
func (w *World) IsMovable(x, y int) bool {
	cell := w.getWorldCell(x, y)
	if cell == nil {
		return false
	}

	return cell.IsMovable(x, y)
}
