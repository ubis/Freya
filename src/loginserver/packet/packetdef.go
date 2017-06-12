package packet

type S2C_HEADER struct {
    MagicKey 	uint16
    Length		uint16
    Opcode		uint16
}

type C2S_HEADER struct {
    MagicKey 	uint16
    Length		uint16
    Checksum	uint32
    Opcode 		uint16
}

type S2C_CONNECT2SVR struct {
    Header 		S2C_HEADER
    XorSeed		uint32
    AuthKey		uint32
    Index		uint16
    XorKeyIdx	uint16
}

type C2S_CONNECT2SVR struct {
    Header 		C2S_HEADER
    AuthKey 	uint32
}