package rpc

import (
	"github.com/ubis/Freya/cmd/gameserver/packet"
	"github.com/ubis/Freya/share/models/account"
	"github.com/ubis/Freya/share/models/message"
	"github.com/ubis/Freya/share/network"
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
	var res = account.OnlineRes{}
	var idx = g_NetworkManager.IsOnline(r.Account)

	if idx < network.INVALID_USER_INDEX {
		// user is online in this server
		if r.Kick {
			// kick user
			if g_NetworkManager.CloseUser(idx) {
				res.Result = true
			}
		} else {
			// notify user about double login
			var m = packet.SystemMessg(message.LoginDuplicate, 0)
			if g_NetworkManager.SendToUser(idx, m) {
				res.Result = true
			}
		}
	}

	*s = res
	return nil
}
