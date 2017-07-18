package packet

import (
    "share/network"
)

// GetMyChartr Packet
func GetMyChartr(session *network.Session, reader *network.Reader) {
    var packet = network.NewWriter(GETMYCHARTR)
    packet.WriteInt32(0x00) // subpassword exists
    packet.WriteBytes(make([]byte, 9))
    packet.WriteInt32(0x00) // selected character id
    packet.WriteInt32(0x00) // slot order

    session.Send(packet)
}