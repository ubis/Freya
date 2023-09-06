package context

type ItemHandler interface {
	GetId() int32
	GetOwner() int32
	GetKind() uint32
	GetOption() int32
	GetPosition() (uint16, uint16)
	GetKey() uint16
}
