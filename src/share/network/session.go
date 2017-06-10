package network

import (
    "net"
    "share/event"
)

const MAX_RECV_BUFFER_SIZE = 4096

type Session struct {
    socket      net.Conn
    buffer      []byte

    UserIdx     uint16
}

func (s *Session) GetEndPnt() string {
    return s.socket.RemoteAddr().String()
}

func (s *Session) Start() {
    // create new receiving buffer
    s.buffer = make([]byte, MAX_RECV_BUFFER_SIZE)

    for {
        // read data
        reqLen, err := s.socket.Read(s.buffer)

        if err != nil {
            log.Error("Error reading: " + err.Error())
            event.Trigger(event.ClientDisconnectEvent, s)
            break
        }


        event.Trigger(event.ClientConnectEvent, nil)
        log.Info("Received Packet Size:", reqLen)

        /*var t = s.enc.Decrypt(buf)
        pk.Handle(t)

        for  i := 0; i < reqLen; i ++ {
            fmt.Print(fmt.Sprintf("%02x ", t[i]))
        }

        fmt.Print("\n")*/
    }
}