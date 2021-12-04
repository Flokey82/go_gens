package gamecs

type CAiState struct {
	ape    *CAiPerception
	states map[string]bool
}

func newCAiState() *CAiState {
	return &CAiState{
		states: make(map[string]bool),
	}
}

func (c *CAiState) init(ape *CAiPerception) {
	c.ape = ape
}

func (c *CAiState) Update(m *CMovable, s *CStatus, delta float64) {
	c.states[sExhausted] = s.Exhaustion > 100
	c.states[sThreatened] = len(c.ape.Entities) > 0
	c.states[sHungry] = s.Hunger > 100
}

const (
	sThreatened = "threatened"
	sHungry     = "hungry"
	sExhausted  = "exhausted"
)
