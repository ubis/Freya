package rpc

import (
    "share/rpc"
    "share/models/account"
)

// FetchSubPassword RPC Call
func FetchSubPassword(c *rpc.Client, r *account.SubPasswordReq, s *account.SubPassword) error {
    var res = account.SubPassword{}

    g_LoginDatabase.Get(&res,
        "SELECT password, answer, question, expires " +
            "FROM accounts_subpassword " +
            "WHERE account = ?", r.Account)

    *s = res
    return nil
}

// SetSubPassword RPC Call
func SetSubPassword(c *rpc.Client, r *account.SetSubPass, s *account.SubPassResp) error {
    var res = account.SubPassResp{true}
    var exist = 0

    g_LoginDatabase.Get(&exist,
        "SELECT account FROM accounts_subpassword WHERE account = ?", r.Account)

    if exist == 1 {
        // changing subpassword
        g_LoginDatabase.MustExec(
            "UPDATE accounts_subpassword " +
                "SET password = ?, answer = ?, question = ?, expires = ? WHERE account = ?",
            r.Password, r.Answer, r.Question, r.Expires, r.Account)
    } else {
        // creating subpassword
        g_LoginDatabase.MustExec(
            "INSERT INTO accounts_subpassword " +
                "(account, password, answer, question, expires)" +
                "VALUES (?, ?, ?, ?, ?)",
            r.Account, r.Password, r.Answer, r.Question, r.Expires)
    }

    *s = res
    return nil
}

// RemoveSubPassword RPC Call
func RemoveSubPassword(c *rpc.Client, r *account.SubPasswordReq, s *account.SubPassResp) error {
    var res = account.SubPassResp{true}

    g_LoginDatabase.MustExec(
        "DELETE FROM accounts_subpassword WHERE account = ?", r.Account)

    *s = res
    return nil
}