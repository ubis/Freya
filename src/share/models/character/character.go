package character

import (
	"share/models/inventory"
	"time"
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
	Alz       uint64
	Nation    byte
	CurrentHP uint16
	MaxHP     uint16
	CurrentMP uint16
	MaxMP     uint16
	STR       uint32
	INT       uint32
	DEX       uint32
	SwordRank byte `db:"sword_rank"`
	MagicRank byte `db:"magic_rank"`
	Equipment inventory.Equipment
	Created   time.Time
}
