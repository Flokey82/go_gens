package gamecs

type CStatus struct {
	Health     int
	Exhaustion float64
	Hunger     float64
}

func newCStatus() *CStatus {
	return &CStatus{
		Health:     100,
		Exhaustion: 0,
	}
}

const (
	hungerRate     = 0.5
	exhaustionRate = 0.4
)

func (c *CStatus) Update(delta float64) {
	c.Hunger += hungerRate * delta
	// TODO: Switch from exhaustion to stamina and regenerate over time.
	c.Exhaustion += exhaustionRate * delta
}

type State struct {
	Name  string
	Value float64
	Max   float64
}
