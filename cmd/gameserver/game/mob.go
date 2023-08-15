package game

import (
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/beefsack/go-astar"
	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/cmd/gameserver/packet"
	"github.com/ubis/Freya/share/log"
)

type MobStateType int

const (
	MobStateNone MobStateType = iota
	MobStateFind
	MobStateMove
)

// Mob struct represents a monster in the game world.
type Mob struct {
	Id           int     `yaml:"-"`
	Species      int     `yaml:"id"`
	SpawnX       int16   `yaml:"posx"`
	SpawnY       int16   `yaml:"posy"`
	MoveSpeed    float32 `yaml:"movespeed"`
	FindCount    int     `yaml:"findcount"`
	FindInterval int64   `yaml:"findinterval"`
	MoveInterval int64   `yaml:"moveinterval"`
	Level        int     `yaml:"lev"`
	MaxHP        int     `yaml:"hp"`
	CurrentHP    int     `yaml:"-"`

	position   *context.Position
	state      MobStateType
	findRemain int
	nextRun    int64

	mutex sync.RWMutex

	world *World
	cell  *Cell
}

// Merge integrates the properties of another mob into the current one,
// without overwriting zero values.
func (m *Mob) Merge(mob *Mob) {
	assignIfNotZero := func(target interface{}, value interface{}) {
		targetValue := reflect.ValueOf(target).Elem()
		valueValue := reflect.ValueOf(value)

		if valueValue.Kind() == reflect.Ptr {
			valueValue = valueValue.Elem()
		}

		if valueValue.Interface() != reflect.Zero(valueValue.Type()).Interface() {
			targetValue.Set(valueValue)
		}
	}

	assignIfNotZero(&m.Species, mob.Species)
	assignIfNotZero(&m.SpawnX, mob.SpawnX)
	assignIfNotZero(&m.SpawnY, mob.SpawnY)
	assignIfNotZero(&m.MoveSpeed, mob.MoveSpeed)
	assignIfNotZero(&m.FindCount, mob.FindCount)
	assignIfNotZero(&m.FindInterval, mob.FindInterval)
	assignIfNotZero(&m.MoveInterval, mob.MoveInterval)
	assignIfNotZero(&m.Level, mob.Level)
	assignIfNotZero(&m.MaxHP, mob.MaxHP)
}

// Initialize sets up a mob with its default values, especially for its position and state.
func (m *Mob) Initialize() {
	m.CurrentHP = m.MaxHP

	// initialize position data
	m.position = &context.Position{
		InitialX: int(m.SpawnX),
		InitialY: int(m.SpawnY),
		CurrentX: int(m.SpawnX),
		CurrentY: int(m.SpawnY),
		FinalX:   int(m.SpawnX),
		FinalY:   int(m.SpawnY),
	}

	m.SetState(MobStateFind)
}

// GetId retrieves the mob's ID.
func (m *Mob) GetId() int {
	return m.Id
}

// GetSpecies retrieves the species of the mob.
func (m *Mob) GetSpecies() int {
	return m.Species
}

// GetHealth safely retrieves the current and max health of the mob.
func (m *Mob) GetHealth() (int, int) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.CurrentHP, m.MaxHP
}

// GetPosition safely retrieves the mob's position in the world.
func (m *Mob) GetPosition() context.Position {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return *m.position
}

// SetPosition safely sets a new position for the mob.
func (m *Mob) SetPosition(pos *context.Position) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	tmp := *pos

	if pos.WayPoints != nil {
		tmp.WayPoints = make([]context.WayPointHandler, len(pos.WayPoints))
		copy(tmp.WayPoints, pos.WayPoints)
	}

	m.position = &tmp
}

// SetState determines the next action the mob will take and when it will do so.
func (m *Mob) SetState(state MobStateType) {
	next := time.Now().UnixMilli()

	switch state {
	case MobStateFind:
		next += m.FindInterval
	case MobStateMove:
		next += m.MoveInterval
	}

	m.mutex.Lock()
	m.state = state
	m.nextRun = next
	m.findRemain = m.FindCount
	m.mutex.Unlock()
}

