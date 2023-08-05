package server

type DelUserType int

const (
	DelUserDead       DelUserType = 16
	DelUserWarp       DelUserType = 17
	DelUserLogout     DelUserType = 18
	DelUserReturn     DelUserType = 19
	DelUserDisappear  DelUserType = 20
	DelUserNotifyDead DelUserType = 21
)
