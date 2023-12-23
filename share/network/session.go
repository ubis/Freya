package network

import (
	"io"
	"net"
	"strings"
	"sync"

	"github.com/ubis/Freya/share/encryption"
	"github.com/ubis/Freya/share/event"
	"github.com/ubis/Freya/share/log"
)

// max buffer size
const MAX_RECV_BUFFER_SIZE = 4096

type Session struct {
	socket net.Conn
	buffer []byte

	Encryption encryption.Encryption
	UserIdx    uint16
	AuthKey    uint32
	Ses        any

	PeriodicJobs map[string]*PeriodicTask
	jobMutex     sync.Mutex
}

func (s *Session) Store(ses any) {
	s.Ses = ses
}

func (s *Session) Retrieve() any {
	return s.Ses
}

func (s *Session) GetUserIdx() uint16 {
	return s.UserIdx
}

func (s *Session) GetAuthKey() uint32 {
	return s.AuthKey
}

func (s *Session) GetSeed() uint32 {
	return s.Encryption.Key.Seed2nd
}

func (s *Session) GetKeyIdx() uint32 {
	return s.Encryption.RecvXorKeyIdx
}

// Starts session goroutine
func (s *Session) Start(table *encryption.XorKeyTable) {
	// create new receiving buffer
	s.buffer = make([]byte, MAX_RECV_BUFFER_SIZE)
	// create map to store periodic tasks
	s.PeriodicJobs = make(map[string]*PeriodicTask)

	// init encryption
	s.Encryption = encryption.Encryption{}
	s.Encryption.Init(table)

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
			arg := &PacketArgs{
				Session: s,
				Length:  int(reader.Size),
				Type:    int(reader.Type),
				Data:    data,
				Reader:  reader,
			}

			// trigger packet received event
			event.Trigger(event.PacketReceiveEvent, arg)

			i += packetLength
		}
	}
}

// Sends specified data to the client
func (s *Session) Send(writer *Writer) {
	data := writer.Finalize()

	// encrypt data
	encrypt, err := s.Encryption.Encrypt(data)
	if err != nil {
		log.Error("Error encrypting packet: " + err.Error())
		return
	}

	// send it...
	length, err := s.socket.Write(encrypt)
	if err != nil {
		log.Error("Error sending packet: " + err.Error())
		return
	}

	// create new packet event argument
	arg := &PacketArgs{
		Session: s,
		Length:  length,
		Type:    writer.Type,
		Data:    data,
		Reader:  nil,
	}

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

// GetLocalEndPntIp returns local end point IP address.
// Local end point is server to which session is connected to.
func (s *Session) GetLocalEndPntIp() string {
	pnt := s.socket.LocalAddr().String()
	ip := strings.Split(pnt, ":")
	return ip[0]
}

// IsLocal checks if session's remote endpoint originated from private network.
func (s *Session) IsLocal() bool {
	return net.IP.IsPrivate(net.ParseIP(s.GetIp()))
}

// Closes session socket
func (s *Session) Close() {
	s.RemoveAllJobs()
	s.socket.Close()
	event.Trigger(event.ClientDisconnectEvent, s)
}

// AddJob adds a periodic task to the session's job list in a thread-safe manner.
// The job can be referenced and managed by its name.
func (s *Session) AddJob(name string, task *PeriodicTask) {
	s.jobMutex.Lock()
	defer s.jobMutex.Unlock()

	s.PeriodicJobs[name] = task
}

// RemoveJob stops and removes a periodic task from the session's job list
// based on its name in a thread-safe manner.
func (s *Session) RemoveJob(name string) {
	s.jobMutex.Lock()
	defer s.jobMutex.Unlock()

	if job, exists := s.PeriodicJobs[name]; exists {
		job.Stop()
		delete(s.PeriodicJobs, name)
	}
}

// RemoveAllJobs stops all periodic tasks and clears the job list
// for the session in a thread-safe manner.
func (s *Session) RemoveAllJobs() {
	s.jobMutex.Lock()
	defer s.jobMutex.Unlock()

	for name, job := range s.PeriodicJobs {
		job.Stop()
		delete(s.PeriodicJobs, name)
	}
}
