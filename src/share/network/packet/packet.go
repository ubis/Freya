package packet

import (
	"share/log"
	"share/network"
)

// info structure for network packet
type info struct {
	name   string
	method interface{}
}

// List structure for storing network packets
type List struct {
	packets map[int]*info
}

// Add a new network packet
func (l *List) Add(code uint16, name string, method interface{}) {
	if l.packets == nil {
		l.packets = make(map[int]*info)
	}

	ccode := int(code)
	l.packets[ccode] = &info{name: name, method: method}

	pType := "CSC"
	if l.packets[ccode].method == nil {
		pType = "NFY"
	}

	log.Debugf("Registered %s packet: %s(%d)", pType, name, code)
}

// Handle specified network packet
func (l *List) Handle(args *network.PacketArgs) {
	// recover on panic
	defer func(s *network.Session) {
		if err := recover(); err != nil {
			log.Warningf("Panic! Recovered from: %s, src: %s, id: %d",
				l.GetName(args.Type), s.GetEndPnt(), s.Data.AccountId)
			s.Close()
		}
	}(args.Session)

	code := int(args.Packet.Type)
	if l.packets[code] == nil {
		// unknown packet received
		log.Errorf("Unknown packet received (Len: %d, type: %d, src: %s)",
			args.Packet.Size, code, args.Session.GetEndPnt())
		return
	}

	invoke := l.packets[code].method
	if invoke == nil {
		s := args.Session
		t := args.Type
		log.Errorf(
			"Trying to access procedure `%s` (Type: %d, src: %s, id: %d)",
			l.GetName(t), t, s.GetEndPnt(), s.Data.AccountId)
		return
	}

	// invoke packets function
	invoke.(func(*network.Session, *network.Reader))(args.Session, args.Packet)
}

// GetName returns packet's name by code
func (l *List) GetName(code int) string {
	if l.packets[code] != nil {
		return l.packets[code].name
	}

	return "Unknown"
}
