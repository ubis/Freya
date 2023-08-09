package rpc

// Server related RPC's
const (
	ServerRegister = "ServerRegister"
	ServerList     = "ServerList"
)

// Account related RPC's
const (
	AuthCheck       = "AuthCheck"
	UserVerify      = "UserVerify"
	PasswdCheck     = "PasswdCheck"
	OnlineCheck     = "OnlineCheck"
	ForceDisconnect = "ForceDisconnect"
)

// SubPassword related RPC's
const (
	FetchSubPassword  = "FetchSubPassword"
	SetSubPassword    = "SetSubPassword"
	RemoveSubPassword = "RemoveSubPassword"
)

// Character related RPC's
const (
	LoadCharacters    = "LoadCharacters"
	CreateCharacter   = "CreateCharacter"
	DeleteCharacter   = "DeleteCharacter"
	SetSlotOrder      = "SetSlotOrder"
	LoadCharacterData = "LoadCharacterData"
)

const (
	MoveItemEquToEqu = "MoveItemEquToEqu"
	MoveItemEquToInv = "MoveItemEquToInv"
	MoveItemInvToEqu = "MoveItemInvToEqu"
	MoveItemInvToInv = "MoveItemInvToInv"
)
