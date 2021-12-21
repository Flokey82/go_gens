package gamecs

import (
	"fmt"
	"log"
)

const perceptionDist = 4.0

type CAiPerception struct {
	Entities []*Agent
	Items    []*Item
	w        *World
	maxDist  float64
}

func newCAiPerception(w *World) *CAiPerception {
	return &CAiPerception{
		w:       w,
		maxDist: perceptionDist,
	}
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
	for _, ce := range c.w.c {
		if ce.CMovable != m && calcDist(ce.CMovable.Pos, m.Pos) < c.maxDist {
			c.Entities = append(c.Entities, ce)
		}
	}

	// Update perceived items.
	c.Items = nil
	for _, it := range c.w.items {
		if it.Location == LocWorld && calcDist(it.Pos, m.Pos) < c.maxDist {
			c.Items = append(c.Items, it)
		}
	}
	log.Println(fmt.Sprintf("saw %d Entities, %d Items", len(c.Entities), len(c.Items)))
}
