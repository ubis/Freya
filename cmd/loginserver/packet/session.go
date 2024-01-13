package packet

import (
	"fmt"

	"github.com/ubis/Freya/cmd/loginserver/server"
	"github.com/ubis/Freya/share/log"
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

func (session *Session) LogError(msg string) {
	log.Error("%s by %s ; account: %d",
		msg, session.GetEndPnt(), session.Account)
}

func (session *Session) LogErrorf(msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	fullMsg := fmt.Sprintf("%s by %s ; account: %d",
		formattedMsg, session.GetEndPnt(), session.Account)

	log.Error(fullMsg)
}

func (session *Session) SerializationError(opcode uint16, err error) {
	session.LogErrorf(
		"An error occurred: %s while trying to deserialize %d packet",
		err.Error(), opcode)
}
