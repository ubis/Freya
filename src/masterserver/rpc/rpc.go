package rpc

import (
    "share/logger"
    "share/rpc"
    "masterserver/def"
)

var log = logger.Instance()

var g_RPCHandler      = def.RPCHandler
var g_LoginDatabase   = def.LoginDatabase
var g_ServerManager   = def.ServerManager
var g_DatabaseManager = def.DatabaseManager
var g_DataLoader      = def.DataLoader

func RegisterPackets() {
    log.Info("Registering RPC packets...")

    g_RPCHandler.Register(rpc.ServerRegister, ServerRegister)
    g_RPCHandler.Register(rpc.ServerList, ServerList)

    g_RPCHandler.Register(rpc.AuthCheck, AuthCheck)
    g_RPCHandler.Register(rpc.UserVerify, UserVerify)

    g_RPCHandler.Register(rpc.FetchSubPassword, FetchSubPassword)
    g_RPCHandler.Register(rpc.SetSubPassword, SetSubPassword)
    g_RPCHandler.Register(rpc.RemoveSubPassword, RemoveSubPassword)

    g_RPCHandler.Register(rpc.LoadCharacters, LoadCharacters)
    g_RPCHandler.Register(rpc.CreateCharacter, CreateCharacter)
}