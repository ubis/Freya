package rpc

import (
	"github.com/ubis/Freya/cmd/gameserver/packet"
	"github.com/ubis/Freya/share/models/account"
	"github.com/ubis/Freya/share/models/message"
	"github.com/ubis/Freya/share/rpc"
)

// UserVerify RPC Call
func UserVerify(c *rpc.Client, r *account.VerifyReq, s *account.VerifyRes) error {
	var verify = g_NetworkManager.VerifyUser(r.UserIdx, r.AuthKey, r.IP, r.DBIdx)
	*s = account.VerifyRes{verify}
	return nil
}

// OnlineCheck RPC Call
func OnlineCheck(c *rpc.Client, r *account.OnlineReq, s *account.OnlineRes) error {
	res := account.OnlineRes{}
	ok, idx := g_NetworkManager.IsOnline(r.Account)

	if !ok {
		*s = res
		return nil
	}

	// user is online in this server
	if r.Kick {
		// kick user
		res.Result = g_NetworkManager.CloseUser(idx)
	} else {
		// notify user about double login
		m := packet.SystemMessg(message.LoginDuplicate, 0)
		res.Result = g_NetworkManager.SendToUser(idx, m)
	}

	*s = res
	return nil
}
