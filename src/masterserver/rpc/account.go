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

        // if subpasswd exist
        var exist int32 = 0
        g_LoginDatabase.Get(&exist,
            "SELECT account FROM sub_password WHERE account = ?", res.Id)

        if exist == res.Id {
            res.SubPassChar = 1
        }

        var count = account.CharCount{}

        // fetch char list on all of servers
        for _, value := range g_DatabaseManager.DBList {
            value.DB.Get(&count.Count,
                    "SELECT COUNT(id) FROM characters WHERE id >= ? AND id <= ?",
                res.Id * 8, res.Id * 8 + 5,
            )

            if count.Count > 0 {
                count.Server = byte(value.Index)
            }

            res.CharList = append(res.CharList, count)
        }
    }

    *s = res
    return nil
}

// UserVerify RPC Call
func UserVerify(c *rpc.Client, r *account.VerifyReq, s *account.VerifyRes) error {
    var t = account.VerifyRes{}
    g_ServerManager.SendToGS(r.ServerId, r.ChannelId, rpc.UserVerify, r, &t)
    *s = t
    return nil
}

// PasswdCheck RPC Call
func PasswdCheck(c *rpc.Client, r *account.AuthCheckReq, s *account.AuthCheckRes) error {
    var res = account.AuthCheckRes{}
    var passHash string

    var err = g_LoginDatabase.Get(&passHash,
        "SELECT password FROM accounts WHERE id = ?", r.Id)

    if err != nil {
        *s = res
        return nil
    }

    err = bcrypt.CompareHashAndPassword([]byte(passHash), []byte(r.Password))
    if err == nil {
        res.Result = true
    }

    *s = res
    return nil
}