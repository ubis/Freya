package packet

import (
    "share/network"
    "share/logger"
    "loginserver/def"
)

var log = logger.Instance()

var g_ServerConfig   = def.ServerConfig
var g_ServerSettings = def.ServerSettings
var g_RPCHandler     = def.RPCHandler

type PacketHandler struct {
    packets network.PacketInfo
}

// Initializes PacketHandler which registers packets
func (pk *PacketHandler) Init() {
    log.Info("Registering packets...")
    pk.packets = make(network.PacketInfo)

    // registering packets
    pk.Register(CONNECT2SVR, "Connect2Svr", Connect2Svr)
    pk.Register(VERIFYLINKS, "VerifyLinks", nil)
    pk.Register(AUTHACCOUNT, "AuthAccount", AuthAccount)
    pk.Register(SYSTEMMESSG, "SystemMessg", nil)
    pk.Register(SERVERSTATE, "ServerState", nil)
    pk.Register(CHECKVERSION, "CheckVersion", CheckVersion)
    pk.Register(URLTOCLIENT, "URLToClient", nil)
    pk.Register(PUBLIC_KEY, "PublicKey", PublicKey)
    pk.Register(PRE_SERVER_ENV_REQUEST, "PreServerEnvRequest", PreServerEnvRequest)

    for code := range pk.packets {
        var pType = "CSC"

        if pk.packets[code].Method == nil {
            pType = "NFY"
        }

        log.Debugf("Registered %s packet: %s(%d)", pType, pk.packets[code].Name, code)
    }
}

/*
    Registers network packet
    @param  code    packet type
    @param  name    packet name
    @param  method  packet processing method
 */
func (pk *PacketHandler) Register(code uint16, name string, method interface{}) {
    pk.packets[code] = &network.PacketData{name, method}
}

/*
    Handles specified network packet
    @param  args    packet args
 */
func (pk *PacketHandler) Handle(args *network.PacketArgs) {
    // recover on panic
    /*defer func() {
        recover()
        log.Warning("Recovered from:", pk.Name(args.Type))
    }()*/

    if pk.packets[args.Packet.Type] == nil {
        // unknown packet received
        log.Errorf("Unknown packet received (Len: %d, Type: %d, Src: %s, UserIdx: %d)",
            args.Packet.Size,
            args.Packet.Type,
            args.Session.GetEndPnt(),
            args.Session.UserIdx,
        )
        return
    }

    var invoke = pk.packets[args.Packet.Type].Method
    if invoke == nil {
        log.Errorf("Trying to access procedure `%s` (Type: %d, Src: %s, UserIdx: %d)",
            pk.Name(args.Type),
            args.Type,
            args.Session.GetEndPnt(),
            args.Session.UserIdx,
        )
        return;
    }

    invoke.(func(*network.Session, *network.Reader))(args.Session, args.Packet)
}

/*
    Returns packet's name by packet type
    @param  code    packet type
    @return packet name and `Unknown` for un-registered packet
 */
func (pk *PacketHandler) Name(code int) string {
    if pk.packets[uint16(code)] != nil {
        return pk.packets[uint16(code)].Name
    }

    return "Unknown"
}