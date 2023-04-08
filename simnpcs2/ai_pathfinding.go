package simnpcs2

import (
	"log"

	"github.com/Flokey82/go_gens/vectors"
)

// Pathfinding represents the pathfinding of an AI.
type Pathfinding struct {
	*AI
	Waypoints   []int // Cell indices of the waypoints.
	WaypointIdx int   // Current waypoint.
}

// newPathfinding creates a new pathfinding for the given AI.
func newPathfinding(ai *AI) *Pathfinding {
	return &Pathfinding{
		AI: ai,
	}
}

// SetDestination sets the destination of the AI, and calculates the path.
func (p *Pathfinding) SetDestination(dest *vectors.Vec2) error {
	p.Destination = dest
	start := p.Being.Pos()
	wp, err := p.World.Pathfind(&start, dest)
	if err != nil {
		return err
	}
	p.Waypoints = wp
	p.WaypointIdx = 0
	return nil
}

// Update updates the pathfinding of the AI.
func (p *Pathfinding) Update(delta float64) {
	// No destination? Do nothing.
	if p.Destination == nil || p.WaypointIdx >= len(p.Waypoints) {
		p.Waypoints = nil
		p.WaypointIdx = 0
		return
	}

	// Check if we reached the destination.
	curWP := p.World.CellIdxToPos(p.Waypoints[p.WaypointIdx])
	if curWP.DistanceTo(p.Being.Pos()) < arrivalRadius {
		log.Println("Reached waypoint", p.WaypointIdx)
		p.WaypointIdx++
		if p.WaypointIdx >= len(p.Waypoints) {
			// We reached the destination.
			p.Destination = nil
			p.Being.Speed.SetXY(0, 0)
			return
		}
		curWP = p.World.CellIdxToPos(p.Waypoints[p.WaypointIdx])
	}
	log.Println("Waypoint", p.WaypointIdx, "is at", curWP, "and I'm at", p.Being.Pos(), "my destination is", p.Destination)

	// Calculate the velocity vector.
	velVec := curWP.Sub(p.Being.Pos())

	// TODO: Check if we can still reach the next waypoint or if it's blocked.
	// If not, change the velocity vector accordingly.

	// Limit the velocity vector to the maximum speed.
	if velVec.Len() > maxSpeed {
		velVec.MulWithThis(maxSpeed / velVec.Len())
	}
	p.Being.Speed.SetXY(velVec.X, velVec.Y)
}
