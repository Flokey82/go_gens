package gamecs

type CStatus struct {
	Health     int
	Exhaustion int
	Hunger     float64
}

const hungerRate = 1.0

func (c *CStatus) Update(delta float64) {
	c.Hunger += hungerRate * delta
}

type State struct {
	Name  string
	Value float64
	Max   float64
}
