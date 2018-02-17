package network

import (
	"fmt"
	"net"
	"share/encryption"
	"share/event"
	"share/log"
	"sync"
)

// Server listener structure
type Server struct {
	port    int
	lock    sync.RWMutex
	connLst map[int]*Session
	connIdx uint16
	keys    encryption.XorKeyTable
}

// Init network server server
func (s *Server) Init(port int) {
	log.Info("Configuring network server...")

	s.port = port
	s.lock = sync.RWMutex{}
	s.connLst = make(map[int]*Session)
	s.connIdx = 0
	s.keys = encryption.XorKeyTable{}

	// init encryption table
	s.keys.Init()

	// register client disconnect event
	event.Register(event.ClientDisconnect, event.Handler(s.onClientDisconnect))
}

// Run the server listener
func (s *Server) Run() {
	// prepare to listen for incoming connections
	// listening on Ip.Any
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))

	if err != nil {
		log.Fatal(err.Error())
	}

	// close the listener when the application closes
	defer l.Close()

	log.Info("Listening on " + l.Addr().String() + "...")

	for {
		// accept incoming connection
		socket, err := l.Accept()
		if err != nil {
			log.Error("Error accepting: " + err.Error())
			continue
		}

		// in case index is used already...
		s.lock.Lock()
		if s.connLst[int(s.connIdx)] != nil {
			s.lock.Unlock()

			// change it
			s.connIdx++

			// if it's still in-use then drop new connection
			s.lock.Lock()
			if s.connLst[int(s.connIdx)] != nil {
				s.lock.Unlock()
				log.Error("Couldn't find any available user index!")
				socket.Close()
				continue
			} else {
				s.lock.Unlock()
			}
		} else {
			s.lock.Unlock()
		}

		// create user session
		session := &Session{socket: socket, UserIdx: s.connIdx}

		// increase connection index
		s.connIdx++

		// trigger client connect event
		event.Trigger(event.ClientConnect, session)

		// handle new client session
		go session.Start(s.keys)

		s.lock.Lock()
		s.connLst[int(s.connIdx)] = session // add new session
		s.lock.Unlock()
	}
}

// GetSessionCount returns online session count
func (s *Server) GetSessionCount() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.connLst)
}

// VerifySession verifies session by conn index and a key
// and restores account index, returns true on success
func (s *Server) VerifySession(idx uint16, key uint32, account int32) bool {
	s.lock.Lock()
	session := s.connLst[int(idx)]
	s.lock.Unlock()

	if session != nil && session.AuthKey == key {
		session.Data.Verified = true
		session.Data.LoggedIn = true
		session.Data.AccountId = account
		return true
	}

	return false
}

// SendToSession sends packet to session specified by index
// returns true on success
func (s *Server) SendToSession(idx uint16, writer *Writer) bool {
	s.lock.Lock()
	session := s.connLst[int(idx)]
	s.lock.Unlock()

	if session != nil && session.Connected {
		session.Send(writer)
		return true
	}

	return false
}

// IsSessionOnline checks if session is online by account index
// and returns user index on success
func (s *Server) IsSessionOnline(account int32) uint16 {
	/*list := make(map[int]*Session)

	// copy session list
	s.lock.Lock()
	copy(list, s.connLst)
	s.lock.Unlock()

	for _, u := range list {
		d := u.Data
		if d.AccountId == account && d.Verified && d.LoggedIn {
			return u.UserIdx
		}
	}*/

	return INVALID_USER_INDEX
}

// CloseSession closes session's connection by it's index
func (s *Server) CloseSession(idx uint16) bool {
	s.lock.Lock()
	session := s.connLst[int(idx)]
	s.lock.Unlock()

	if session != nil {
		session.Close()
		return true
	}

	return false
}

// onClientDisconnect event informs server about disconnected client
func (s *Server) onClientDisconnect(event event.Event) {
	session, ok := event.(*Session)
	if ok != true {
		return
	}

	s.lock.Lock()
	delete(s.connLst, int(session.UserIdx))
	s.lock.Unlock()
}
