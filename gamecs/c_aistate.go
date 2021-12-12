package gamecs

type CAiState struct {
	ape    *CAiPerception
	states map[string]bool

	eat bool // This is a temp hack to work around the lack of event/command/message system
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
	if c.eat {
		s.Hunger = 0
		c.eat = false
	}
	c.states[sExhausted] = s.Exhaustion > 100
	c.states[sThreatened] = len(c.ape.Entities) > 0
	c.states[sHungry] = s.Hunger > 10
}

func (c *CAiState) Eat() {
	// This is a hack to reset hunger.
	c.eat = true
}

const (
	sThreatened = "threatened"
	sHungry     = "hungry"
	sExhausted  = "exhausted"
)
