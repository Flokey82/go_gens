package gamecs

import (
	"log"
	"math/rand"

	"github.com/Flokey82/go_gens/gamesheet"
)

func randByte() byte {
	return byte(rand.Intn(255))
}

type CStatus struct {
	Sleeping bool
	cs       *gamesheet.CharacterSheet
}

func newCStatus() *CStatus {
	return &CStatus{
		cs: gamesheet.New(100, 100, 0, randByte(), randByte(), randByte(), randByte()),
	}
}

func (c *CStatus) Update(delta float64) {
	c.cs.Tick(delta)
	if c.Sleeping && c.cs.StatExhaustion.Val <= 0 {
		c.Wake()
	} else if c.Sleeping {
		log.Println("Sleeping!")
	}
}

func (c *CStatus) Sleep() {
	c.Sleeping = true
	c.cs.SetState(gamesheet.StateAsleep)
}

func (c *CStatus) Wake() {
	c.Sleeping = false
	c.cs.SetState(gamesheet.StateAwake)
}

func (c *CStatus) Health() float64 {
	return float64(c.cs.HP.Value())
}

func (c *CStatus) MaxHealth() float64 {
	return float64(c.cs.HP.Max())
}

func (c *CStatus) Hunger() float64 {
	return float64(c.cs.StatHunger.Val)
}

func (c *CStatus) Thirst() float64 {
	return float64(c.cs.StatThirst.Val)
}

func (c *CStatus) Stress() float64 {
	return float64(c.cs.StatStress.Val)
}

func (c *CStatus) Exhaustion() float64 {
	return float64(c.cs.StatExhaustion.Val)
}

func (c *CStatus) Dead() bool {
	return c.cs.Dead
}
