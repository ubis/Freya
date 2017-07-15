package server

import (
    "sync"
    "share/logger"
    "share/rpc"
)

type Server struct {
    Client     *rpc.Client
    ServerType byte
    ServerIdx  byte
    ChannelIdx byte
}

type ServerManager struct {
    servers map[string]*Server
    lock    sync.RWMutex
}

var log = logger.Instance()

// Initializes ServerManager
func (sm *ServerManager) Init() {
    sm.servers = make(map[string]*Server)
    sm.lock    = sync.RWMutex{}
}

/*
    Adds a new server to the server list
    @param  server  a server to add
 */
func (sm *ServerManager) NewServer(server Server) {
    var endpnt = server.Client.GetEndPnt()

    sm.lock.RLock()
    var err = sm.servers[endpnt]
    sm.lock.RUnlock()

    if err != nil {
        log.Errorf("Server with endpoint exist: %s", endpnt)
        return
    }

    sm.lock.Lock()
    sm.servers[endpnt] = &server
    sm.lock.Unlock()
}

/*
    Removes the server from the server list by server's endpoint
    @param  endpnt  server's endpoint
 */
func (sm *ServerManager) RemoveServer(endpnt string) {
    sm.lock.Lock()
    delete(sm.servers, endpnt)
    sm.lock.Unlock()
}

/*
    Retrieves the server from the server list by server's endpoint
    @param  endpnt  server's endpoint
    @return server struct or nil if not found
 */
func (sm *ServerManager) GetServer(endpnt string) *Server {
    sm.lock.RLock()
    var server = sm.servers[endpnt]
    sm.lock.RUnlock()

    if server != nil {
        log.Errorf("Server with endpoint doesn't exist: %s", endpnt)
        return nil
    }

    return server
}