package rpc

// Server related RPC's
const (
    ServerRegister  = "ServerRegister"
    ServerList      = "ServerList"
)

// Account related RPC's
const (
    AuthCheck  = "AuthCheck"
    UserVerify = "UserVerify"
)

// SubPassword related RPC's
const (
    FetchSubPassword  = "FetchSubPassword"
    SetSubPassword    = "SetSubPassword"
    RemoveSubPassword = "RemoveSubPassword"
)