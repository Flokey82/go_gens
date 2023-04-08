package simnpcs2

import (
	"log"

	"github.com/Flokey82/go_gens/vectors"
)

type Need int

const (
	NeedSurvival Need = iota
	NeedHealth
	NeedSatiation
	NeedRest
	NeedRevenge
	NeedMax
)

// Needs represents the basic needs of an AI.
type Needs struct {
	*AI
	Enemy     *Being        // current enemy
	Needs     [NeedMax]bool // determines which needs need to be met
	Prioities []Need        // determines the order of needs
}

// newNeeds returns a new Needs system.
func newNeeds(ai *AI) *Needs {
	return &Needs{
		AI: ai,
		Prioities: []Need{
			NeedSurvival,
			NeedHealth,
			NeedSatiation,
			NeedRest,
			NeedRevenge,
		},
	}
}

// Update updates the state of the needs.
func (n *Needs) Update(delta float64) {
	// Evaluate basic needs.

	// Observe values changed during the being.act.
	// TODO: Register the observation of values -> needs
	// instead of hardcoding.
	n.Needs[NeedHealth] = n.Being.Health < n.Being.HealthMax
	n.Needs[NeedSurvival] = n.Being.Health < 3
	n.Needs[NeedSatiation] = n.Being.Starvation > HungerStarving
	n.Needs[NeedRest] = n.Being.Exhaustion > ExhaustionTired

	// TODO: Poll recent events for incoming damage or active threats
	// instead of setting it when TakeDamage is called.
	if n.Enemy != nil && n.Enemy.Dead() {
		n.Enemy = nil
		n.Needs[NeedRevenge] = false
	}

	// Act on the needs.
	for _, need := range n.Prioities {
		if n.Needs[need] {
			// n.Being.Log(fmt.Sprintf("needs %d", need))
			n.ActOnNeed(need)
			return
		}
	}
	// n.Being.Log("idle")

	// All needs are met, so we can do other things.

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
}

// ActOnNeed causes the AI to act on the given need.
func (n *Needs) ActOnNeed(need Need) {
	// TODO: Register needs somehow instead of hardcoding
	// needs -> actions.
	switch need {
	case NeedSurvival, NeedHealth:
		// a.Heal()
		log.Println("I'm hurt, heal!")
		n.Being.Health = n.Being.HealthMax
	case NeedSatiation:
		// a.Eat()
		log.Println("Eat!")
		n.Being.Starvation = 0
	case NeedRest:
		// a.Sleep()
		log.Println("Sleep!")
		n.Being.Exhaustion = 0
	case NeedRevenge:
		// TODO: Either pursue or flee the enemy.
		// a.Damage(a.Enemy)
		log.Println("Revenge!")
		n.Enemy.TakeDamage(1, n.Being)
	}
}
