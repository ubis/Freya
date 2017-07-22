package server

import "share/rpc"

const (
    LOGIN_SERVER = 1
    GAME_SERVER  = 2
)

type RegisterReq struct {
    Type            byte
    ServerType      byte
    ServerId        byte
    ChannelId       byte
    PublicIp        uint32
    PublicPort      uint16
    CurrentUsers    uint16
    MaxUsers        uint16
}

type RegisterResp struct {
    Registered bool
}

type Server struct {
    *RegisterReq
    Client       *rpc.Client
}

type ListReq struct {

}

type ListResp struct {
    List []ServerItem
}

type ServerItem struct {
    Id  byte
    Hot byte
    List []ChannelItem
}

type ChannelItem struct {
    Id           byte
    Type         byte
    Ip           uint32
    Port         uint16
    CurrentUsers uint16
    MaxUsers     uint16
}

// for sorting
type ByServer []ServerItem
func (a ByServer) Len() int           { return len(a) }
func (a ByServer) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByServer) Less(i, j int) bool { return a[i].Id < a[j].Id }

type ByChannel []ChannelItem
func (a ByChannel) Len() int           { return len(a) }
func (a ByChannel) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByChannel) Less(i, j int) bool { return a[i].Id < a[j].Id }