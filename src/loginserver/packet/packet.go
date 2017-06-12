package packet

import (
    "share/network"
    "share/encryption"
)

func Connect2Svr(session *network.Session, data []uint8) {
    var c2s = C2S_CONNECT2SVR{}

    if err := network.Deserialize(data, &c2s); err != nil {
        log.Errorf(
            "Cannot deserialize(Len: %d, Type: %d, Src: %s, UserIdx: %d)\n%s",
            len(data),
            0,
            session.GetEndPnt(),
            session.UserIdx,
            err.Error(),
        )
        return
    }

    var s2c = S2C_CONNECT2SVR{
        S2C_HEADER{encryption.MagicKey, 0, CONNECT2SVR},
        session.Encryption.Key.Seed2nd,
        0x11223344,
        session.UserIdx,
        uint16(session.Encryption.RecvXorKeyIdx),
    }

    var recvData, err = network.Serialize(s2c)

    if err != nil {
        log.Error("Error serializing packet: " + err.Error())
        return
    }

    session.Send(recvData)
}
