package packet

import (
    "share/network"
    "time"
)

// Connect2Svr Packet
func Connect2Svr(session *network.Session, data []uint8) {
    var c2s = C2S_CONNECT2SVR{}

    if err := network.Deserialize(data, &c2s); err != nil {
        log.Errorf("%s; Src: %s, UserIdx: %d",
            err.Error(),
            session.GetEndPnt(),
            session.UserIdx,
        )
        return
    }

    session.AuthKey = uint32(time.Now().Unix())

    var s2c = S2C_CONNECT2SVR{
        S2C_HEADER{MAGIC_KEY, 0, CONNECT2SVR},
        session.Encryption.Key.Seed2nd,
        session.AuthKey,
        session.UserIdx,
        uint16(session.Encryption.RecvXorKeyIdx),
    }

    var recvData, err = network.Serialize(&s2c)

    if err != nil {
        log.Error(err.Error())
        return
    }

    session.Send(recvData)
}