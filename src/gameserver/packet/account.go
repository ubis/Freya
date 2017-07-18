package packet

import (
    "share/network"
)

// ChargeInfo Packet
func ChargeInfo(session *network.Session, reader *network.Reader) {
    var packet = network.NewWriter(CHARGEINFO)
    packet.WriteInt32(0x00)
    packet.WriteInt32(0x00)     // service kind
    packet.WriteUint32(0x00)    // service expire

    session.Send(packet)
}
