package gamecs

import "log"

type CStatus struct {
	Health     int
	MaxHealth  int
	Exhaustion float64
	Hunger     float64
	Sleeping   bool
}

func newCStatus() *CStatus {
	return &CStatus{
		Health:     100,
		MaxHealth:  100,
		Exhaustion: 0,
	}
}

const (
	healthRecovery     = 0.7 // TODO: Implement
	hungerRate         = 0.5
	exhaustionRate     = 0.4
	exhaustionRecovery = 0.9
)

func (c *CStatus) Update(delta float64) {
	log.Println(delta)
	if c.Sleeping {
		log.Println("Sleeping!")
		if c.Health < c.MaxHealth {
			c.Health += 10 // TODO: Take delta in account
		}
		if c.Health > c.MaxHealth {
			c.Health = c.MaxHealth
		}
		// TODO: Reduce exhaustion.
		c.Exhaustion -= exhaustionRecovery * delta // TODO: Take delta in account
		if c.Exhaustion <= 0 {
			c.Exhaustion = 0
			c.Sleeping = false
		}
	} else {
		// TODO: Switch from exhaustion to stamina and regenerate over time.
		c.Exhaustion += exhaustionRate * delta
	}
	c.Hunger += hungerRate * delta
}

func (c *CStatus) Sleep() {
	c.Sleeping = true
}

func (c *CStatus) Dead() bool {
	return c.Health <= 0
}

type State struct {
	Name  string
	Value float64
	Max   float64
}
