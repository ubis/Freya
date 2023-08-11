package context

// Warp represents a warp point that allows players to move between worlds and locations.
type Warp struct {
	Id       byte
	World    byte
	Location []struct {
		X byte
		Y byte
	}
}
