package gamecs

import (
	"fmt"
	"log"
)

const perceptionDist = 4.0

type CAiPerception struct {
	ai       *CAi
	w        *World
	Entities []*Agent
	Items    []*Item
	maxDist  float64
}

func newCAiPerception() *CAiPerception {
	return &CAiPerception{
		maxDist: perceptionDist,
	}
}

func (c *CAiPerception) init(ai *CAi) {
	c.ai = ai
	c.w = ai.w
}

func (c *CAiPerception) CanSeeEntity(it *Agent) bool {
	for _, cit := range c.Entities {
		if it == cit {
			return true
		}
	}
	return false
}

func (c *CAiPerception) CanSee(it *Item) bool {
	for _, cit := range c.Items {
		if it == cit {
			return true
		}
	}
	return false
}

func (c *CAiPerception) Update(m *CMovable, delta float64) {
	// Update perceived agents.
	c.Entities = nil
	for _, ce := range c.w.mgr.Entities() {
		if ce.CMovable != m && calcDist(ce.CMovable.Pos, m.Pos) < c.maxDist {
			c.Entities = append(c.Entities, ce)
		}
	}

	// Update perceived items.
	c.Items = nil
	for _, it := range c.w.mgr.Items() {
		if it.Location == LocWorld && calcDist(it.Pos, m.Pos) < c.maxDist {
			c.Items = append(c.Items, it)
		}
	}
	log.Println(fmt.Sprintf("saw %d Entities, %d Items", len(c.Entities), len(c.Items)))
}
