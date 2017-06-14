package network

import (
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
)

/*
    Attempts to deserialize binary array to specified struct
    @param  data        source data array
    @param  structure   specified struct
    @return error if deserialization failed
 */
func Deserialize(data []uint8, structure interface{}) error {
    // first serialize empty structure
    var buffer = new(bytes.Buffer)
    binary.Write(buffer, binary.LittleEndian, structure);

    // check data length
    if len(buffer.Bytes()) != len(data) {
        return errors.New(
            fmt.Sprintf("Invalid len (Detected: %d, required: %d)",
                len(data),
                len(buffer.Bytes()),
            ),
        )
    }

    // now try to deserialize real data...
    buffer = bytes.NewBuffer(data)
    if err := binary.Read(buffer, binary.LittleEndian, structure); err != nil {
        return errors.New(
            fmt.Sprintf(
            "Cannot deserialize (Array len: %d, Detected len: %d, Type: %d)\n%s",
                len(data),
                int(data[2] + (data[3] >> 16)),
                int(data[8] + (data[9] >> 16)),
                err.Error(),
            ),
        )
    }

    return nil
}

/*
    Attempts to serialize specified struct into binary array
    @param  structure   specified struct
    @return binary array with serialized data and error if serialization failed
 */
func Serialize(structure interface{}) ([]byte, error) {
    var buffer = new(bytes.Buffer)
    if err := binary.Write(buffer, binary.LittleEndian, structure); err != nil {
        return nil, errors.New("Error serializing packet: " + err.Error())
    }

    // update packet length
    // since Golang adds padding to struct's,
    // we need to update data manually
    var data = buffer.Bytes()
    var len  = len(data)

    if len >= 4 {
        data[2] = uint8(len)
        data[3] = uint8(len >> 18)
    }

    return data, nil
}