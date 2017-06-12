package packet

import (
    "share/logger"
    "share/network"
)

var log = logger.Instance()

type PacketHandler struct {
    packets network.PacketInfo
}

func (pk *PacketHandler) Init() {
    pk.packets = make(network.PacketInfo)

    // registering packets
    pk.Register(CONNECT2SVR, "Connect2Svr", Connect2Svr)
    pk.Register(VERIFYLINKS, "VerifyLinks", nil)
    pk.Register(AUTHACCOUNT, "AuthAccount", nil)
    pk.Register(SYSTEMMESSG, "SystemMessg", nil)
    pk.Register(SERVERSTATE, "ServerState", nil)
    pk.Register(CHECKVERSION, "CheckVersion", nil)
    pk.Register(URLTOCLIENT, "URLToClient", nil)
    pk.Register(PUBLIC_KEY, "PublicKey", nil)
    pk.Register(PRE_SERVER_ENV_REQUEST, "PreServerEnvRequest", nil)

    for code := range pk.packets {
        log.Debugf("Registered packet: %s(%d)", pk.packets[code].Name, code)
    }
}

func (pk *PacketHandler) Register(code int, name string, method interface{}) {
    pk.packets[code] = &network.PacketData{name, method}
}

func (pk *PacketHandler) Handle(args *network.PacketArgs) {
    if pk.packets[args.Type] == nil {
        // unknown packet received
        log.Errorf("Unknown packet received (Len: %d, Type: %d, Src: %s, UserIdx: %d)",
            args.Length,
            args.Type,
            args.Session.GetEndPnt(),
            args.Session.UserIdx,
        )
        return
    }

    var invoke = pk.packets[args.Type].Method
    invoke.(func(*network.Session, []uint8))(args.Session, *args.Data)
}

func (pk *PacketHandler) Name(code int) string {
    if pk.packets[code] != nil {
        return pk.packets[code].Name
    }

    return "Unknown"
}