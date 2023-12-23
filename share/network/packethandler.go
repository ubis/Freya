package network

import (
	"runtime/debug"

	"github.com/ubis/Freya/share/log"
)

type PacketData struct {
	Name   string
	Method interface{}
}

type PacketInfo map[uint16]*PacketData

type PacketHandler struct {
	packets PacketInfo
}

// Initializes PacketHandler
func (pk *PacketHandler) Init() {
	pk.packets = make(PacketInfo)
}

// Registers a new network packet
func (pk *PacketHandler) Register(code uint16, name string, method interface{}) {
	pk.packets[code] = &PacketData{name, method}

	var pType = "CSC"

	if pk.packets[code].Method == nil {
		pType = "NFY"
	}

	log.Debugf("Registered %s packet: %s(%d)", pType, pk.packets[code].Name, code)
}

// Handles specified network packet
func (pk *PacketHandler) Handle(args *PacketArgs) {
	// recover on panic
	defer func() {
		if err := recover(); err != nil {
			log.Warningf("Panic! Recovered from: %s, src: %s",
				pk.Name(args.Type), args.Session.GetEndPnt(),
			)
			log.Error(err)
			log.Error(string(debug.Stack()))

			args.Session.Close()
		}
	}()

	if pk.packets[args.Reader.Type] == nil {
		// unknown packet received
		log.Errorf("Unknown packet received (Len: %d, type: %d, src: %s)",
			args.Reader.Size, args.Reader.Type, args.Session.GetEndPnt(),
		)

		DumpPacket(args.Reader)

		return
	}

	var invoke = pk.packets[args.Reader.Type].Method
	if invoke == nil {
		log.Errorf("Trying to access procedure `%s` (Type: %d, src: %s)",
			pk.Name(args.Type), args.Type, args.Session.GetEndPnt(),
		)

		return
	}

	// invoke packets function
	invoke.(func(*Session, *Reader))(args.Session, args.Reader)
}

// Returns packet's name by packet type
func (pk *PacketHandler) Name(code int) string {
	if pk.packets[uint16(code)] != nil {
		return pk.packets[uint16(code)].Name
	}

	return "Unknown"
}
