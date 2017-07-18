package server

import "share/rpc"

const (
    LOGIN_SERVER_TYPE = 1
    GAME_SERVER_TYPE  = 2
)

type Server struct {
    ServerData
    Client      *rpc.Client
}

type ServerData struct {
    Type            byte
    ServerType      byte
    ServerId        byte
    ChannelId       byte
    PublicIp        string
    PublicPort      uint16
    CurrentUsers    int16
    MaxUsers        int16
}

type RegResponse struct {
    Registered bool
}

type SvrListRequest struct {

}

type SvrListResponse struct {
    Servers []ServerItem
}

type ServerItem struct {
    Id  byte
    Hot byte
    Channels []ChannelItem
}

type ChannelItem struct {
    Id           byte
    Type         byte
    Ip           string
    Port         uint16
    CurrentUsers int16
    MaxUsers     int16
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