package network

import (
    "net"
    "share/event"
    "share/encryption"
    "share/models/session"
)

// max buffer size
const MAX_RECV_BUFFER_SIZE = 4096

type Session struct {
    socket      net.Conn
    buffer      []byte

    Encryption  encryption.Encryption
    UserIdx     uint16
    AuthKey     uint32
    Data        *session.Data
}

/*
    Starts session goroutine
    @param  key encryption XorKeyTable
 */
func (s *Session) Start(table encryption.XorKeyTable) {
    // create new receiving buffer
    s.buffer = make([]byte, MAX_RECV_BUFFER_SIZE)

    // init encryption
    s.Encryption = encryption.Encryption{}
    s.Encryption.Init(&table)

    s.Data = &session.Data{}

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
            log.Error("Error decrypting: " + error.Error())
            event.Trigger(event.ClientDisconnectEvent, s)
            break
        }

        // create new packet reader
        var reader = NewReader(data)

        // create new packet event argument
        var arg = &PacketArgs{
            s,
            int(reader.Size),
            int(reader.Type),
            reader,
        }

        // trigger packet received event
        event.Trigger(event.PacketReceiveEvent, arg)
    }
}

/*
    Sends specified data to the client
    @param  writer  a pointer to Writer so that byte array of data could be received from it
 */
func (s *Session) Send(writer *Writer) {
    // encrypt data
    var encrypt, err = s.Encryption.Encrypt(writer.Finalize())
    if err != nil {
        log.Error("Error encrypting packet: " + err.Error())
        return
    }

    // send it...
    var length, err2 = s.socket.Write(encrypt)
    if err2 != nil {
        log.Error("Error sending packet: " + err2.Error())
        return
    }

    // create new packet event argument
    var arg = &PacketArgs{
        s,
        length,
        writer.Type,
        nil,
    }

    // trigger packet sent event
    event.Trigger(event.PacketSendEvent, arg)
}

/*
    Returns session's remote endpoint
    @return remote endpoint
 */
func (s *Session) GetEndPnt() string {
    return s.socket.RemoteAddr().String()
}

// Closes session socket
func (s *Session) Close() {
    s.socket.Close()
}