package server

type DelUserType int

const (
	DelUserDead DelUserType = iota + 16
	DelUserWarp
	DelUserLogout
	DelUserReturn
	DelUserDisappear
	DelUserNotifyDead
)

type NewUserType int

const (
	NewUserNone NewUserType = iota
	NewUserInit NewUserType = iota + 47
	NewUserWarp
	NewUserMove
	NewUserReset
)
