package rpc

import (
    "share/models/server"
    "share/rpc"
)

func ServerRegister(c *rpc.Client, r server.RegRequest, s *server.RegResponse) error {
    var response = server.RegResponse{false}

    switch r.Type {
    case server.LOGIN_SERVER_TYPE:
        response.Registered = true
        log.Infof("Server type: LoginServer (src: %s)", c.GetEndPnt())
    }

    *s = response
    return nil
}