package packet

const (
	MagicKey = 0xB7E2
)

type C2SHeader struct {
	MagicCode uint16
	Length    uint16
	Checksum  uint32
	Opcode    uint16
}

type S2CHeader struct {
	MagicCode uint16
	Length    uint16
	Opcode    uint16
}

type C2SConnect2Svr struct {
	C2SHeader

	AuthKey uint32
}

type S2CConnect2Svr struct {
	S2CHeader

	XorSeed   uint32
	AuthKey   uint32
	UserIdx   uint16
	XorKeyIdx uint16
}
