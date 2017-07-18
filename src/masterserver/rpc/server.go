package rpc

import (
    "share/models/server"
    "share/rpc"
)

// ServerRegister RPC Call
func ServerRegister(c *rpc.Client, r *server.ServerData, s *server.RegResponse) error {
    var response = server.RegResponse{false}

    switch r.Type {
    case server.LOGIN_SERVER_TYPE:
        response.Registered = true
        g_ServerManager.NewServer(server.Server{r, c})
        log.Infof("Server type: LoginServer (src: %s)", c.GetEndPnt())
    case server.GAME_SERVER_TYPE:
        response.Registered = true
        g_ServerManager.NewServer(server.Server{r, c})
        log.Infof("Server type: GameServer (type: %d, server: %d, channel: %d, src: %s)",
            r.ServerType, r.ServerId, r.ChannelId, c.GetEndPnt())
    default:
        log.Errorf("Unknown server type (src %s)", c.GetEndPnt())
    }

    *s = response
    return nil
}

// ServerList RPC Call
func ServerList(c *rpc.Client, r *server.SvrListRequest, s *server.SvrListResponse) error {
    var list = g_ServerManager.GetGameServerList()
    *s = server.SvrListResponse{list}
    return nil
}