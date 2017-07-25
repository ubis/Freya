package network

import (
    "net"
    "fmt"
    "sync"
    "share/logger"
    "share/event"
    "share/models/server"
)

type Network struct {
    lock     sync.RWMutex
    clients  map[uint16]*Session
    userIdx  uint16
    settings *server.Settings
}

var log = logger.Instance()

// Network initialization
func (n *Network) Init(port int, s *server.Settings) {
    log.Info("Configuring network...")

    n.lock     = sync.RWMutex{}
    n.clients  = make(map[uint16]*Session)
    n.userIdx  = 0
    n.settings = s

    // register client disconnect event
    event.Register(event.ClientDisconnectEvent, event.Handler(n.onClientDisconnect))

    // prepare to listen for incoming connections
    // listening on Ip.Any
    var l, err = net.Listen("tcp", fmt.Sprintf(":%d", port))

    if err != nil {
        log.Fatal(err.Error())
    }

    // close the listener when the application closes
    defer l.Close()

    log.Info("Listening on " + l.Addr().String() + "...")

    for {
        // accept incoming connection
        var socket, err = l.Accept()
        if err != nil {
            log.Error("Error accepting: " + err.Error())
        }

        // create user session
        var session = Session{socket: socket}

        n.lock.RLock()
        // in case its used already...
        if n.clients[n.userIdx] != nil {
            n.lock.RUnlock()
            n.lock.Lock()
            // warning: blocked till loop is ended
            // loop till find free one
            for n.clients[n.userIdx] != nil {
                n.userIdx++
            }
            n.lock.Unlock()

            n.lock.RLock()
            // if still didn't find... ops shouldn't happen at all
            if n.clients[n.userIdx] != nil {
                n.lock.RUnlock()
                log.Error("Can't find any available user indexes!")
                session.Close()
                continue
            } else {
                n.lock.RUnlock()
            }
        } else {
            n.lock.RUnlock()
        }

        n.lock.Lock()
        n.clients[n.userIdx] = &session  // add new session
        session.UserIdx      = n.userIdx // update session user index
        n.userIdx ++
        n.lock.Unlock()

        // trigger client connect event
        event.Trigger(event.ClientConnectEvent, &session)

        // handle new client session
        go session.Start(n.settings.XorKeyTable)
    }
}

// Returns current online user count
func (n *Network) GetOnlineUsers() int {
    n.lock.RLock()
    var users = len(n.clients)
    n.lock.RUnlock()

    return users
}

// Verifies user specified by index, key and sets it's database index
func (n *Network) VerifyUser(i uint16, k uint32, ip string, db_idx int32) (bool, *Session) {
    n.lock.Lock()
    if n.clients[i] != nil && n.clients[i].AuthKey == k && n.clients[i].GetIp() == ip {
        n.clients[i].Data.Verified  = true
        n.clients[i].Data.LoggedIn  = true
        n.clients[i].Data.AccountId = db_idx
        var tmp = n.clients[i]
        n.lock.Unlock()
        return true, tmp
    }

    n.lock.Unlock()
    return false, nil
}

// onClientDisconnect event informs server about disconnected client
func (n *Network) onClientDisconnect(event event.Event) {
    var session, err = event.(*Session)
    if err != true {
        log.Error("Couldn't parse onClientDisconnect event!")
        return
    }

    n.lock.Lock()
    delete(n.clients, session.UserIdx)
    session = nil
    n.lock.Unlock()
}