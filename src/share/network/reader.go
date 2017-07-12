package network

// Packet Reader struct
type Reader struct {
    buffer  []byte
    index    int

    MagicKey uint16
    Size     uint16
    CheckSum uint32
    Type     uint16
}

/*
    Attempts to create a new packet reader and read packet header
    @param  buffer  packet buffer to be readed
    @return pointer to the Reader
 */
func NewReader(buffer []byte) *Reader {
    var r = &Reader{}

    r.buffer = buffer
    r.index   = 0

    r.MagicKey  = r.ReadUint16()
    r.Size      = r.ReadUint16()
    r.CheckSum  = r.ReadUint32()
    r.Type      = r.ReadUint16()

    return r
}

/*
    Attempts to read an signed byte
    @return signed byte
 */
func (r *Reader) ReadSByte() int8 {
    if len(r.buffer) <= r.index {
        var err = "Error reading sbyte: buffer is too small!"
        log.Panic(err)
    }

    var data = int8(r.buffer[r.index])
    r.index ++

    return data
}

/*
    Attempts to read an unsigned byte
    @return unsigned byte
 */
func (r *Reader) ReadByte() byte {
    if len(r.buffer) <= r.index {
        var err = "Error reading byte: buffer is too small!"
        log.Panic(err)
    }

    var data = r.buffer[r.index]
    r.index ++

    return data
}

/*
    Attempts to read an signed 16-bit integer
    @return signed 16-bit integer
 */
func (r *Reader) ReadInt16() int16 {
    if len(r.buffer) <= r.index + 1 {
        var err = "Error reading int16: buffer is too small!"
        log.Panic(err)
    }

    var data = int16(r.buffer[r.index])
    data     += int16(r.buffer[r.index + 1]) << 8
    r.index += 2

    return data
}

/*
    Attempts to read an unsigned 16-bit integer
    @return unsigned 16-bit integer
 */
func (r *Reader) ReadUint16() uint16 {
    if len(r.buffer) <= r.index + 1 {
        var err = "Error reading uint16: buffer is too small!"
        log.Panic(err)
    }

    var data = uint16(r.buffer[r.index])
    data     += uint16(r.buffer[r.index + 1]) << 8
    r.index += 2

    return data
}

/*
    Attempts to read an signed 32-bit integer
    @return signed 32-bit integer
 */
func (r *Reader) ReadInt32() int32 {
    if len(r.buffer) <= r.index + 3 {
        var err = "Error reading int32: buffer is too small!"
        log.Panic(err)
    }

    var data = int32(r.buffer[r.index])
    data     += int32(r.buffer[r.index + 1]) << 8
    data     += int32(r.buffer[r.index + 2]) << 16
    data     += int32(r.buffer[r.index + 3]) << 24
    r.index += 4

    return data
}

/*
    Attempts to read an unsigned 32-bit integer
    @return unsigned 32-bit integer
 */
func (r *Reader) ReadUint32() uint32 {
    if len(r.buffer) <= r.index + 3 {
        var err = "Error reading uint32: buffer is too small!"
        log.Panic(err)
    }

    var data = uint32(r.buffer[r.index])
    data     += uint32(r.buffer[r.index + 1]) << 8
    data     += uint32(r.buffer[r.index + 2]) << 16
    data     += uint32(r.buffer[r.index + 3]) << 24
    r.index += 4

    return data
}

/*
    Attempts to read an signed 64-bit integer
    @return signed 64-bit integer
 */
func (r *Reader) ReadInt64() int64 {
    if len(r.buffer) <= r.index + 7 {
        var err = "Error reading int64: buffer is too small!"
        log.Panic(err)
    }

    var data = int64(r.buffer[r.index])
    data     += int64(r.buffer[r.index + 1]) << 8
    data     += int64(r.buffer[r.index + 2]) << 16
    data     += int64(r.buffer[r.index + 3]) << 24
    data     += int64(r.buffer[r.index + 4]) << 32
    data     += int64(r.buffer[r.index + 5]) << 40
    data     += int64(r.buffer[r.index + 6]) << 48
    data     += int64(r.buffer[r.index + 7]) << 56
    r.index += 8

    return data
}

/*
    Attempts to read an unsigned 64-bit integer
    @return unsigned 64-bit integer
 */
func (r *Reader) ReadUint64() uint64 {
    if len(r.buffer) <= r.index + 7 {
        var err = "Error reading uint64: buffer is too small!"
        log.Panic(err)
    }

    var data = uint64(r.buffer[r.index])
    data     += uint64(r.buffer[r.index + 1]) << 8
    data     += uint64(r.buffer[r.index + 2]) << 16
    data     += uint64(r.buffer[r.index + 3]) << 24
    data     += uint64(r.buffer[r.index + 4]) << 32
    data     += uint64(r.buffer[r.index + 5]) << 40
    data     += uint64(r.buffer[r.index + 6]) << 48
    data     += uint64(r.buffer[r.index + 7]) << 56
    r.index += 8

    return data
}

/*
    Attempts to read a string with given length
    @param length   string's length
    @return string
 */
func (r *Reader) ReadString(length int) string {
    if len(r.buffer) <= r.index + length - 1 {
        var err = "Error reading string: buffer is too small!"
        log.Panic(err)
    }

    var data = r.buffer[r.index:r.index + length]
    r.index += length

    return string(data)
}

/*
    Attempts to read an byte array with given length
    @param length   byte array length
    @return byte array
 */
func (r *Reader) ReadBytes(length int) []byte {
    if len(r.buffer) <= r.index + length - 1 {
        var err = "Error reading []byte: buffer is too small!"
        log.Panic(err)
    }

    var data = r.buffer[r.index:r.index + length]
    r.index += length

    return data
}