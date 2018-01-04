package network

import (
	"io"
	"net"
	"share/encryption"
	"share/event"
	"share/log"
	"share/models/character"
	"share/models/subpasswd"
	"strings"
)

// max buffer size
const MAX_RECV_BUFFER_SIZE = 4096

type Session struct {
	socket net.Conn
	buffer []byte

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

// Starts session goroutine
func (s *Session) Start(table encryption.XorKeyTable) {
	// create new receiving buffer
	s.buffer = make([]byte, MAX_RECV_BUFFER_SIZE)

	s.Connected = true

	// init encryption
	s.Encryption = encryption.Encryption{}
	s.Encryption.Init(&table)

	for {
		// read data
		var length, err = s.socket.Read(s.buffer)

		if err != nil {
			if err != io.EOF {
				log.Error("Error reading: " + err.Error())
			}
			s.Close()
			break
		}

		var i = 0
		for i < length {
			// get packet length
			var packetLength = s.Encryption.GetPacketSize(s.buffer[i:])

			// check length
			if i < 0 || i > len(s.buffer) || i+packetLength > len(s.buffer) {
				log.Error("Error parsing packet: slice bounds out of range!")
				s.Close()
				break
			}

			// attempt to decrypt packet
			var data, error = s.Encryption.Decrypt(s.buffer[i : i+packetLength])

			if error != nil {
				log.Error("Error decrypting: " + error.Error())
				s.Close()
				break
			}

			// create new packet reader
			var reader = NewReader(data)

			// create new packet event argument
			var arg = &PacketArgs{s, int(reader.Size), int(reader.Type), reader}

			// trigger packet received event
			event.Trigger(event.PacketReceiveEvent, arg)

			i += packetLength
		}
	}
}

// Sends specified data to the client
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
	var arg = &PacketArgs{s, length, writer.Type, nil}

	// trigger packet sent event
	event.Trigger(event.PacketSendEvent, arg)
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
	event.Trigger(event.ClientDisconnectEvent, s)
}
