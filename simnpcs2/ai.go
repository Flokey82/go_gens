package simnpcs2

import (
	"github.com/Flokey82/go_gens/vectors"
)

const (
	perceptionRadius = 1.0
	arrivalRadius    = 0.01
	maxSpeed         = 0.1
)

type AI struct {
	World       *World        // underlying world
	Being       *Being        // underlying being
	Perception  *Perception   // perception of the world
	Pathfinding *Pathfinding  // steering behaviors
	Needs       *Needs        // basic needs
	Destination *vectors.Vec2 // current destination
}

// NewAI returns a new AI.
func NewAI(being *Being) *AI {
	ai := &AI{
		World: being.World,
		Being: being,
	}
	ai.Perception = newPerception(ai)
	ai.Pathfinding = newPathfinding(ai)
	ai.Needs = newNeeds(ai)
	return ai
}

// ID returns the ID of the being this AI is controlling.
func (a *AI) ID() int64 {
	return a.Being.ID()
}

// Pos returns the position of the being this AI is controlling.
func (a *AI) Pos() vectors.Vec2 {
	return a.Being.Pos()
}

// Type returns the type of the being this AI is controlling.
func (a *AI) Type() EntityType {
	return a.Being.Type()
}

// Update updates the state of the AI and its underlying components.
func (a *AI) Update(delta float64) {
	// First, update our knowledge of the world.
	a.Perception.Update(delta)

	// Re-evaluate our needs. Instincts trump everything.
	a.Needs.Update(delta)

	// Update our pathfinding.
	a.Pathfinding.Update(delta)

	// Update our being.
	a.Being.Update(delta)
}
