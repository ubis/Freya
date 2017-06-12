package network

type PacketArgs struct {
    Session *Session
    Length  int
    Type    int
    Data    *[]uint8
}