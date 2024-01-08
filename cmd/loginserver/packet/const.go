package packet

// Packet opcode values
const (
	CSCConnect2Svr         = 101
	CSCVerifyLinks         = 102
	CSCAuthAccount         = 103
	CSCForceDisconnect     = 110
	NFYSystemMessage       = 120
	NFYServerState         = 121
	CSCCheckVersion        = 122
	NFYUrlToClient         = 128
	CSCPublicKey           = 2001
	CSCPreServerEnvRequest = 2002
	NFYDisconnectTimer     = 2005
	CSCAuthenticate        = 2006
	NFYAuthTimer           = 2009
	CSCUnknown3383         = 3383
	CSCUnknown5383         = 5383
)

var opcodeNames = map[uint16]string{
	CSCConnect2Svr:         "Connect2Svr",
	CSCVerifyLinks:         "VerifyLinks",
	CSCAuthAccount:         "AuthAccount",
	CSCForceDisconnect:     "ForceDisconnect",
	NFYSystemMessage:       "SystemMessage",
	NFYServerState:         "ServerState",
	CSCCheckVersion:        "CheckVersion",
	NFYUrlToClient:         "UrlToClient",
	CSCPublicKey:           "PublicKey",
	CSCPreServerEnvRequest: "PreServerEnvRequest",
	NFYDisconnectTimer:     "DisconnectTimer",
	CSCAuthenticate:        "Authenticate",
	NFYAuthTimer:           "AuthTimer",
	CSCUnknown3383:         "Unknown3383",
	CSCUnknown5383:         "Unknown5383",
}
