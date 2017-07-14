package packet

import (
    "bytes"
    "share/network"
    "share/rpc"
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

    log.Infof("AuthRequest with id `%s`", userId)

    var r = account.AuthResponse{Status: account.None}
    err = g_RPCHandler.Call(rpc.AuthCheck,
        account.AuthRequest{userId, userPw}, &r)

    if err != nil {
        r.Status = account.OutOfService
    }

    var packet = network.NewWriter(AUTHACCOUNT)

    packet.WriteByte(r.Status)
    packet.WriteInt32(r.Id)
    packet.WriteInt16(0x00)
    packet.WriteByte(0x00)  // server count
    packet.WriteInt64(0x00)
    packet.WriteInt32(0x00) // premium service id
    packet.WriteInt32(0x00) // premium service expire date
    packet.WriteByte(0x00)
    packet.WriteByte(0x00)  // subpw exist
    packet.WriteBytes(make([]byte, 7))
    packet.WriteInt32(0x00) // language
    packet.WriteString(r.AuthKey + "\x00")

    if r.Status == account.Normal {
        log.Infof("User `%s` succesfully logged in.", userId)
    }

    session.Send(packet)
}