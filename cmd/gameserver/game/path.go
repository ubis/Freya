package game

import (
	"math"
	"time"

	"github.com/ubis/Freya/cmd/gameserver/context"
)

// Helper function to compute distance, sine, and cosine between two points.
func computeDSC(x1, y1, x2, y2 int) (float32, float32, float32) {
	deltaX := float64(x2 - x1)
	deltaY := float64(y2 - y1)

	distance := float32(math.Hypot(deltaX, deltaY))
	sin := float32(deltaY) / distance
	cos := float32(deltaX) / distance

	return distance, sin, cos
}

// OpenDeadReckoning initializes the dead reckoning process for a given position.
func OpenDeadReckoning(pos *context.Position) {
	pos.Distance, pos.Sin, pos.Cos =
		computeDSC(pos.InitialX, pos.InitialY, pos.FinalX, pos.FinalY)

	pos.IsMoving = true
	pos.IsDeadReckoning = true
	pos.MoveBegin = time.Now().UnixMilli()
}

// DeadReckoning performs the dead reckoning algorithm to predict entity's future
// position based on its current state and motion.
func DeadReckoning(pos *context.Position) {
	timeElapsed := float32(time.Now().UnixMilli()-pos.MoveBegin) / 1000
	travelledDistance := pos.Speed * timeElapsed
	remainingDistance := travelledDistance - pos.Base

	if remainingDistance < 0 {
		return
	}

	for remainingDistance >= pos.Distance {
		pos.Base += pos.Distance
		pos.CurrentWayPoint++

		if pos.CurrentWayPoint+1 >= len(pos.WayPoints) {
			pos.IsDeadReckoning = false
			pos.CurrentX, pos.CurrentY = pos.FinalX, pos.FinalY
			return
		}

		// Update the initial and final coordinates for the next waypoint
		pos.InitialX, pos.InitialY = pos.FinalX, pos.FinalY
		pos.FinalX, pos.FinalY = pos.WayPoints[pos.CurrentWayPoint+1].Get()
		pos.Distance, pos.Sin, pos.Cos =
			computeDSC(pos.InitialX, pos.InitialY, pos.FinalX, pos.FinalY)

		remainingDistance -= pos.Distance
	}

	moveX := int(remainingDistance * pos.Cos)
	moveY := int(remainingDistance * pos.Sin)

	pos.CurrentX = pos.InitialX + moveX
	pos.CurrentY = pos.InitialY + moveY
}
