package packet

import (
	"runtime"
	"runtime/debug"

	"github.com/ubis/Freya/cmd/loginserver/def"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/network"
)

var g_ServerConfig = def.ServerConfig
var g_ServerSettings = def.ServerSettings
var g_PacketHandler = def.PacketHandler
var g_RPCHandler = def.RPCHandler

func NotifyServerInfo(session *network.Session) {
	msg := "Welcome to Freya - CABAL Server Emulator!"
	session.Send(SystemMessgEx(msg))

	msg = "Running on " + runtime.GOOS + " OS with " + runtime.Version()
	session.Send(SystemMessgEx(msg))

	var buildCommit = func() string {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return ""
		}

		for _, setting := range info.Settings {
			if setting.Key != "vcs.revision" {
				continue
			}
			return setting.Value

		}

		return ""
	}()

	msg = "Build: #" + buildCommit[:6]
	session.Send(SystemMessgEx(msg))
}

// Registers network packets
func RegisterPackets() {
	log.Info("Registering packets...")

	var pk = g_PacketHandler
	pk.Register(CONNECT2SVR, "Connect2Svr", Connect2Svr)
	pk.Register(VERIFYLINKS, "VerifyLinks", VerifyLinks)
	pk.Register(AUTHACCOUNT, "AuthAccount", AuthAccount)
	pk.Register(FDISCONNECT, "FDisconnect", FDisconnect)
	pk.Register(SYSTEMMESSG, "SystemMessg", nil)
	pk.Register(SERVERSTATE, "ServerState", nil)
	pk.Register(CHECKVERSION, "CheckVersion", CheckVersion)
	pk.Register(URLTOCLIENT, "URLToClient", nil)
	pk.Register(PUBLIC_KEY, "PublicKey", PublicKey)
	pk.Register(PRE_SERVER_ENV_REQUEST, "PreServerEnvRequest", PreServerEnvRequest)
}
