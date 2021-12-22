package gamecs

import (
	"github.com/Flokey82/go_gens/vectors"
	"math/rand"
)

type Agent struct {
	*CMovable
	*CStatus
	*CAi

	id int
}

func (w *World) NewChar() *Agent {
	c := newAgent(w)
	w.Add(c)
	return c
}

func newAgent(w *World) *Agent {
	a := &Agent{
		id: w.mgr.NextID(),
		CMovable: newCMovable(vectors.Vec2{
			X: float64(rand.Intn(w.Height)),
			Y: float64(rand.Intn(w.Width)),
		}),
		CStatus: newCStatus(),
		CAi:     newCAi(w),
	}
	a.SetLocation("home", vectors.Vec2{
		X: float64(rand.Intn(w.Height)),
		Y: float64(rand.Intn(w.Width)),
	})
	return a
}

func (c *Agent) ID() int {
	return c.id
}

func (c *Agent) Update(delta float64) {
	c.CAi.Update(c.CMovable, c.CStatus, delta)
	c.CMovable.Update(delta)
	c.CStatus.Update(delta)
}
