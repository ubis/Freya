package network

type PacketData struct {
    Name    string
    Method  interface{}
}

type PacketInfo map[uint16]*PacketData