package encryption

type KeyRand struct {
    holdRand uint32
}

// Sets initial seed value
func (k *KeyRand) Seed(seed uint32) {
    k.holdRand = seed
}

// Returns random int64 number
func (k *KeyRand) Rand() (int64) {
    // some hardcoded function from ep3 sources...
    // probably should be changed with other(better) algorithm
    k.holdRand = k.holdRand * 49723125 + 21403831
    return (((((int64)(k.holdRand) >> 16) * 41894339 + 11741117) >> 16) & 0xFFFF)
}