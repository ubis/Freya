package network

import (
	"reflect"

	"github.com/ubis/Freya/share/encryption"
	"github.com/ubis/Freya/share/log"
)

const DEFAULT_BUFFER_SIZE = 1024

type Writer struct {
	buffer []byte
	index  int

	Type int
}

// Attempts to create a new packet writer and write initial packet header
func NewWriter(code uint16, magic ...uint16) *Writer {
	var w = &Writer{}

	w.buffer = make([]byte, DEFAULT_BUFFER_SIZE)
	w.index = 0

	if magic == nil {
		// grab magic key from encryption package
		w.WriteUint16(encryption.MagicKey)
	} else {
		// write custom magic key
		w.WriteUint16(magic[0])
	}

	w.WriteUint16(0x00) // size
	w.WriteUint16(code) // packet type

	w.Type = int(code)

	return w
}

// Checks buffer length and if it's too small, it will resize it
func (w *Writer) checkLength(length int) {
	if len(w.buffer) < w.index+length {
		// resize...
		var tmp = make([]byte, len(w.buffer)+DEFAULT_BUFFER_SIZE)
		copy(tmp[:len(w.buffer)], w.buffer)
		w.buffer = tmp

		// recursion
		w.checkLength(length)
	}
}

// Attempts to read specified interface type and serializes it into byte array.
// It has a length parameter, which tells the required size of interface type.
// If interface type length is smaller or higher than required, the correct
// length will be returned
func (w *Writer) getType(obj interface{}, length int) []byte {
	// check length
	w.checkLength(length)
	var tmp = make([]byte, 8)

	switch objType := obj.(type) {
	case int8:
		tmp[0] = byte(objType)
	case uint8:
		tmp[0] = byte(objType)
	case int16:
		tmp[0] = byte(objType)
		tmp[1] = byte(objType >> 8)
	case uint16:
		tmp[0] = byte(objType)
		tmp[1] = byte(objType >> 8)
	case int32:
		tmp[0] = byte(objType)
		tmp[1] = byte(objType >> 8)
		tmp[2] = byte(objType >> 16)
		tmp[3] = byte(objType >> 24)
	case uint32:
		tmp[0] = byte(objType)
		tmp[1] = byte(objType >> 8)
		tmp[2] = byte(objType >> 16)
		tmp[3] = byte(objType >> 24)
	case int:
		tmp[0] = byte(objType)
		tmp[1] = byte(objType >> 8)
		tmp[2] = byte(objType >> 16)
		tmp[3] = byte(objType >> 24)
	case int64:
		tmp[0] = byte(objType)
		tmp[1] = byte(objType >> 8)
		tmp[2] = byte(objType >> 16)
		tmp[3] = byte(objType >> 24)
		tmp[4] = byte(objType >> 32)
		tmp[5] = byte(objType >> 40)
		tmp[6] = byte(objType >> 48)
		tmp[7] = byte(objType >> 56)
	case uint64:
		tmp[0] = byte(objType)
		tmp[1] = byte(objType >> 8)
		tmp[2] = byte(objType >> 16)
		tmp[3] = byte(objType >> 24)
		tmp[4] = byte(objType >> 32)
		tmp[5] = byte(objType >> 40)
		tmp[6] = byte(objType >> 48)
		tmp[7] = byte(objType >> 56)
	default:
		log.Error("Unknown data type:", reflect.TypeOf(obj))
		return nil
	}

	return tmp[:length]
}

// Writes a bool
func (w *Writer) WriteBool(data bool) {
	value := 0
	if data {
		value = 1
	}

	w.buffer[w.index] = byte(value)
	w.index++
}

// Writes an signed byte
func (w *Writer) WriteSbyte(data interface{}) {
	var t = w.getType(data, 1)

	w.buffer[w.index] = byte(t[0])
	w.index++
}

// Writes an unsigned byte
func (w *Writer) WriteByte(data interface{}) {
	var t = w.getType(data, 1)

	w.buffer[w.index] = byte(t[0])
	w.index++
}

// Writes an signed 16-bit integer
func (w *Writer) WriteInt16(data interface{}) {
	var t = w.getType(data, 2)

	w.buffer[w.index] = byte(t[0])
	w.buffer[w.index+1] = byte(t[1])
	w.index += 2
}

// Writes an unsigned 16-bit integer
func (w *Writer) WriteUint16(data interface{}) {
	var t = w.getType(data, 2)

	w.buffer[w.index] = byte(t[0])
	w.buffer[w.index+1] = byte(t[1])
	w.index += 2
}

// Writes an signed 32-bit integer
func (w *Writer) WriteInt32(data interface{}) {
	var t = w.getType(data, 4)

	w.buffer[w.index] = byte(t[0])
	w.buffer[w.index+1] = byte(t[1])
	w.buffer[w.index+2] = byte(t[2])
	w.buffer[w.index+3] = byte(t[3])
	w.index += 4
}

// Writes an unsigned 32-bit integer
func (w *Writer) WriteUint32(data interface{}) {
	var t = w.getType(data, 4)

	w.buffer[w.index] = byte(t[0])
	w.buffer[w.index+1] = byte(t[1])
	w.buffer[w.index+2] = byte(t[2])
	w.buffer[w.index+3] = byte(t[3])
	w.index += 4
}

// Writes an signed 64-bit integer
func (w *Writer) WriteInt64(data interface{}) {
	var t = w.getType(data, 8)

	w.buffer[w.index] = byte(t[0])
	w.buffer[w.index+1] = byte(t[1])
	w.buffer[w.index+2] = byte(t[2])
	w.buffer[w.index+3] = byte(t[3])
	w.buffer[w.index+4] = byte(t[4])
	w.buffer[w.index+5] = byte(t[5])
	w.buffer[w.index+6] = byte(t[6])
	w.buffer[w.index+7] = byte(t[7])
	w.index += 8
}

// Writes an unsigned 64-bit integer
func (w *Writer) WriteUint64(data interface{}) {
	var t = w.getType(data, 8)

	w.buffer[w.index] = byte(t[0])
	w.buffer[w.index+1] = byte(t[1])
	w.buffer[w.index+2] = byte(t[2])
	w.buffer[w.index+3] = byte(t[3])
	w.buffer[w.index+4] = byte(t[4])
	w.buffer[w.index+5] = byte(t[5])
	w.buffer[w.index+6] = byte(t[6])
	w.buffer[w.index+7] = byte(t[7])
	w.index += 8
}

// Writes a string
func (w *Writer) WriteString(data string) {
	// check length
	w.checkLength(len(data))

	var bytes = []byte(data)
	copy(w.buffer[w.index:], bytes)
	w.index += len(bytes)
}

// Writes an byte array
func (w *Writer) WriteBytes(data []byte) {
	// check length
	w.checkLength(len(data))

	copy(w.buffer[w.index:], data)
	w.index += len(data)
}

/*
Updates packet length and returns byte array
@return byte array of packet
*/
func (w *Writer) Finalize() []byte {
	// update size
	var length = w.index
	w.buffer[2] = byte(length)
	w.buffer[3] = byte(length >> 8)

	return w.buffer[:length]
}
