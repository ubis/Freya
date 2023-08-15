package context

// WayPointHandler defines the interface for handling waypoint data in the
// game world.
type WayPointHandler interface {
	Get() (int, int)
	World() WorldHandler
}

// Position represents object's position data in the game world.
type Position struct {
	IsMoving        bool
	IsDeadReckoning bool

	InitialX int
	InitialY int
	CurrentX int
	CurrentY int
	FinalX   int
	FinalY   int

	WayPoints       []WayPointHandler
	CurrentWayPoint int

	MoveBegin int64
	Speed     float32
	Base      float32
	Distance  float32
	Sin       float32
	Cos       float32
}
