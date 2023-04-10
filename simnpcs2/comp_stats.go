package simnpcs2

import "math/rand"

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
	Health         float64 // Current health.
	HealthMax      float64 // Maximum health.
	StarvationRate float64 // How fast we starve.
	Starvation     float64 // Current level of starvation.
	ExhaustionRate float64 // How fast we get tired.
	Exhaustion     float64 // Current level of exhaustion.
}

// newCompStats returns a new stats component.
func newCompStats(maxHealth float64) *CompStats {
	return &CompStats{
		Health:         maxHealth,
		HealthMax:      maxHealth,
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