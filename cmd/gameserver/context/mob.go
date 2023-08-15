package context

// MobHandler defines the interface for handling mob data in the game world.
type MobHandler interface {
	GetId() int
	GetSpecies() int
	GetHealth() (int, int)
	GetPosition() Position
}
