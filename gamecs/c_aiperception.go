package gamecs

import (
	"fmt"
	"log"
)

const perceptionDist = 4.0

type CAiPerception struct {
	Entities []*Character
}

func (c *CAiPerception) Update(w *World, m *CMovable, delta float64) {
	c.Entities = nil
	for _, ce := range w.c {
		if &ce.CMovable == m {
			continue
		}
		if calcDist(ce.CMovable.Pos, m.Pos) < perceptionDist {
			c.Entities = append(c.Entities, ce)
		}
	}
	log.Println(fmt.Sprintf("saw %d Entities", len(c.Entities)))
}
