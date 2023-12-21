package character

import (
	"time"

	"github.com/ubis/Freya/share/models/inventory"
	"github.com/ubis/Freya/share/models/skills"
)

type ListReq struct {
	Account int32
	Server  byte
}

type ListRes struct {
	List      []Character
	LastId    int32
	SlotOrder int32
}

type CreateReq struct {
	Server byte
	Character
}

type CreateRes struct {
	Result byte
	Character
}

type DeleteReq struct {
	Server byte
	CharId int32
}

type DeleteRes struct {
	Result byte
}

type SetOrderReq struct {
	Server  byte
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
}

type DataReq struct {
	Server byte
	Id     int32
}

type DataRes struct {
	Inventory inventory.Inventory
	Skills    skills.SkillList
	Links     skills.Links
}
