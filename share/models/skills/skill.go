package skills

type Skill struct {
	Id    uint16 `db:"skill"`
	Level byte
	Slot  uint16
}
