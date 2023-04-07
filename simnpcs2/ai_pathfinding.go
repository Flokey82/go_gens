package simnpcs2

// Pathfinding represents the pathfinding of an AI.
type Pathfinding struct {
	*AI
}

// newPathfinding creates a new pathfinding for the given AI.
func newPathfinding(ai *AI) *Pathfinding {
	return &Pathfinding{
		AI: ai,
	}
}

// Update updates the pathfinding of the AI.
func (p *Pathfinding) Update(delta float64) {
	// No destination? Do nothing.
	if p.Destination == nil {
		return
	}

	// TODO: Check if we can still reach the next waypoint.
	// If not, change the velocity vector to the next waypoint.

	// Check if we reached the destination.
	if p.Destination.DistanceTo(p.Being.Pos()) < arrivalRadius {
		// We reached the destination.
		p.Destination = nil
		p.Being.Speed.SetXY(0, 0)
		return
	}

	// Calculate the velocity vector.
	velVec := p.Destination.Sub(p.Being.Pos())
	if velVec.Len() > maxSpeed {
		velVec.MulWithThis(maxSpeed / velVec.Len())
	}
	p.Being.Speed.SetXY(velVec.X, velVec.Y)
}
