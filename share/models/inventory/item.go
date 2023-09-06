package inventory

type Item struct {
	Kind    uint32
	Serials uint32
	Option  int32 `db:"opt"`
	Slot    uint16
	Expire  uint32
}
