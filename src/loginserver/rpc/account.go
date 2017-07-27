package rpc

import (
    "time"
    "share/models/account"
    "share/rpc"
    "loginserver/packet"
)

// UserVerify RPC Call
func UserVerify(c *rpc.Client, r *account.VerifyReq, s *account.VerifyRes) error {
    var verify = g_NetworkManager.VerifyUser(r.UserIdx, r.AuthKey, r.IP, r.DBIdx)

    if verify {
        // send server list periodically
        var t = time.NewTicker(time.Second * 5)
        go func(id uint16) {
            for {
                if !g_NetworkManager.SendToUser(id, packet.ServerSate()) {
                    break
                }

                <-t.C
            }
        }(r.UserIdx)
    }

    *s = account.VerifyRes{verify}
    return nil
}