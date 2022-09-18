package gamecs

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/Flokey82/go_gens/vectors"
)

type Agent struct {
	id int
	w  *World
	*CMovable
	*CStatus
	*CInventory
	*CAi
	*Profession
}

func (w *World) NewChar() *Agent {
	c := newAgent(w)
	w.mgr.RegisterEntity(c)
	l := newLocation(w, w.mgr.NextID(), vectors.NewVec2(
		float64(rand.Intn(w.Height/2)),
		float64(rand.Intn(w.Width/2)),
	))
	w.mgr.RegisterLocation(l)
	c.SetLocation("home", l)
	return c
}

func newAgent(w *World) *Agent {
	id := w.mgr.NextID()
	a := &Agent{
		id: id,
		CMovable: newCMovable(vectors.NewVec2(
			float64(rand.Intn(w.Height)),
			float64(rand.Intn(w.Width)),
		)),
		CStatus:    newCStatus(),
		CInventory: newCInventory(w, id, 3),
		CAi:        newCAi(w, id),
		w:          w,
	}
	return a
}

func (c *Agent) SetProfession(w *World, p *ProfessionType) {
	c.Profession = p.New(w, c, c.GetLocation("home"))
	// We currently only work if we don't have any other worries.
	c.CAi.CAiScheduler.AddAnyTransition(c.Profession, func() bool {
		return !c.CAi.states[sExhausted] &&
			!c.CAi.states[sThreatened] &&
			!c.CAi.states[sHungry] &&
			!c.CAi.states[sInjured]
	})
}

func (c *Agent) ID() int {
	return c.id
}

func (c *Agent) Update(delta float64) {
	c.CAi.Update(c.CMovable, c.CStatus, delta)
	c.CMovable.Update(delta)
	c.CStatus.Update(delta)
}

func (c *Agent) Injure(amount, srcID int) {
	// TODO: Should this be messaging?
	c.cs.TakeDamage(amount)

	// TODO: Change opinion of related tags / terms / individuals.
	log.Println(fmt.Sprintf("%d: remember %d caused %d damage", c.id, srcID, amount))
}
