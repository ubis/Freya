package inventory

type Item struct {
	Kind    uint32
	Serials uint32
	Option  uint32 `db:"opt"`
	Slot    uint16
	Expire  uint32
}
