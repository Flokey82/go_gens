package gamecs

import (
	"fmt"
	"log"
	"sort"

	"github.com/Flokey82/go_gens/vectors"
)

const perceptionDist = 4.0

type CAiPerception struct {
	ai       *CompAi
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

// init initializes the perception using the CAi struct.
func (c *CAiPerception) init(ai *CompAi) {
	c.ai = ai
	c.w = ai.w
}

// CanSeeEntity returns true if we can percieve the given entity (Agent).
// TODO: Deduplicate with CanSee.
func (c *CAiPerception) CanSeeEntity(it *Agent) bool {
	for _, cit := range c.Entities {
		if it == cit {
			return true
		}
	}
	return false
}

// CanSee returns true if we can see the given item.
// TODO: Deduplicate with CanSeeEntity.
func (c *CAiPerception) CanSee(it *Item) bool {
	for _, cit := range c.Items {
		if it == cit {
			return true
		}
	}
	return false
}

// Update updates the list of visible items / entities.
func (c *CAiPerception) Update(m *CompMovable, delta float64) {
	// TODO: Send notifications on newly spotted entities and items and those we lost sight of.

	// Update perceived agents.
	// TODO: Deduplicate with percieved items.
	c.Entities = nil
	for _, ce := range c.w.mgr.Entities() {
		if ce.CompMovable != m && vectors.Dist2(ce.CompMovable.Pos, m.Pos) < c.maxDist {
			c.Entities = append(c.Entities, ce)
		}
	}
	sort.Slice(c.Entities, func(i, j int) bool {
		return vectors.Dist2(c.Entities[i].Pos, m.Pos) < vectors.Dist2(c.Entities[j].Pos, m.Pos)
	})

	// Update perceived items.
	// TODO: Deduplicate with percieved agents.
	c.Items = nil
	for _, it := range c.w.mgr.Items() {
		if it.Location == LocWorld && vectors.Dist2(it.Pos, m.Pos) < c.maxDist {
			c.Items = append(c.Items, it)
		}
	}
	sort.Slice(c.Items, func(i, j int) bool {
		return vectors.Dist2(c.Items[i].Pos, m.Pos) < vectors.Dist2(c.Items[j].Pos, m.Pos)
	})
	log.Println(fmt.Sprintf("saw %d Entities, %d Items", len(c.Entities), len(c.Items)))
}
