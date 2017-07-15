package rpc

import (
    "share/models/server"
    "share/rpc"
    srv "masterserver/server"
)

// ServerRegister RPC Call
func ServerRegister(c *rpc.Client, r server.RegRequest, s *server.RegResponse) error {
    var response = server.RegResponse{false}

    switch r.Type {
    case server.LOGIN_SERVER_TYPE:
        response.Registered = true
        g_ServerManager.NewServer(srv.Server{c, r.Type, 0, 0})
        log.Infof("Server type: LoginServer (src: %s)", c.GetEndPnt())
    case server.GAME_SERVER_TYPE:
        response.Registered = true
        g_ServerManager.NewServer(srv.Server{c, r.Type, r.ServerId, r.ChannelId})
        log.Infof("Server type: GameServer (type: %d, server: %d, channel: %d, src: %s)",
            r.ServerType, r.ServerId, r.ChannelId, c.GetEndPnt())
    default:
        log.Errorf("Unknown server type (src %s)", c.GetEndPnt())
    }

    *s = response
    return nil
}