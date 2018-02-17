package network

import (
	"fmt"
	"io"
	"net"
	"share/encryption"
	"share/event"
	"share/log"
	"share/models/character"
	"share/models/subpasswd"
	"strings"
)

type Session struct {
	socket  net.Conn
	buffer  []byte
	txQueue chan *Writer

	Encryption encryption.Encryption
	UserIdx    uint16
	AuthKey    uint32
	Connected  bool
	Data       struct {
		AccountId     int32 // database account id
		Verified      bool  // version verification
		LoggedIn      bool  // auth verification
		CharVerified  bool  // character delete password verification
		SubPassword   *subpasswd.Details
		CharacterList []character.Character
	}
}

// Start session worker
func (s *Session) Start(table encryption.XorKeyTable) {
	// create new receiving buffer
	s.buffer = make([]byte, maxRecvBufferSize)
	s.txQueue = make(chan *Writer, txQueueMaxNum)
	s.Connected = true

	// init encryption
	s.Encryption = encryption.Encryption{}
	s.Encryption.Init(table)

	go s.reader() // start reader
	go s.writer() // start writer
}

// start socket writer
func (s *Session) writer() {
	for {
		w := <-s.txQueue

		// encrypt data
		enc, err := s.Encryption.Encrypt(w.Finalize())
		if err != nil {
			log.Error("Error encrypting packet: " + err.Error())
			continue
		}

		// send it
		len, err := s.socket.Write(enc)
		if err != nil {
			log.Error("Error sending packet: " + err.Error())
			return
		}

		// create new packet event argument
		arg := &PacketArgs{s, len, w.Type, nil}

		// trigger packet sent event
		event.Trigger(event.PacketSend, arg)
	}
}

// start socket reader
func (s *Session) reader() {
	for {
		// read data
		length, err := s.socket.Read(s.buffer)

		if err != nil {
			if err != io.EOF {
				log.Error("Error reading: " + err.Error())
			}
			s.Close()
			break
		}

		i := 0
		for i < length {
			// get packet length
			pLen := s.Encryption.GetPacketSize(s.buffer[i:])

			// check length
			if i < 0 || i > len(s.buffer) || i+pLen > len(s.buffer) {
				log.Error("Error parsing packet: slice bounds out of range!")
				s.Close()
				break
			}

			// attempt to decrypt packet
			data, err := s.Encryption.Decrypt(s.buffer[i : i+pLen])

			if err != nil {
				log.Error("Error decrypting: " + err.Error())
				s.Close()
				break
			}

			// create new packet reader
			reader := NewReader(data)

			// create new packet event argument
			arg := &PacketArgs{s, int(reader.Size), int(reader.Type), reader}

			// trigger packet received event
			event.Trigger(event.PacketReceive, arg)

			i += pLen
		}
	}
}

// Send data packet to the client
func (s *Session) Send(w *Writer) {
	// post into channel
	s.txQueue <- w
}

// Returns session's remote endpoint
func (s *Session) GetEndPnt() string {
	return s.socket.RemoteAddr().String()
}

// Returns session's ip address
func (s *Session) GetIp() string {
	var ip = strings.Split(s.GetEndPnt(), ":")
	return ip[0]
}

// Closes session socket
func (s *Session) Close() {
	s.Connected = false
	s.socket.Close()
	event.Trigger(event.ClientDisconnect, s)
}

// Info returns session's data as a string
func (s *Session) Info() string {
	return fmt.Sprintf("[conn: %d, endpnt: %s, account: %d]",
		s.UserIdx, s.GetEndPnt(), s.Data.AccountId)
}
