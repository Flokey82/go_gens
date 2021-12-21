package gamecs

type iCStatus interface {
	GetCStatus() *CStatus
}

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

func (c *CStatus) GetCStatus() *CStatus {
	return c
}

const (
	hungerRate     = 1.0
	exhaustionRate = 0.2
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
