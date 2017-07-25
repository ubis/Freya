package rpc

import (
    "time"
    "share/models/account"
    "share/rpc"
    "share/network"
    "loginserver/packet"
)

// UserVerify RPC Call
func UserVerify(c *rpc.Client, r *account.VerifyReq, s *account.VerifyRes) error {
    var verify, session = g_NetworkManager.VerifyUser(r.UserIdx, r.AuthKey, r.IP, r.DBIdx)

    if verify {
        // send server list periodically
        var t = time.NewTicker(time.Second * 5)
        go func(s *network.Session) {
            for {
                if !s.Connected {
                    break
                }

                packet.ServerSate(s)
                <-t.C
            }
        }(session)
    }

    *s = account.VerifyRes{verify}
    return nil
}