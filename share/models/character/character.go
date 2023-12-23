package character

import (
	"sync"
	"time"

	"github.com/ubis/Freya/share/models/inventory"
	"github.com/ubis/Freya/share/models/skills"
)

type ListReq struct {
	Account int32
}

type ListRes struct {
	List      []*Character
	LastId    int32
	SlotOrder int32
}

type CreateReq struct {
	Character
}

type CreateRes struct {
	Result byte
	Character
}

type DeleteReq struct {
	CharId int32
}

type DeleteRes struct {
	Result byte
}

type SetOrderReq struct {
	Account int32
	Order   int32
}

type SetOrderRes struct {
	Result bool
}

type Character struct {
	Id        int32
	Name      string
	Level     uint16
	World     byte
	X         byte
	Y         byte
	Style     Style
	LiveStyle int32 `db:"-"`
	Alz       uint64
	Nation    byte
	SwordRank byte   `db:"sword_rank"`
	MagicRank byte   `db:"magic_rank"`
	CurrentHP uint16 `db:"current_hp"`
	MaxHP     uint16 `db:"max_hp"`
	CurrentMP uint16 `db:"current_mp"`
	MaxMP     uint16 `db:"max_mp"`
	CurrentSP uint16 `db:"current_sp"`
	MaxSP     uint16 `db:"max_sp"`
	STR       uint32 `db:"str_stat"`
	INT       uint32 `db:"int_stat"`
	DEX       uint32 `db:"dex_stat"`
	PNT       uint32 `db:"pnt_stat"`
	Exp       uint64
	WarExp    uint64 `db:"war_exp"`
	Equipment inventory.Equipment
	Inventory *inventory.Inventory
	Skills    skills.SkillList
	Links     *skills.Links
	Created   time.Time

	// movement data
	BeginX int16 `db:"-"`
	BeginY int16 `db:"-"`
	EndX   int16 `db:"-"`
	EndY   int16 `db:"-"`

	mutex sync.RWMutex
}

func (c *Character) SetWorld(world byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.World = world
}

func (c *Character) GetWorld() byte {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.World
}

func (c *Character) SetLevel(level uint16) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Level = level
}

func (c *Character) GetLevel() uint16 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.Level
}

func (c *Character) GetHealth() (uint16, uint16) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.CurrentHP, c.MaxHP
}

func (c *Character) GetMana() (uint16, uint16) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.CurrentMP, c.MaxMP
}

func (c *Character) SetLiveStyle(style int32) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.LiveStyle = style
}

func (c *Character) GetStyle() (Style, int32) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.Style, c.LiveStyle
}

func (c *Character) SetPosition(x, y byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.X = x
	c.Y = y
}

func (c *Character) GetPosition() (byte, byte) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.X, c.Y
}

func (c *Character) SetMovement(sx, sy, dx, dy byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.BeginX = int16(sx)
	c.BeginY = int16(sy)
	c.EndX = int16(sx)
	c.EndY = int16(sy)
}

func (c *Character) GetMovement() (byte, byte, byte, byte) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return byte(c.BeginX), byte(c.BeginY), byte(c.EndX), byte(c.EndY)
}

type DataReq struct {
	Id int32
}

type DataRes struct {
	Inventory inventory.Inventory
	Skills    skills.SkillList
	Links     skills.Links
}
