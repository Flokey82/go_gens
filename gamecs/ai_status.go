package gamecs

import "log"

type CAiStatus struct {
	ai     *CAi
	ape    *CAiPerception
	states map[string]bool
	eat    bool // This is a temp hack to work around the lack of event/command/message system
	sleep  bool
}

func newCAiStatus() *CAiStatus {
	return &CAiStatus{
		states: make(map[string]bool),
	}
}

func (c *CAiStatus) init(ai *CAi) {
	c.ai = ai
	c.ape = ai.CAiPerception
}

func (c *CAiStatus) Update(s *CStatus, delta float64) {
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
	c.states[sHungry] = s.Hunger > 20
}

func (c *CAiStatus) HasFood() bool {
	a := c.ai.w.mgr.GetEntityFromID(c.ai.id)
	if a.CInventory.Find("food") != nil {
		log.Println(a.CInventory.Find("food"))
	}
	return a.CInventory.Find("food") != nil
}

func (c *CAiStatus) Eat() {
	log.Println("Eat!")
	// This is a hack to reset hunger.
	c.eat = true
}

func (c *CAiStatus) Sleep() {
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
