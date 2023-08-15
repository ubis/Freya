package game

import (
	"math"

	"github.com/beefsack/go-astar"
	"github.com/ubis/Freya/cmd/gameserver/context"
)

// WayPoint represents a point in a game world and holds its coordinates.
type WayPoint struct {
	X int
	Y int

	world context.WorldHandler
}

// Get returns the coordinates (X, Y) of the waypoint.
func (wp WayPoint) Get() (int, int) {
	return wp.X, wp.Y
}

// World returns the world associated with the waypoint.
func (wp WayPoint) World() context.WorldHandler {
	return wp.world
}

// PathNeighborCost returns the cost of traveling from this waypoint to a
// neighboring waypoint.
// Currently, this is set to a constant value of 1.0 for all neighbors.
func (wp WayPoint) PathNeighborCost(point astar.Pather) float64 {
	return 1.0
}

// PathEstimatedCost calculates and returns an estimate of the least cost from
// this waypoint to another.
// It uses the maximum change in either X or Y coordinate as the heuristic.
func (wp WayPoint) PathEstimatedCost(point astar.Pather) float64 {
	to := point.(WayPoint)

	deltaX := AbsInt(to.X - wp.X)
	deltaY := AbsInt(to.Y - wp.Y)

	return math.Max(float64(deltaX), float64(deltaY))
}

// PathNeighbors identifies and returns the neighboring waypoints around this waypoint.
// It considers 8 directions: North, South, East, West, and the four diagonals.
// It then filters out any non-movable points based on the world's IsMovable method.
func (wp WayPoint) PathNeighbors() []astar.Pather {
	x, y := wp.X, wp.Y

	var neighbors []astar.Pather

	directions := []struct {
		dx int
		dy int
	}{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		{1, -1}, {-1, 1}, {1, 1}, {-1, -1},
	}

	for _, dir := range directions {
		newX, newY := x+dir.dx, y+dir.dy

		if !wp.world.IsMovable(newX, newY) {
			continue
		}

		waypoint := WayPoint{X: newX, Y: newY, world: wp.world}
		neighbors = append(neighbors, waypoint)

	}

	return neighbors
}
