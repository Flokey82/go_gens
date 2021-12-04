package gamecs

import (
	"fmt"
	"log"
)

const perceptionDist = 4.0

type CAiPerception struct {
	Entities []*Character
	w        *World
}

func newCAiPerception(w *World) *CAiPerception {
	return &CAiPerception{
		w: w,
	}
}

func (c *CAiPerception) Update(m *CMovable, delta float64) {
	c.Entities = nil
	for _, ce := range c.w.c {
		if &ce.CMovable == m {
			continue
		}
		if calcDist(ce.CMovable.Pos, m.Pos) < perceptionDist {
			c.Entities = append(c.Entities, ce)
		}
	}
	log.Println(fmt.Sprintf("saw %d Entities", len(c.Entities)))
}
