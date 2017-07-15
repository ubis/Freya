package network

import (
    "net"
    "fmt"
    "sync"
    "share/logger"
    "share/event"
    "share/models/server"
)

var log     = logger.Instance()
var lock    = sync.RWMutex{}
var clients = make(map[uint16]*Session)

var userIdx  uint16 = 0
var settings *server.Settings
/*
    Network initialization
    @param  port    Server port to listen on
 */
func Init(port int, s *server.Settings) {
    settings = s
    log.Info("Configuring network...")

    // register client disconnect event
    event.Register(event.ClientDisconnectEvent, event.Handler(OnClientDisconnect))

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

        lock.RLock()
        // in case its used already...
        if clients[userIdx] != nil {
            // loop till find free one
            for clients[userIdx] != nil {
                userIdx++
            }
        } else {
            clients[userIdx] = &session
        }

        // set client session's user index
        session.UserIdx       = userIdx
        settings.CurrentUsers = len(clients)
        userIdx++
        lock.RUnlock()

        // trigger client connect event
        event.Trigger(event.ClientConnectEvent, &session)

        // handle new client session
        go session.Start(settings.XorKeyTable)
    }
}

/*
    OnClientDisconnect event, informs server about disconnected client
    @param  event   Event interface, which is later parsed into Session struct
 */
func OnClientDisconnect(event event.Event) {
    var session, err = event.(*Session)
    if err != true {
        log.Error("Couldn't parse onClientDisconnect event!")
        return
    }

    lock.Lock()
    delete(clients, session.UserIdx)
    session = nil
    settings.CurrentUsers = len(clients)
    lock.Unlock()
}