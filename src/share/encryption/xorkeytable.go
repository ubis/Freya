package encryption

import "share/logger"

var log = logger.Instance()

// XorKeyTable structure
type XorKeyTable struct {
    KeyTable []uint32
    Seed2nd  uint32
}

// Initializes XorKeyTable
func (x *XorKeyTable) Init() {
    log.Info("Generating XorKeyTable...")
    x.XorKeyTable()
    x.Generate2ndXorKeyTable(Recv2ndXorSeed)
}

// Generates 1st XorKeyTable
func (x *XorKeyTable) XorKeyTable() {
    var keyRand   = KeyRand{}
    x.KeyTable = make([]uint32, RecvXorKeyNum * 2)

    keyRand.Seed(RecvXorSeed)
    log.Debugf("1st XorSeed: %d", RecvXorSeed)

    for i := 0; i < RecvXorKeyNum; i++ {
        var low  = keyRand.Rand()
        var high = keyRand.Rand()
        x.KeyTable[i] = (uint32)((low & 0xFFFF) | ((high & 0xFFFF) << 16))
    }
}

/*
    Generates 2nd XorKeyTable
    @param  seed    seed value to set and use
 */
func (x *XorKeyTable) Generate2ndXorKeyTable(seed uint32) {
    var keyRand = KeyRand{}
    keyRand.Seed(seed)
    log.Debugf("2nd XorSeed: %d", seed)

    for i := RecvXorKeyNum; i < RecvXorKeyNum * 2; i++ {
        var low  = keyRand.Rand()
        var high = keyRand.Rand()
        x.KeyTable[i] = (uint32)((low & 0xFFFF) | ((high & 0xFFFF) << 16))
    }

    x.Seed2nd = seed
}