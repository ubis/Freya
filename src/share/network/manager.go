package network

import (
	"fmt"
	"net"
	"share/encryption"
	"share/event"
	"share/log"
	"sync"
)

// Manager structure
type Manager struct {
	port     int
	lock     sync.RWMutex
	clients  map[uint16]*Session
	userIdx  uint16
	packets  map[uint16]*PacketData
	keytable *encryption.XorKeyTable
}

// Init network manager
func (m *Manager) Init(port int) {
	log.Info("Configuring network...")

	m.port = port
	m.lock = sync.RWMutex{}
	m.clients = make(map[uint16]*Session)
	m.userIdx = 0
	m.packets = make(map[uint16]*PacketData)
	m.keytable = &encryption.XorKeyTable{}

	// init encryption table
	m.keytable.Init()

	// register client disconnect event
	event.Register(event.ClientDisconnect, event.Handler(m.onClientDisconnect))
}

// Run network server
func (m *Manager) Run() {
	// prepare to listen for incoming connections
	// listening on Ip.Any
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", m.port))

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

		// create user session
		session := &Session{socket: socket}

		m.lock.RLock()
		// in case its used already...
		if m.clients[m.userIdx] != nil {
			m.lock.RUnlock()
			m.lock.Lock()
			// warning: blocked till loop is ended
			// loop till find free one
			for m.clients[m.userIdx] != nil {
				m.userIdx++
			}
			m.lock.Unlock()

			m.lock.RLock()
			// if still didn't found... ops this shouldn't happen at all
			if m.clients[m.userIdx] != nil {
				m.lock.RUnlock()
				log.Error("Couldn't find any available user indexes!")
				session.Close()
				continue
			} else {
				m.lock.RUnlock()
			}
		} else {
			m.lock.RUnlock()
		}

		m.lock.Lock()
		m.clients[m.userIdx] = session // add new session
		session.UserIdx = m.userIdx    // update session user index
		m.userIdx++
		m.lock.Unlock()

		// trigger client connect event
		event.Trigger(event.ClientConnect, &session)

		// handle new client session
		go session.Start(m.keytable)
	}
}

// GetSessionCount returns online session count
func (m *Manager) GetSessionCount() int {
	m.lock.RLock()
	users := len(m.clients)
	m.lock.RUnlock()

	return users
}

// VerifySession verifies session by index and key, restores account index
func (m *Manager) VerifySession(idx uint16, key uint32, account int32) bool {
	m.lock.RLock()
	session := m.clients[idx]
	m.lock.RUnlock()

	if session != nil && session.AuthKey == key {
		session.Data.Verified = true
		session.Data.LoggedIn = true
		session.Data.AccountId = account
		return true
	}

	return false
}

// SendToSession sends packet to session specified by index
func (m *Manager) SendToSession(idx uint16, writer *Writer) bool {
	m.lock.RLock()
	session := m.clients[idx]
	m.lock.RUnlock()

	if session != nil && session.Connected {
		session.Send(writer)
		return true
	}

	return false
}

// IsOnline checks if account is online and returns user index
func (m *Manager) IsOnline(account int32) uint16 {
	// copy session list
	/*list := make(map[uint16]*Session)

	m.lock.RLock()
	copy(list, m.clients)
	m.lock.RUnlock()

	for _, s := range list {
		if s.Data.AccountId == account && s.Data.Verified && s.Data.LoggedIn {
			return s.UserIdx
		}
	}*/

	return INVALID_USER_INDEX
}

// CloseSession connection by index
func (m *Manager) CloseSession(idx uint16) bool {
	m.lock.RLock()
	if m.clients[idx] != nil {
		m.clients[idx].Close()
		m.lock.RUnlock()
		return true
	}

	m.lock.RUnlock()
	return false
}

// onClientDisconnect event informs server about disconnected client
func (m *Manager) onClientDisconnect(event event.Event) {
	session, err := event.(*Session)
	if err != true {
		log.Error("Couldn't parse onClientDisconnect event!")
		return
	}

	m.lock.Lock()
	delete(m.clients, session.UserIdx)
	m.lock.Unlock()
}
