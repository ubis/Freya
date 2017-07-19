package encryption

import (
    "encoding/binary"
    "errors"
    "fmt"
    "math/rand"
)

var MASK = []uint32 { 0xFFFFFFFF, 0xFFFFFF00, 0xFFFF0000, 0xFF000000 };

// Encryption structure
type Encryption struct {
    recvXorKey              uint32
    isFirstPacket           bool
    xorKeyTableBaseMultiple uint32

    Key                     *XorKeyTable
    RecvXorKeyIdx           uint32
}

/*
    Initializes encryption
    @param  key     or key table
 */
func (e *Encryption) Init(key *XorKeyTable) {
    e.isFirstPacket           = true
    e.xorKeyTableBaseMultiple = 1
    e.Key                     = key
    e.RecvXorKeyIdx           = uint32(rand.Intn(RecvXorKeyNum))
}

/*
    Returns packet size
    @param  data    data array from which size will be taken
    @return packet size
 */
func (e *Encryption) GetPacketSize(data []uint8) int {
    if !e.isFirstPacket {
        var header = binary.LittleEndian.Uint32(data) ^ e.recvXorKey
        return int(header >> 16)
    }

    return Connect2SvrSize
}

/*
    Encrypts data
    @param  data    data array to be encrypted
    @return encrypted data or error, if any
 */
func (e *Encryption) Encrypt(data []uint8) ([]uint8, error) {
    var iLen = len(data)

    if iLen < S2CHeaderSize {
        // packet is smaller than header...
        return nil, errors.New(
            fmt.Sprintf(
                "Packet length is smaller than required! Detected: %d, required: %d",
                iLen, S2CHeaderSize,
            ),
        )
    }

    var pPayload = make([]uint8, iLen)

    var xorKey uint32
    var xorNum uint32

    // header
    xorNum  = binary.LittleEndian.Uint32(data) ^ SendXorKey
    binary.LittleEndian.PutUint32(pPayload, xorNum)
    xorKey = e.Key.KeyTable[xorNum & RecvXorKeyNumMask]

    var payloadLen = iLen - (S2CHeaderSize - MainCMDSize);
    var j int32    = 4

    for i := 0; i < payloadLen / 4; i++ {
        xorNum = binary.LittleEndian.Uint32(data[j:j + 4]) ^ xorKey;
        binary.LittleEndian.PutUint32(pPayload[j:j + 4], xorNum)
        xorKey = e.Key.KeyTable[xorNum & RecvXorKeyNumMask]
        j += 4
    }

    var remainLen  = int32(iLen % 4)
    var remainData = make([]uint8, 4)

    copy(remainData, data[j:j + remainLen])
    var r = binary.LittleEndian.Uint32(remainData)
    var m = MASK[remainLen];

    xorNum = r ^ xorKey;
    r      = xorNum ^ (xorKey & m);

    binary.LittleEndian.PutUint32(remainData, r)
    copy(pPayload[j:j + remainLen], remainData[0:remainLen])

    return pPayload, nil
}

/*
    Decrypts data
    @param  data    data array to be decrypted
    @return decrypted data or error, if any
 */
func (e *Encryption) Decrypt(data []uint8) ([]uint8, error) {
    var packetLen int
    var header    uint32

    if !e.isFirstPacket {
        header = binary.LittleEndian.Uint32(data) ^ e.recvXorKey

        if uint16(header) != MagicKey {
            // invalid magic key
            return nil, errors.New(
                fmt.Sprintf("MagicKey is invalid! Detected: %d, required: %d",
                    uint16(header), MagicKey,
                ),
            )
        }

        packetLen = int(header >> 16)
    } else {
        header       = binary.LittleEndian.Uint32(data)
        e.recvXorKey = header ^ (MagicKey | (Connect2SvrSize << 16))

        packetLen = Connect2SvrSize

        if uint16(header ^ e.recvXorKey) != MagicKey {
            // invalid magic key
            return nil, errors.New(
                fmt.Sprintf("MagicKey is invalid! Detected: %d, required: %d",
                    uint16(header ^ e.recvXorKey), MagicKey,
                ),
            )
        }

        if int(header ^ e.recvXorKey) >> 16 != packetLen || packetLen < C2SHeaderSize  {
            // invalid packet size
            return nil, errors.New(
                fmt.Sprintf("Packet size is invalid! Detected: %d, required: %d",
                    int(header ^ e.recvXorKey) >> 16, packetLen,
                ),
            )
        }
    }

    var pHeader  = binary.LittleEndian.Uint32(data)
    //var pCheckSum = binary.LittleEndian.Uint32(data[4:8])
    var pPayload = make([]uint8, packetLen)
    var xorKey   = e.recvXorKey
    var xorNum uint32

    // decrypt header
    xorNum  = pHeader;
    pHeader = xorNum ^ xorKey;
    xorKey  = e.Key.KeyTable[(xorNum & RecvXorKeyNumMask) * e.xorKeyTableBaseMultiple];
    binary.LittleEndian.PutUint32(pPayload, pHeader)

    // decrypt body
    var payloadLen = packetLen - (C2SHeaderSize - MainCMDSize);
    var j int32    = 8

    for i := 0; i < payloadLen / 4; i++ {
        xorNum = binary.LittleEndian.Uint32(data[j:j + 4]);
        binary.LittleEndian.PutUint32(pPayload[j:j + 4], xorNum ^ xorKey)
        xorKey = e.Key.KeyTable[(xorNum & RecvXorKeyNumMask) * e.xorKeyTableBaseMultiple];
        j += 4
    }

    var remainLen  = int32(packetLen % 4)
    var remainData = make([]uint8, 4)

    copy(remainData, data[j:j + remainLen])
    var r = binary.LittleEndian.Uint32(remainData)
    var m = MASK[remainLen];

    xorNum = r;
    r      = (xorNum ^ xorKey) ^ (xorKey & m);

    binary.LittleEndian.PutUint32(remainData, r)
    copy(pPayload[j:j + remainLen], remainData[0:remainLen])

    if e.isFirstPacket {
        e.isFirstPacket           = false
        e.xorKeyTableBaseMultiple = 2;
    }

    e.recvXorKey = e.Key.KeyTable[(e.RecvXorKeyIdx) * e.xorKeyTableBaseMultiple]
    e.RecvXorKeyIdx++

    if e.RecvXorKeyIdx >= RecvXorKeyNum {
        e.RecvXorKeyIdx = 0
    }

    return pPayload, nil
}