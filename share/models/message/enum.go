package message

const (
	Normal           = 0
	LoginDuplicate   = 1
	ForceDisconnect  = 2
	Shutdown         = 3
	ShutdownWoNotice = 4 // shutdown without noticing client..?
	War_Cappela      = 5 // broadcast only to capella nation
	War_Procyon      = 6 // broadcast only to procyon nation
)
