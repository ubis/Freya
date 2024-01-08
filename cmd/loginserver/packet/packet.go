package packet

import (
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/script"
)

type PacketFunc func(*Session, *network.Reader)

func register(p *network.PacketHandler, opcode uint16, method PacketFunc) {
	name, ok := opcodeNames[opcode]
	if !ok {
		name = "Unknown"
	}

	if method == nil {
		p.Register(opcode, name, nil)
		return
	}

	p.Register(opcode, name, func(s *network.Session, r *network.Reader) {
		session, ok := s.Retrieve().(*Session)
		if !ok {
			log.Error("Unable to parse client session!")
			return
		}

		method(session, r)
	})
}

func verifyState(session *Session, state SessionState) bool {
	if session.state != state {
		log.Errorf("Invalid client state [need: %d ; have: %d] src: %s",
			state, session.state, session.GetEndPnt())
		session.Close()
		return false
	}

	return true
}

// Registers network packets
func RegisterPackets(h *network.PacketHandler) {
	log.Info("Registering packets...")

	register(h, CSCConnect2Svr, Connect2Svr)
	register(h, CSCVerifyLinks, VerifyLinks)
	register(h, CSCAuthAccount, AuthAccount)
	register(h, CSCForceDisconnect, ForceDisconnect)
	register(h, NFYSystemMessage, nil)
	register(h, NFYServerState, nil)
	register(h, CSCCheckVersion, CheckVersion)
	register(h, NFYUrlToClient, nil)
	register(h, CSCPublicKey, PublicKey)
	register(h, CSCPreServerEnvRequest, PreServerEnvRequest)
	register(h, NFYDisconnectTimer, nil)
	register(h, CSCAuthenticate, Authenticate)
	register(h, NFYAuthTimer, nil)
	register(h, CSCUnknown3383, nil)
	register(h, CSCUnknown5383, nil)
}

func RegisterFunc() {
	script.RegisterFunc("sendClientPacket", sessionPacketFunc{})
	script.RegisterFunc("sendClientMessage", clientMessageFunc{Fn: SystemMessageEx})
}
