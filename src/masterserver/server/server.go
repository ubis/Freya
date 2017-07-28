package server

import (
    "sync"
    "sort"
    "share/logger"
    "share/rpc"
    "share/models/server"
    "share/models/account"
)

type ServerManager struct {
    servers map[string]*server.Server
    lock    sync.RWMutex
}

var log = logger.Instance()

// Initializes ServerManager
func (sm *ServerManager) Init() {
    sm.servers = make(map[string]*server.Server)
    sm.lock    = sync.RWMutex{}
}

// Adds a new server to the server list
func (sm *ServerManager) NewServer(server server.Server) {
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

// Removes the server from the server list by server's endpoint
func (sm *ServerManager) RemoveServer(endpnt string) {
    sm.lock.Lock()
    delete(sm.servers, endpnt)
    sm.lock.Unlock()
}

// Retrieves the server from the server list by server's endpoint
func (sm *ServerManager) GetServer(endpnt string) *server.Server {
    sm.lock.RLock()
    var server = sm.servers[endpnt]
    sm.lock.RUnlock()

    if server != nil {
        log.Errorf("Server with endpoint doesn't exist: %s", endpnt)
        return nil
    }

    return server
}

// Returns sorted game server list
func (sm *ServerManager) GetGameServerList() []server.ServerItem {
    var serverList []server.ServerItem

    sm.lock.RLock()
    for _, value := range sm.servers {
        // this is login server
        if value.ServerId == 0 {
            continue
        }

        var chn = server.ChannelItem{
            value.ChannelId,
            value.Type,
            value.PublicIp,
            value.PublicPort,
            value.CurrentUsers,
            value.MaxUsers,
        }

        var svr   = &server.ServerItem{Id: value.ServerId}
        var found = false

        for i := 0; i < len(serverList); i ++ {
            if serverList[i].Id == value.ServerId {
                svr   = &serverList[i]
                found = true
                break
            }
        }

        svr.List = append(svr.List, chn)

        if !found {
            serverList = append(serverList, *svr)
        }
    }
    sm.lock.RUnlock()

    // sort servers by id
    sort.Sort(server.ByServer(serverList))

    // sort each server channels by id
    for i := 0; i < len(serverList); i ++ {
        sort.Sort(server.ByChannel(serverList[i].List))
    }

    return serverList
}

// Attempts to find GameServer with server and channel id's and send RPC call
func (sm *ServerManager) SendToGS(s byte, c byte, m string, a interface{}, r interface{}) error {
    sm.lock.RLock()
    for _, value := range sm.servers {
        if value.ServerId == s && value.ChannelId == c {
            var c = value.Client
            sm.lock.RUnlock()
            return c.Call(m, a, r)
        }
    }
    sm.lock.RUnlock()
    
    return nil
}

// Attempts to find LoginServer and send RPC call
func (sm *ServerManager) SendToLS(m string, a interface{}, r interface{}) error  {
    sm.lock.RLock()
    for _, value := range sm.servers {
        if value.ServerId == 0 {
            var c = value.Client
            sm.lock.RUnlock()
            return c.Call(m, a, r)
        }
    }
    sm.lock.RUnlock()

    return nil
}

// Attempts to find logged in user in all servers by account id
func (sm *ServerManager) IsOnline(id int32) bool {
    sm.lock.RLock()
    for _, value := range sm.servers {
        var c   = value.Client
        var res = account.OnlineRes{}
        var err = c.Call(rpc.OnlineCheck, account.OnlineReq{id, false}, &res)

        if err == nil && res.Result {
            sm.lock.RUnlock()
            return true
        }
    }
    sm.lock.RUnlock()

    return false
}

// Attempts to kick out user in all servers by account id
func (sm *ServerManager) KickAccount(id int32) bool {
    sm.lock.RLock()
    for _, value := range sm.servers {
        var c   = value.Client
        var res = account.OnlineRes{}
        var err = c.Call(rpc.OnlineCheck, account.OnlineReq{id, true}, &res)

        if err == nil && res.Result {
            sm.lock.RUnlock()
            return true
        }
    }
    sm.lock.RUnlock()

    return false
}