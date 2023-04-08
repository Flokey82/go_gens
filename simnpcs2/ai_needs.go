package simnpcs2

import (
	"log"

	"github.com/Flokey82/go_gens/vectors"
)

// Needs represents the basic needs of an AI.
type Needs struct {
	*AI
}

// newNeeds returns a new Needs system.
func newNeeds(ai *AI) *Needs {
	return &Needs{
		AI: ai,
	}
}

// Update updates the state of the needs.
func (n *Needs) Update(delta float64) {
	// Check if we are hungry, thirsty, tired, etc.
	// If so, make sure that the AI knows about it.
	if n.Destination == nil {
		log.Println("I'm bored!")
		// Wander to a random destination.
		if err := n.Pathfinding.SetDestination(&vectors.Vec2{
			X: randFloat(float64(n.World.Width)),
			Y: randFloat(float64(n.World.Height)),
		}); err != nil {
			log.Println(err)
			n.Pathfinding.Destination = nil
		}
	}

	// Check if we are hurt.
	if n.Being.Health <= 0 {
		log.Println("I'm dead!")
		// Heal up to a random amount.
		n.Being.Health = randFloat(n.Being.HealthMax)
	} else if n.Being.Health < n.Being.HealthMax/2 {
		log.Println("I'm hurt!")
	}
}
