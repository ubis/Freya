package rpc

import (
    "share/models/account"
    "share/rpc"
)

// UserVerify RPC Call
func UserVerify(c *rpc.Client, r *account.VerifyReq, s *account.VerifyRes) error {
    *s = account.VerifyRes{g_NetworkManager.VerifyUser(r.UserIdx, r.AuthKey, r.DBIdx)}
    return nil
}