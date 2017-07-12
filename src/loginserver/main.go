package main

import (
    "share/logger"
    "share/network"
    "loginserver/def"
    "loginserver/packet"
    "net"
    "share/lib/rpc2"
    "fmt"
)

var log = logger.Instance()

var g_ServerConfig   = def.ServerConfig
var g_ServerSettings = def.ServerSettings
var g_PacketHandler  = packet.PacketHandler{}

type Args struct{ A, B int }
type Reply int

func main() {
    log.Info("LoginServer init")

    conn, _ := net.Dial("tcp", "127.0.0.1:5000")

    clt := rpc2.NewClient(conn)
    clt.Handle("mult", func(client *rpc2.Client, args *Args, reply *Reply) error {
        *reply = Reply(args.A * args.B)
        return nil
    })
    go clt.Run()

    var rep Reply
    clt.Call("add", Args{1, 2}, &rep)
    fmt.Println("add result:", rep)

    // read config
    g_ServerConfig.Read()

    // set server settings
    g_ServerSettings.XorKeyTable.Init()
    g_ServerSettings.RSA.Init()

    // register events
    RegisterEvents()

    // init packet handler
    g_PacketHandler.Init()

    // create network and start listening for connections
    network.Init(g_ServerConfig.Port, g_ServerSettings.XorKeyTable)
}