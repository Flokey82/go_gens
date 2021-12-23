package gamecs

import (
	"github.com/Flokey82/go_gens/vectors"
	"math/rand"
)

type Agent struct {
	id int
	*CMovable
	*CStatus
	*CInventory
	*CAi
}

func (w *World) NewChar() *Agent {
	c := newAgent(w)
	w.mgr.RegisterEntity(c)
	return c
}

func newAgent(w *World) *Agent {
	id := w.mgr.NextID()
	a := &Agent{
		id: id,
		CMovable: newCMovable(vectors.Vec2{
			X: float64(rand.Intn(w.Height)),
			Y: float64(rand.Intn(w.Width)),
		}),
		CStatus:    newCStatus(),
		CInventory: newCInventory(w, id),
		CAi:        newCAi(w, id),
	}
	a.SetLocation("home", vectors.Vec2{
		X: float64(rand.Intn(w.Height / 2)),
		Y: float64(rand.Intn(w.Width / 2)),
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
