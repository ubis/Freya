package encryption

import "github.com/ubis/Freya/share/log"

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
	var keyRand = KeyRand{}
	x.KeyTable = make([]uint32, RecvXorKeyNum*2)

	keyRand.Seed(RecvXorSeed)
	log.Debugf("1st XorSeed: %d", RecvXorSeed)

	for i := 0; i < RecvXorKeyNum; i++ {
		var low = keyRand.Rand()
		var high = keyRand.Rand()
		x.KeyTable[i] = (uint32)((low & 0xFFFF) | ((high & 0xFFFF) << 16))
	}
}

// Generates 2nd XorKeyTable
func (x *XorKeyTable) Generate2ndXorKeyTable(seed uint32) {
	var keyRand = KeyRand{}
	keyRand.Seed(seed)
	log.Debugf("2nd XorSeed: %d", seed)

	for i := RecvXorKeyNum; i < RecvXorKeyNum*2; i++ {
		var low = keyRand.Rand()
		var high = keyRand.Rand()
		x.KeyTable[i] = (uint32)((low & 0xFFFF) | ((high & 0xFFFF) << 16))
	}

	x.Seed2nd = seed
}
