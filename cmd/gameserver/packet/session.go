package packet

import (
	"fmt"
	"sync"

	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/cmd/gameserver/server"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/character"
	"github.com/ubis/Freya/share/models/subpasswd"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

type SessionState int

const (
	StateUnknown SessionState = iota
	StateConnected
	StateVerified
	StateInGame
)

type Session struct {
	network.SessionHandler

	state SessionState

	ServerConfig   *server.Config
	ServerInstance *server.Instance
	RPC            *rpc.Client
	Account        int32

	SubPassword      *subpasswd.Details
	Characters       []*character.Character
	PasswordVerified bool

	Character *character.Character
	mutex     sync.RWMutex

	Cell         context.CellHandler
	World        context.WorldHandler
	WorldManager context.WorldManagerHandler
}

// Create a new server-specified client session
func NewSession(s *network.Session, inst *server.Instance, wm context.WorldManagerHandler) *Session {
	return &Session{
		SessionHandler: s,
		ServerConfig:   inst.Config,
		ServerInstance: inst,
		RPC:            inst.RPC,
		WorldManager:   wm,
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

func (session *Session) LogFatal(msg string) {
	log.Error("%s by %s ; account: %d",
		msg, session.GetEndPnt(), session.Account)
	session.Close()
}

func (session *Session) LogFatalf(msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	fullMsg := fmt.Sprintf("%s by %s ; account: %d",
		formattedMsg, session.GetEndPnt(), session.Account)

	log.Error(fullMsg)
	session.Close()
}

func (session *Session) SetState(state SessionState) {
	session.mutex.Lock()
	defer session.mutex.Unlock()

	session.state = state
}

func (session *Session) GetState() SessionState {
	session.mutex.RLock()
	defer session.mutex.RUnlock()

	return session.state
}

func (session *Session) FindPlayerByIndex(index uint16) *Session {
	player := session.ServerInstance.Server.GetSession(index)
	if player == nil {
		return nil
	}

	ses, ok := session.Retrieve().(*Session)
	if !ok {
		log.Error("Unable to parse client session!")
		return nil
	}

	if ses.GetState() != StateInGame {
		return nil
	}

	return ses
}

func SetCurrentWorld(session network.SessionHandler, world context.WorldHandler, cell context.CellHandler) {
	ses, ok := session.Retrieve().(*Session)
	if !ok {
		log.Error("Unable to parse client session!")
		return
	}

	ses.mutex.Lock()
	defer ses.mutex.Unlock()

	ses.Cell = cell
	ses.World = world
}

func GetCurrentWorld(session network.SessionHandler) context.WorldHandler {
	ses, ok := session.Retrieve().(*Session)
	if !ok {
		log.Error("Unable to parse client session!")
		return nil
	}

	ses.mutex.RLock()
	defer ses.mutex.RUnlock()

	return ses.World
}

func GetCurrentCell(session network.SessionHandler) context.CellHandler {
	ses, ok := session.Retrieve().(*Session)
	if !ok {
		log.Error("Unable to parse client session!")
		return nil
	}

	return ses.Cell
}

func GetCurrentCellByPos(session network.SessionHandler) context.CellHandler {
	ses, ok := session.Retrieve().(*Session)
	if !ok {
		log.Error("Unable to parse client session!")
		return nil
	}

	x, y := ses.Character.GetPosition()

	id := ses.Character.GetWorld()
	world := ses.WorldManager.FindWorld(id)
	if world == nil {
		ses.LogFatalf("Unable to find world: %d for character %d",
			id, ses.Character.Id)
		return nil
	}

	return world.FindCell(int(x), int(y))
}
