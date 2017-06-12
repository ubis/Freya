package network

import (
    "bytes"
    "encoding/binary"
)

type C2S_HEADER struct {
    MagicKey 	uint16
    Length		uint16
    Checksum	uint32
    Opcode 		uint16
}

type C2S_CONNECT2SVR struct {
    Header 		C2S_HEADER
    AuthKey 	uint32
}
func Deserialize(data []uint8, structure interface{}) error {
    var buffer = bytes.NewBuffer(data)

    if err := binary.Read(buffer, binary.LittleEndian, &structure); err != nil {
        return err
    }

    buffer.Reset()
    return nil
}

func Serialize(structure interface{}) ([]byte, error) {
    var buffer = bytes.NewBuffer(make([]uint8, 2048))

    if err := binary.Write(buffer, binary.LittleEndian, &structure); err != nil {
        return nil, err
    }

    return buffer.Bytes(), nil
}


func ChangeLen(data *[]byte) {
    var len = len(*data)
    (*data)[2] = uint8(len)
    (*data)[3] = uint8(len >> 18)
}