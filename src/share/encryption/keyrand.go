package encryption

// KeyRand structure
type KeyRand struct {
    holdRand uint32
}

/*
    Sets initial seed value
    @param  seed    initial seed value to set
 */
func (k *KeyRand) Seed(seed uint32) {
    k.holdRand = seed
}

// Returns random integer
func (k *KeyRand) Rand() (int64) {
    // some hardcoded func...
    k.holdRand = k.holdRand * 49723125 + 21403831
    return (((((int64)(k.holdRand) >> 16) * 41894339 + 11741117) >> 16) & 0xFFFF)
}