package gamecs

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/Flokey82/go_gens/vectors"
)

// Agent is an independent entity in the world.
type Agent struct {
	id int
	w  *World
	*CompMovable
	*CompStatus
	*CompInventory
	*CompAi
	*Profession
}

// NewChar adds a new Agent to the world and returns it.
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

// newAgent returns a new Agent placed in the world.
func newAgent(w *World) *Agent {
	id := w.mgr.NextID()
	a := &Agent{
		id: id,
		CompMovable: newCompMovable(vectors.NewVec2(
			float64(rand.Intn(w.Height)),
			float64(rand.Intn(w.Width)),
		)),
		CompStatus:    newCompStatus(),
		CompInventory: newCompInventory(w, id, 3),
		CompAi:        newCompAi(w, id),
		w:             w,
	}
	return a
}

// SetProfession assigns a profession to the agent.
// NOTE: This is just for experimentation and will
// probably be refactored into a more generic function
// that allows the extension of the AI.
func (c *Agent) SetProfession(w *World, p *ProfessionType) {
	c.Profession = p.New(w, c, c.GetLocation("home"))

	// We currently only work if we don't have any other worries.
	c.CompAi.CAiScheduler.AddAnyTransition(c.Profession, c.CompAi.Idle)
}

// ID returns the unique identifier for this Agent.
func (c *Agent) ID() int {
	return c.id
}

// Update updates the state of the agent.
func (c *Agent) Update(delta float64) {
	c.CompAi.Update(c.CompMovable, c.CompStatus, delta)
	c.CompMovable.Update(delta)
	c.CompStatus.Update(delta)
}

// Injure causes damage to the agent from a given source.
func (c *Agent) Injure(amount, srcID int) {
	// TODO: Should this be messaging?
	c.cs.TakeDamage(amount)

	// TODO: Change opinion of related tags / terms / individuals.
	log.Println(fmt.Sprintf("%d: remember %d caused %d damage", c.id, srcID, amount))
}
