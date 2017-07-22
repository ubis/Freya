package rpc

import (
    "share/models/account"
    "share/rpc"
)

// UserVerify RPC Call
func UserVerify(c *rpc.Client, r *account.VerifyReq, s *account.VerifyResp) error {
    *s = account.VerifyResp{g_NetworkManager.VerifyUser(r.UserIdx, r.AuthKey, r.DBIdx)}
    return nil
}