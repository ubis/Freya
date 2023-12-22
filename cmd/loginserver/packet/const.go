package packet

// Packet opcode values
const (
	CSCConnect2Svr         = 101
	CSCVerifyLinks         = 102
	CSCAuthAccount         = 103
	CSCForceDisconnect     = 109
	NFYSystemMessage       = 120
	NFYServerState         = 121
	CSCCheckVersion        = 122
	NFYUrlToClient         = 128
	CSCPublicKey           = 2001
	CSCPreServerEnvRequest = 2002
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
}
