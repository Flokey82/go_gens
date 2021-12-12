package gamecs

type CStatus struct {
	Health     int
	Exhaustion int
	Hunger     float64
}

func newCStatus() *CStatus {
	return &CStatus{
		Health:     100,
		Exhaustion: 0,
	}
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
