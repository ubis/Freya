package packet

import (
    "bytes"
    "share/network"
    "share/models/account"
    "loginserver/rsa"
)

// PreServerEnvRequest Packet
func PreServerEnvRequest(session *network.Session, reader *network.Reader) {
    var packet = network.NewWriter(PRE_SERVER_ENV_REQUEST)

    packet.WriteBytes(make([]byte, 4113))

    session.Send(packet)
}

// PublicKey Packet
func PublicKey(session *network.Session, reader *network.Reader) {
    var rsa = g_ServerSettings.RSA
    var key = rsa.PublicKey

    var packet = network.NewWriter(PUBLIC_KEY)
    packet.WriteByte(0x01)
    packet.WriteUint16(len(key))
    packet.WriteBytes(key[:])

    session.Send(packet)
}

// AuthAccount Packet
func AuthAccount(session *network.Session, reader *network.Reader) {
    if session.Data.Verified != true {
        log.Errorf("Session version is not verified! Src: %s, UserIdx: %d",
            session.GetEndPnt(),
            session.UserIdx,
        )
        return
    }

    reader.ReadUint16()

    var loginData = reader.ReadBytes(rsa.RSA_LOGIN_LENGTH)

    var rsa = g_ServerSettings.RSA
    var data, err = rsa.Decrypt(loginData[:])
    if err != nil {
        log.Errorf("%s; Src: %s, UserIdx: %d",
            err.Error(),
            session.GetEndPnt(),
            session.UserIdx,
        )
        return
    }

    var userId = string(bytes.Trim(data[:32], "\x00"))
    var userPw = string(bytes.Trim(data[32:], "\x00"))

    var packet = network.NewWriter(AUTHACCOUNT)

    if userId == "root" && userPw == "root" {
        packet.WriteByte(account.Normal)
    } else {
        packet.WriteByte(account.Incorrect)
    }

    packet.WriteInt32(0x00)   // account id
    packet.WriteByte(0x00)
    packet.WriteByte(0x00)
    packet.WriteByte(0x00)    // server count
    packet.WriteInt32(0x00); packet.WriteInt32(0x00)
    packet.WriteInt32(0x00)   // premium service id
    packet.WriteInt32(0x00)   // premium service expire date
    packet.WriteByte(0x00)
    packet.WriteByte(0x00)
    packet.WriteInt32(0x00); packet.WriteInt16(0x00); packet.WriteByte(0x00)
    packet.WriteInt32(0x00)   // language
    packet.WriteString("ASDASDDADSSADSAASDASASD")


    log.Infof("Logging in with: %s and %s", userId, userPw)
    session.Send(packet)
}