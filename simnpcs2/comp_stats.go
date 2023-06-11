package simnpcs2

import (
	"fmt"
	"math/rand"
)

const (
	HungerPeckish   = 20
	HungerHungry    = 40
	HungerStarving  = 60
	HungerStarved   = 100
	HungerPerSecond = 1.0
)

const (
	ExhaustionTired     = 10
	ExhaustionExhausted = 40
	ExhaustionDead      = 60
	ExhaustionPerSecond = 1.0
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

// String returns a string representation of the stats.
func (c *CompStats) String() string {
	return fmt.Sprintf("HLT: %.0f/%.0f, STV: %.0f, EXH: %.0f", c.Health, c.HealthMax, c.Starvation, c.Exhaustion)
}

// Update updates the state of the stats.
func (c *CompStats) Update(delta float64) {
	// Sometimes we are clumsy and hurt ourselves.
	if rand.Float64() < 0.01 {
		c.Health -= 10
	}

	// If we are lucky we find some food.
	if rand.Float64() < 0.01 && c.Starvation > 20 {
		c.Starvation -= 1.0
	}

	// If we are lucky we find some rest.
	if rand.Float64() < 0.01 && c.Exhaustion > 20 {
		c.Exhaustion -= 1.0
	}

	// If we are dead we are dead.
	if c.Health < 0 {
		c.Health = 0
	}

	// Update starvation (if there is a starvation rate).
	if c.StarvationRate > 0 {
		c.Starvation += c.StarvationRate * delta
		if c.Starvation > HungerStarved {
			c.Health -= 10 * delta
		}
	}

	// Update exhaustion (if there is an exhaustion rate).
	if c.ExhaustionRate > 0 {
		c.Exhaustion += c.ExhaustionRate * delta
		if c.Exhaustion > ExhaustionDead {
			c.Health -= 10 * delta
		}
	}
}
