package simnpcs2

import (
	"log"
	"math/rand"

	"github.com/Flokey82/go_gens/vectors"
)

// CompMoveable is a movable component.
type CompMoveable struct {
	Pos   vectors.Vec2
	Speed vectors.Vec2
}

// newCompMoveable returns a new moveable component.
func newCompMoveable(pos vectors.Vec2) *CompMoveable {
	return &CompMoveable{
		Pos: pos,
	}
}

// Update moves the position in the component by the speed within the
// given time.
func (c *CompMoveable) Update(delta float64) {
	c.Pos.AddToThis(c.Speed)
	log.Println("speed:", c.Speed)

	// TODO: Check for collisions and provide a way to react to them.
}

const (
	HungerPeckish   = 20
	HungerHungry    = 40
	HungerStarving  = 60
	HungerStarved   = 100
	HungerPerSecond = 0.1
)

const (
	ExhaustionTired     = 20
	ExhaustionExhausted = 40
	ExhaustionDead      = 60
	ExhaustionPerSecond = 0.1
)

// CompStats is a stats component.
type CompStats struct {
	Health         float64
	HealthMax      float64
	StarvationRate float64 // How fast we starve.
	Starvation     float64 // Current level of starvation.
	ExhaustionRate float64 // How fast we get tired.
	Exhaustion     float64 // Current level of exhaustion.
}

// newCompStats returns a new stats component.
func newCompStats() *CompStats {
	return &CompStats{
		Health:         100,
		HealthMax:      100,
		StarvationRate: HungerPerSecond,
		ExhaustionRate: ExhaustionPerSecond,
	}
}

// Update updates the state of the stats.
func (c *CompStats) Update(delta float64) {
	// Sometimes we are clumsy and hurt ourselves.
	if rand.Float64() < 0.01 {
		c.Health -= 10
	}
	// If we are lucky we find some food.
	if rand.Float64() < 0.01 {
		c.Starvation -= 10
	}
	// If we are lucky we find some rest.
	if rand.Float64() < 0.01 {
		c.Exhaustion -= 10
	}

	// If we are dead we are dead.
	if c.Health < 0 {
		c.Health = 0
	}

	// Update starvation (if there is a starvation rate).
	if c.StarvationRate > 0 {
		c.Starvation += c.StarvationRate * delta
		if c.Starvation > HungerStarved {
			c.Health -= 10
		}
	}

	// Update exhaustion (if there is an exhaustion rate).
	if c.ExhaustionRate > 0 {
		c.Exhaustion += c.ExhaustionRate * delta
		if c.Exhaustion > ExhaustionDead {
			c.Health -= 10
		}
	}
}

// StarvationLevel returns the current level of starvation.
func (c *CompStats) StarvationLevel() float64 {
	if c.Starvation >= HungerStarved {
		return HungerStarved
	}
	if c.Starvation >= HungerStarving {
		return HungerStarving
	}
	if c.Starvation >= HungerHungry {
		return HungerHungry
	}
	if c.Starvation >= HungerPeckish {
		return HungerPeckish
	}
	return 0
}

// ExhaustionLevel returns the current level of exhaustion.
func (c *CompStats) ExhaustionLevel() float64 {
	if c.Exhaustion >= ExhaustionDead {
		return ExhaustionDead
	}
	if c.Exhaustion >= ExhaustionExhausted {
		return ExhaustionExhausted
	}
	if c.Exhaustion >= ExhaustionTired {
		return ExhaustionTired
	}
	return 0
}
