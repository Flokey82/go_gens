package gamecs

import (
	"log"
	"math/rand"

	"github.com/Flokey82/go_gens/gamesheet"
)

func randByte() byte {
	return byte(rand.Intn(255))
}

// CompStatus is a component that handles the status or state
// of an entity (Health, stamina, hunger, thirst...).
type CompStatus struct {
	Sleeping bool
	cs       *gamesheet.CharacterSheet
}

// newCompStatus returns a new status component.
func newCompStatus() *CompStatus {
	return &CompStatus{
		cs: gamesheet.New(100, 100, 0, randByte(), randByte(), randByte(), randByte()),
	}
}

// Update updates the status of the entity.
func (c *CompStatus) Update(delta float64) {
	c.cs.Tick(delta)
	if c.Sleeping && c.cs.StatExhaustion.Val <= 0 {
		c.Wake()
	} else if c.Sleeping {
		log.Println("Sleeping!")
	}
}

// Sleep puts the entity to sleep.
func (c *CompStatus) Sleep() {
	c.Sleeping = true
	c.cs.SetState(gamesheet.StateAsleep)
}

// Wake wakes the entity up.
func (c *CompStatus) Wake() {
	c.Sleeping = false
	c.cs.SetState(gamesheet.StateAwake)
}

// Health returns the current health of the entity.
func (c *CompStatus) Health() float64 {
	return float64(c.cs.HP.Value())
}

// MaxHealth returns the maximum health of the entity.
func (c *CompStatus) MaxHealth() float64 {
	return float64(c.cs.HP.Max())
}

// Hunger returns the current hunger of the entity.
func (c *CompStatus) Hunger() float64 {
	return float64(c.cs.StatHunger.Val)
}

// Thirst returns the current thirst of the entity.
func (c *CompStatus) Thirst() float64 {
	return float64(c.cs.StatThirst.Val)
}

// Stress returns the current stress of the entity.
func (c *CompStatus) Stress() float64 {
	return float64(c.cs.StatStress.Val)
}

// Exhaustion returns the current exhaustion of the entity.
func (c *CompStatus) Exhaustion() float64 {
	return float64(c.cs.StatExhaustion.Val)
}

// Dead returns true if the entity is dead.
func (c *CompStatus) Dead() bool {
	return c.cs.Dead
}
