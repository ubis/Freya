package rpc

import (
	"github.com/ubis/Freya/share/models/subpasswd"
	"github.com/ubis/Freya/share/rpc"
)

// FetchSubPassword RPC Call
func FetchSubPassword(c *rpc.Client, r *subpasswd.FetchReq, s *subpasswd.FetchRes) error {
	var res = subpasswd.Details{}

	g_LoginDatabase.Get(&res,
		"SELECT password, answer, question, expires "+
			"FROM sub_password WHERE account = ?", r.Account)

	*s = subpasswd.FetchRes{res}
	return nil
}

// SetSubPassword RPC Call
func SetSubPassword(c *rpc.Client, r *subpasswd.SetReq, s *subpasswd.SetRes) error {
	var res = subpasswd.SetRes{true}
	var exist = 0

	g_LoginDatabase.Get(&exist,
		"SELECT account FROM sub_password WHERE account = ?", r.Account)

	if exist > 0 {
		// changing subpassword
		g_LoginDatabase.MustExec(
			"UPDATE sub_password "+
				"SET password = ?, answer = ?, question = ?, expires = ? WHERE account = ?",
			r.Password, r.Answer, r.Question, r.Expires, r.Account)
	} else {
		// creating subpassword
		g_LoginDatabase.MustExec(
			"INSERT INTO sub_password "+
				"(account, password, answer, question, expires)"+
				"VALUES (?, ?, ?, ?, ?)",
			r.Account, r.Password, r.Answer, r.Question, r.Expires)
	}

	*s = res
	return nil
}

// RemoveSubPassword RPC Call
func RemoveSubPassword(c *rpc.Client, r *subpasswd.SetReq, s *subpasswd.SetRes) error {
	var res = subpasswd.SetRes{true}

	g_LoginDatabase.MustExec("DELETE FROM sub_password WHERE account = ?", r.Account)

	*s = res
	return nil
}
