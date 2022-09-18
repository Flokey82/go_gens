package gamecs

import "log"

type CAiStatus struct {
	ai     *CAi
	ape    *CAiPerception
	states map[string]bool
	eat    bool // This is a temp hack to work around the lack of event/command/message system
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
		s.cs.StatHunger.Val = 0.0
		s.cs.StatThirst.Val = 0.0
		c.eat = false
	}

	// TODO: Improve needs evaluation.
	// Check 'aineeds' repo, add priority list
	c.states[sExhausted] = s.Exhaustion() > 10
	// TODO: Evaluate if the entities we see are actually a threat.
	c.states[sThreatened] = false
	for _, a := range c.ape.Entities {
		// Dead folks are not a threat.
		if a.Dead() {
			continue
		}
		c.states[sThreatened] = true
		break
	}
	// c.states[sInteracting] = len(c.ape.Entities) > 0
	c.states[sHungry] = s.Hunger() > 20
	c.states[sInjured] = s.Health()/s.MaxHealth() <= 0.2
}

func (c *CAiStatus) HasFood() bool {
	a := c.ai.w.mgr.GetEntityFromID(c.ai.id)
	return a.CInventory.Find("food") != nil
}

func (c *CAiStatus) Eat() {
	log.Println("Eat!")
	// This is a hack to reset hunger.
	c.eat = true
}

func (c *CAiStatus) Sleep() {
	log.Println("Sleep!")
	c.ai.w.mgr.GetEntityFromID(c.ai.id).Sleep()
}

const (
	sThreatened = "threatened"
	sExhausted  = "exhausted"
	sHungry     = "hungry"
	sIdle       = "idle"
	sInjured    = "injured"
	// sInteracting = "interacting"
	// sTired
	// sThirsty
)
