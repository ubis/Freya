package server

const (
    LOGIN_SERVER_TYPE = 1
    GAME_SERVER_TYPE  = 2
)

type RegRequest struct {
    Type            byte
    ServerType      byte
    ServerId        byte
    ChannelId       byte
    PublicIp        string
    PublicPort      int16
    CurrentUsers    int16
    MaxUsers        int16

}

type RegResponse struct {
    Registered bool
}