// Update checks if the mob needs to take action based on its state and the
// current time.
func (m *Mob) Update() {
	now := time.Now().UnixMilli()

	m.mutex.RLock()
	nextRun := m.nextRun
	state := m.state
	m.mutex.RUnlock()

	if nextRun > now {
		return
	}

	switch state {
	case MobStateFind:
		m.handleFind(now)
	case MobStateMove:
		m.handleMove(now)
	}
}

// handleFind deals with the logic when the mob is in the 'find' state.
func (m *Mob) handleFind(now int64) {
	m.mutex.RLock()
	pending := m.findRemain
	m.mutex.RUnlock()

	if pending <= 1 {
		m.SetState(MobStateMove)
		return
	}

	next := now + m.FindInterval
	m.mutex.Lock()
	m.findRemain--
	m.nextRun = next
	m.mutex.Unlock()
}

// handleMove handles the logic of a mob's movement through the world.
func (m *Mob) handleMove(now int64) {
	next := now + m.MoveInterval

	m.mutex.Lock()
	m.nextRun = next
	m.mutex.Unlock()

	pos := m.GetPosition()

	if !pos.IsMoving {
		if !m.findNewPath(&pos) {
			m.SetState(MobStateFind)
			return
		}

		path, distance := m.findMovePath(&pos)
		if distance == 0 {
			m.SetState(MobStateFind)
			return
		}

		start := path[int(distance)].(WayPoint)
		end := path[0].(WayPoint)

		pos.WayPoints = pos.WayPoints[:0]
		pos.WayPoints = append(pos.WayPoints, start, end)
		pos.CurrentWayPoint = 0
		pos.Speed = m.MoveSpeed

		OpenDeadReckoning(&pos)
		m.SetPosition(&pos)

		column, row := m.cell.GetId()

		pkt := packet.MobMoveBegin(m)
		m.world.BroadcastPacket(column, row, pkt)
		return
	}

	if pos.CurrentX == pos.FinalX && pos.CurrentY == pos.FinalY {
		pos.IsDeadReckoning = false
	} else {
		DeadReckoning(&pos)
		m.SetPosition(&pos)
		m.adjustCell(&pos)
	}

	if pos.IsDeadReckoning {
		m.SetPosition(&pos)
		return
	}

	pos.InitialX = pos.FinalX
	pos.InitialY = pos.FinalY
	pos.IsMoving = false

	m.SetPosition(&pos)

	column, row := m.cell.GetId()

	pkt := packet.MobMoveEnd(m)
	m.world.BroadcastPacket(column, row, pkt)

	m.SetState(MobStateFind)
}

// findNewPath determines a new path for the mob to follow.
func (m *Mob) findNewPath(pos *context.Position) bool {
	newX := pos.CurrentX + rand.Intn(11) - 5
	newY := pos.CurrentY + rand.Intn(11) - 5

	if newX < 0 || newY < 0 {
		return false
	}

	if !m.world.IsMovable(newX, newY) {
		return false
	}

	pos.FinalX = newX
	pos.FinalY = newY

	return true
}

// findMovePath computes the path a mob should take to reach its
// destination using the A* algorithm.
func (m *Mob) findMovePath(pos *context.Position) ([]astar.Pather, float64) {
	start := WayPoint{
		X:     pos.InitialX,
		Y:     pos.InitialY,
		world: m.world,
	}

	end := WayPoint{
		X:     pos.FinalX,
		Y:     pos.FinalY,
		world: m.world,
	}

	path, distance, _ := astar.Path(start, end)

	return path, distance
}

// adjustCell updates the mob's cell location within the world.
func (m *Mob) adjustCell(pos *context.Position) {
	cell := m.world.getWorldCell(pos.CurrentX, pos.CurrentY)
	if cell == nil {
		log.Error("Invalid world cell")
		return
	}

	if cell == m.cell {
		return
	}

	m.cell.RemoveMob(m)

	m.mutex.Lock()
	m.cell = cell
	m.mutex.Unlock()

	cell.AddMob(m)
}
