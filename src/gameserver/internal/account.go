package internal

import (
	"share/models/account"
	"share/models/message"
	"share/network"
	"share/rpc"
)

// UserVerify RPC Call
func (cm *Comm) UserVerify(c *rpc.Client, r *account.VerifyReq,
	s *account.VerifyRes) error {
	verify := cm.Net.VerifySession(r.UserIdx, r.AuthKey, r.DBIdx)
	*s = account.VerifyRes{Verified: verify}
	return nil
}

// OnlineCheck RPC Call
func (cm *Comm) OnlineCheck(c *rpc.Client, r *account.OnlineReq,
	s *account.OnlineRes) error {
	res := account.OnlineRes{}
	idx := cm.Net.IsSessionOnline(r.Account)

	if idx < network.INVALID_USER_INDEX {
		// user is online in this server
		if r.Kick {
			// kick user
			if cm.Net.CloseSession(idx) {
				res.Result = true
			}
		} else {
			// notify user about double login
			msg := cm.Lst.SystemMessg(message.LoginDuplicate, 0)
			if cm.Net.SendToSession(idx, msg) {
				res.Result = true
			}
		}
	}

	*s = res
	return nil
}
