package packet

import (
	"github.com/ubis/Freya/cmd/loginserver/server"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

type SessionState int

const (
	StateUnknown SessionState = iota
	StateConnected
	StateVerified
)

type Session struct {
	network.SessionHandler

	state SessionState

	ServerConfig   *server.Config
	ServerInstance *server.Instance
	RPC            *rpc.Client
	Account        int32
}

// Create a new server-specified client session
func NewSession(s *network.Session, inst *server.Instance) *Session {
	return &Session{
		SessionHandler: s,
		ServerConfig:   inst.Config,
		ServerInstance: inst,
		RPC:            inst.RPC,
	}
}
