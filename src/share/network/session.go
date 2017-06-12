package network

import (
    "net"
    "share/event"
    "share/encryption"
)

const MAX_RECV_BUFFER_SIZE = 4096

type Session struct {
    socket      net.Conn
    buffer      []byte

    Encryption  encryption.Encryption
    UserIdx     uint16
}

func (s *Session) Start(key *encryption.XorKeyTable) {
    // create new receiving buffer
    s.buffer    = make([]byte, MAX_RECV_BUFFER_SIZE)
    s.Encryption = encryption.Encryption{}

    // init encryption
    s.Encryption.Init(key)

    for {
        // read data
        var _, err = s.socket.Read(s.buffer)

        if err != nil {
            log.Error("Error reading: " + err.Error())
            event.Trigger(event.ClientDisconnectEvent, s)
            break
        }

        // attempt to decrypt packet
        var data, error = s.Encryption.Decrypt(s.buffer)

        if error != nil {
            log.Error("Error decrypting: " + err.Error())
            event.Trigger(event.ClientDisconnectEvent, s)
            break
        }

        // create new packet event argument
        var arg = &PacketArgs{
            s,
            len(data),
            int(data[8] + (data[9] >> 16)),
            &data,
        }

        // trigger packet received event
        event.Trigger(event.PacketReceiveEvent, arg)
    }
}

func (s *Session) Send(data []uint8) {
    var length = 0

    var encrypt, err = s.Encryption.Encrypt(data)
    if err != nil {
        log.Error("Error encrypting packet: " + err.Error())
        return
    }

    length, err = s.socket.Write(encrypt)
    if err != nil {
        log.Error("Error sending packet: " + err.Error())
        return
    }

    // create new packet event argument
    var arg = &PacketArgs{
        s,
        length,
        int(data[4] + (data[5] >> 16)),
        nil,
    }

    // trigger packet received event
    event.Trigger(event.PacketSendEvent, arg)
}

func (s *Session) GetEndPnt() string {
    return s.socket.RemoteAddr().String()
}