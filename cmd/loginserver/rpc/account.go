package rpc

import (
	"time"

	"github.com/ubis/Freya/cmd/loginserver/packet"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/account"
	"github.com/ubis/Freya/share/models/message"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

// UserVerify RPC Call
func (e *RPC) UserVerify(c *rpc.Client, r *account.VerifyReq, s *account.VerifyRes) error {
	net := e.Server

	state := net.VerifyUser(r.UserIdx, r.AuthKey, r.IP)
	*s = account.VerifyRes{Verified: state}

	if !state {
		return nil
	}

	session := net.GetSession(r.UserIdx)

	// session is verified
	// create new periodic task to send server list periodically
	task := network.NewPeriodicTask(time.Second*5, func() {
		ses, ok := session.Retrieve().(*packet.Session)
		if !ok {
			log.Error("Unable to parse client session!")
			return
		}

		// restore account data
		ses.Account = r.DBIdx

		ses.Send(packet.ServerSate(ses))
	})

	session.AddJob("ServerState", task)
	return nil
}

// OnlineCheck RPC Call
func (e *RPC) OnlineCheck(c *rpc.Client, r *account.OnlineReq, s *account.OnlineRes) error {
	users := e.Server.GetUsers()
	online := false
	var index uint16

	for _, v := range users {
		ses, ok := v.Retrieve().(*packet.Session)
		if !ok {
			log.Error("Unable to parse client session!")
			continue
		}

		if ses.Account == r.Account {
			online = true
			index = v.GetUserIdx()
			break
		}
	}

	res := account.OnlineRes{}

	if !online {
		*s = res
		return nil
	}

	// user is online in this server
	if r.Kick {
		// kick user
		res.Result = e.Server.CloseUser(index)
	} else {
		// notify user about double login
		m := packet.SystemMessage(message.LoginDuplicate, 0)
		res.Result = e.Server.SendToUser(index, m)
	}

	*s = res
	return nil
}
