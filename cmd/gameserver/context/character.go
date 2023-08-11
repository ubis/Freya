package context

import "github.com/ubis/Freya/share/network"

// GetCharId retrieves the character ID from the session's context.
func GetCharId(session *network.Session) (int32, error) {
	ctx, err := Parse(session)
	if err != nil {
		return 0, err
	}

	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()

	id := ctx.Char.Id

	return id, nil
}

// GetCharPosition retrieves the character's position (x, y) from the session's context.
func GetCharPosition(session *network.Session, x, y *byte) error {
	ctx, err := Parse(session)
	if err != nil {
		return err
	}

	ctx.Mutex.RLock()
	defer ctx.Mutex.RUnlock()

	*x = ctx.Char.X
	*y = ctx.Char.Y

	return nil
}
