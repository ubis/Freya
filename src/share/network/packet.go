package network

import "share/log"

// PacketData structure
type PacketData struct {
	Name   string
	Method interface{}
}

// RegisterPacket registers a new network packet
func (m *Manager) RegisterPacket(code uint16, name string, method interface{}) {
	m.packets[code] = &PacketData{name, method}

	pType := "CSC"

	if m.packets[code].Method == nil {
		pType = "NFY"
	}

	log.Debugf("Registered %s packet: %s(%d)", pType, m.packets[code].Name, code)
}

// HandlePacket handles specified network packet
func (m *Manager) HandlePacket(args *PacketArgs) {
	// recover on panic
	defer func() {
		if err := recover(); err != nil {
			s := args.Session
			log.Warningf("Panic! Recovered from: %s, src: %s, id: %d",
				m.GetPacketName(args.Type), s.GetEndPnt(), s.Data.AccountId)
			s.Close()
		}
	}()

	if m.packets[args.Packet.Type] == nil {
		// unknown packet received
		log.Errorf("Unknown packet received (Len: %d, type: %d, src: %s)",
			args.Packet.Size, args.Packet.Type, args.Session.GetEndPnt())
		return
	}

	invoke := m.packets[args.Packet.Type].Method
	if invoke == nil {
		s := args.Session
		t := args.Type
		log.Errorf("Trying to access procedure `%s` (Type: %d, src: %s, id: %d)",
			m.GetPacketName(t), t, s.GetEndPnt(), s.Data.AccountId)
		return
	}

	// invoke packets function
	invoke.(func(*Session, *Reader))(args.Session, args.Packet)
}

// GetPacketName returns packet's name by code
func (m *Manager) GetPacketName(code int) string {
	if m.packets[uint16(code)] != nil {
		return m.packets[uint16(code)].Name
	}

	return "Unknown"
}
