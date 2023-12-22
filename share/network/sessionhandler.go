package network

type SessionHandler interface {
	Store(any)
	Retrieve() any

	GetUserIdx() uint16
	GetAuthKey() uint32

	GetSeed() uint32
	GetKeyIdx() uint32

	Send(*Writer)
	Close()

	GetEndPnt() string
	GetIp() string
	GetLocalEndPntIp() string
	IsLocal() bool

	AddJob(string, *PeriodicTask)
}
