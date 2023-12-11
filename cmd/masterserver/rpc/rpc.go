package rpc

import (
	"github.com/ubis/Freya/cmd/masterserver/def"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/rpc"
)

var g_RPCHandler = def.RPCHandler
var g_LoginDatabase = def.LoginDatabase
var g_ServerManager = def.ServerManager
var g_DatabaseManager = def.DatabaseManager
var g_DataLoader = def.DataLoader

func RegisterPackets() {
	log.Info("Registering RPC packets...")

	g_RPCHandler.Register(rpc.ServerRegister, ServerRegister)
	g_RPCHandler.Register(rpc.ServerList, ServerList)

	g_RPCHandler.Register(rpc.AuthCheck, AuthCheck)
	g_RPCHandler.Register(rpc.UserVerify, UserVerify)
	g_RPCHandler.Register(rpc.PasswdCheck, PasswdCheck)
	g_RPCHandler.Register(rpc.ForceDisconnect, ForceDisconnect)

	g_RPCHandler.Register(rpc.FetchSubPassword, FetchSubPassword)
	g_RPCHandler.Register(rpc.SetSubPassword, SetSubPassword)
	g_RPCHandler.Register(rpc.RemoveSubPassword, RemoveSubPassword)

	g_RPCHandler.Register(rpc.LoadCharacters, LoadCharacters)
	g_RPCHandler.Register(rpc.CreateCharacter, CreateCharacter)
	g_RPCHandler.Register(rpc.DeleteCharacter, DeleteCharacter)
	g_RPCHandler.Register(rpc.SetSlotOrder, SetSlotOrder)
	g_RPCHandler.Register(rpc.LoadCharacterData, LoadCharacterData)

	g_RPCHandler.Register(rpc.EquipItem, EquipItem)
	g_RPCHandler.Register(rpc.UnEquipItem, UnEquipItem)
	g_RPCHandler.Register(rpc.MoveEquipmentItem, MoveEquipmentItem)

	g_RPCHandler.Register(rpc.AddItem, AddItem)
	g_RPCHandler.Register(rpc.StackItem, StackItem)
	g_RPCHandler.Register(rpc.RemoveItem, RemoveItem)
	g_RPCHandler.Register(rpc.SwapItem, SwapItem)
	g_RPCHandler.Register(rpc.MoveItem, MoveItem)
}
