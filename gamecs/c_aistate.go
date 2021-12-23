package gamecs

import "log"

type CAiState struct {
	ai     *CAi
	ape    *CAiPerception
	states map[string]bool
	eat    bool // This is a temp hack to work around the lack of event/command/message system
	sleep  bool
}

func newCAiState() *CAiState {
	return &CAiState{
		states: make(map[string]bool),
	}
}

func (c *CAiState) init(ai *CAi) {
	c.ai = ai
	c.ape = ai.CAiPerception
}

func (c *CAiState) Update(s *CStatus, delta float64) {
	if c.eat {
		s.Hunger = 0
		c.eat = false
	}
	if c.sleep {
		s.Exhaustion = 0
		c.sleep = false
	}
	c.states[sExhausted] = s.Exhaustion > 10
	c.states[sThreatened] = len(c.ape.Entities) > 0
	c.states[sHungry] = s.Hunger > 10
}

func (c *CAiState) Eat() {
	log.Println("Eat!")
	// This is a hack to reset hunger.
	c.eat = true
}

func (c *CAiState) Sleep() {
	log.Println("Sleep!")
	// This is a hack to reset exhaustion.
	c.sleep = true
}

const (
	sThreatened = "threatened"
	sExhausted  = "exhausted"
	sHungry     = "hungry"
	sIdle       = "idle"
	// sTired
	// sThirsty
)
