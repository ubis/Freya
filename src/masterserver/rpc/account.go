package rpc

import (
    "share/rpc"
    "share/models/account"
    "golang.org/x/crypto/bcrypt"
)

// AuthCheck RPC Call
func AuthCheck(c *rpc.Client, r *account.AuthRequest, s *account.AuthResponse) error {
    var res = account.AuthResponse{Status: account.Incorrect}
    var passHash string

    var err = g_LoginDatabase.Get(&passHash,
        "SELECT password FROM accounts WHERE username = ?", r.UserId)

    if err != nil {
        *s = res
        return nil
    }

    err = bcrypt.CompareHashAndPassword([]byte(passHash), []byte(r.Password))
    if err == nil {
        g_LoginDatabase.Get(&res,
            "SELECT id, status, auth_key FROM accounts WHERE username = ?", r.UserId)

        var exist int32 = 0

        g_LoginDatabase.Get(&exist,
            "SELECT account FROM accounts_subpassword WHERE account = ?", res.Id)

        if exist == res.Id {
            res.SubPassChar = 1
        }
    }

    *s = res
    return nil
}

// UserVerify RPC Call
func UserVerify(c *rpc.Client, r *account.UserVerify, s *account.UserVerifyRecv) error {
    var t = account.UserVerifyRecv{}
    g_ServerManager.SendToGS(r.ServerId, r.ChannelId, rpc.UserVerify, r, &t)
    *s = t
    return nil
}