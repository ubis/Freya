package encryption

// XorKeyNum constants
const (
    RecvXorKeyNum     = 0x4000
    RecvXorKeyNumMask = 0x00003FFF
    RecvXorSeed       = (0x8F54C37B | 1)
    Recv2ndXorSeed    = (0x34BC821A | 1)
    SendXorKey        = 0x7AB38CF1
)

// Constant sizes
const (
    S2CHeaderSize     = 0x06
    C2SHeaderSize     = 0x0A
    MainCMDSize       = 0x02
    Connect2SvrSize   = 0x0E
)

// Packet constants
const (
    MagicKey          = 0xB7E2
)