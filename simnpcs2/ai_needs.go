package simnpcs2

import (
	"log"
	"math/rand"

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
	Enemy      Entity        // current enemy
	Needs      [NeedMax]bool // determines which needs need to be met
	Prioities  []Need        // determines the order of needs
	Aggression float64       // determines how aggressive the AI is
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
		Aggression: rand.Float64(),
	}
}

// Update updates the state of the needs.
func (n *Needs) Update(delta float64) {
	// Evaluate basic needs.
	n.EvalNeeds()

	// Evaluate threats.
	n.EvalThreats()

	// Act on the needs.
	for _, need := range n.Prioities {
		if n.Needs[need] {
			// n.Being.Log(fmt.Sprintf("needs %d", need))
			n.ActOnNeed(need)
			return
		}
	}

	// All needs are met, so we can do other things.
	// n.Being.Log("idle")

	// Check if we are hungry, thirsty, tired, etc.
	// If so, make sure that the AI knows about it.
	if n.Destination == nil {
		log.Println("I'm bored!")
		// Wander to a random destination.
		if err := n.Pathfinding.SetDestination(&vectors.Vec2{
			X: randFloat(float64(n.World.Width)),
			Y: randFloat(float64(n.World.Height)),
		}, PathfindingModeMoveTo); err != nil {
			log.Println(err)
			n.Pathfinding.Destination = nil
		}
	}
}

func (n *Needs) EvalNeeds() {
	// Observe values changed during the being.act.
	// TODO: Register the observation of values -> needs
	// instead of hardcoding.
	n.Needs[NeedHealth] = n.Being.Health < n.Being.HealthMax
	n.Needs[NeedSurvival] = n.Being.Health < 30
	n.Needs[NeedSatiation] = n.Being.Starvation > HungerStarving
	n.Needs[NeedRest] = n.Being.Exhaustion > ExhaustionTired
}

func (n *Needs) EvalThreats() {
	// TODO: Poll recent events for incoming damage or active threats
	// instead of setting it when TakeDamage is called.
	if evt := n.Being.FindType(EventAttack); evt != nil {
		if evt.Source != nil {
			n.Enemy = evt.Source
			n.Needs[NeedRevenge] = true
		}
	} else if n.Enemy == nil && rand.Float64() < n.Aggression {
		if len(n.Perception.Entities) > 0 {
			n.Enemy = n.Perception.Entities[0].(*AI).Being
			n.Needs[NeedRevenge] = true
		}
	}

	// If we can't see the enemy anymore (or if the enemy is dead), stop fleeing (or attacking).
	if n.Enemy != nil && (n.Enemy.Dead() || !n.Perception.CanSeeEntity(n.Enemy)) {
		n.Enemy = nil
		n.Needs[NeedRevenge] = false
	}
}

// ActOnNeed causes the AI to act on the given need.
func (n *Needs) ActOnNeed(need Need) {
	// TODO: Register needs somehow instead of hardcoding
	// needs -> actions.
	switch need {
	case NeedSurvival:
		// a.Heal()
		log.Println("I'm VERY hurt, heal!", n.Being.Health, n.Being.HealthMax)
		n.Being.Health += 1
		n.ActOnSurvival()
	case NeedHealth:
		// a.Heal()
		log.Println("I'm hurt, heal!", n.Being.Health, n.Being.HealthMax)
		n.Being.Health += 1
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
		n.ActOnRevenge()
	}
}

func (n *Needs) ActOnSurvival() {
	if n.Enemy == nil {
		return
	}

	// Flee from the enemy if we are hurt.
	log.Println("Flee!")
	// NOTE: Find a better way to determine if we are already fleeing.
	if n.AI.Pathfinding.Destination == nil || n.AI.Pathfinding.Mode != PathfindingModeFleeTo {
		// TODO: This should be "flee from enemy"
		pos := n.AI.World.findValidPos()
		n.AI.Pathfinding.SetDestination(&pos, PathfindingModeFleeTo)
	}
}

func (n *Needs) ActOnRevenge() {
	if n.Being.InMeleeRange(n.Enemy) {
		// Attack the enemy.
		dmg := 1 + randFloat(2)

		// Make sure the world knows about the attack.
		// This will be used for nearby entities to react to the attack.
		n.World.Events.Add(EventAttack, n.Being, n.Enemy, &EventAttackData{
			Damage: dmg,
		})

		// Damage the enemy.
		n.Enemy.TakeDamage(dmg, n.Being)
	} else {
		// Chase the enemy.
		pos := n.Enemy.Pos()
		n.AI.Pathfinding.SetDestination(&pos, PathfindingModeChase)
	}
}
