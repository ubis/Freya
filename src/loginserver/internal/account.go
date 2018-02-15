package internal

import (
	"loginserver/net"
	"share/models/account"
	"share/models/message"
	"share/network"
	"share/rpc"
	"time"
)

// UserVerify RPC Call
func (m *Manager) UserVerify(c *rpc.Client, r *account.VerifyReq,
	s *account.VerifyRes) error {
	var verify = m.Network.VerifySession(r.UserIdx, r.AuthKey, r.DBIdx)

	if verify {
		// send server list periodically
		var t = time.NewTicker(time.Second * 5)
		go func(id uint16, n *network.Manager, p *net.Packet) {
			for {
				if !n.SendToSession(id, p.ServerState()) {
					break
				}

				<-t.C
			}
		}(r.UserIdx, m.Network, m.Packets)
	}

	*s = account.VerifyRes{Verified: verify}
	return nil
}

// OnlineCheck RPC Call
func (m *Manager) OnlineCheck(c *rpc.Client, r *account.OnlineReq,
	s *account.OnlineRes) error {
	var res = account.OnlineRes{}
	var idx = m.Network.IsOnline(r.Account)

	if idx < network.INVALID_USER_INDEX {
		// user is online in this server
		if r.Kick {
			// kick user
			if m.Network.CloseSession(idx) {
				res.Result = true
			}
		} else {
			// notify user about double login
			var msg = m.Packets.SystemMessg(message.LoginDuplicate, 0)
			if m.Network.SendToSession(idx, msg) {
				res.Result = true
			}
		}
	}

	*s = res
	return nil
}
