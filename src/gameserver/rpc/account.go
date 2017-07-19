package rpc

import (
    "share/models/account"
    "share/rpc"
)

// UserVerify RPC Call
func UserVerify(c *rpc.Client, r *account.UserVerify, s *account.UserVerifyRecv) error {
    *s = account.UserVerifyRecv{g_NetworkManager.VerifyUser(r.UserIdx, r.AuthKey, r.DBIdx)}
    return nil
}