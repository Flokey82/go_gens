package gamecs

import (
	"github.com/Flokey82/go_gens/aifiver"
	"math/rand"
)

type CAi struct {
	id int
	w  *World

	*CAiPerception
	*CAiScheduler
	*CAiStatus
	*CAiMemory
	*CAiPath
	aifiver.SmallModel
}

func newCAi(w *World, id int) *CAi {
	c := &CAi{
		id:            id,
		w:             w,
		CAiPerception: newCAiPerception(),
		CAiScheduler:  newCAiScheduler(),
		CAiStatus:     newCAiStatus(),
		CAiMemory:     newCAiMemory(),
		CAiPath:       newCAiPath(),
	}
	// Randomize.
	c.SmallModel[aifiver.FactorAgreeableness] = rand.Intn(10) - 5

	c.CAiPerception.init(c)
	c.CAiScheduler.init(c)
	c.CAiStatus.init(c)
	c.CAiMemory.init(c)
	c.CAiPath.init(c)
	return c
}

func (c *CAi) Conflict() bool {
	return c.Get(aifiver.FactorAgreeableness) <= 0
}

func (c *CAi) Update(m *CMovable, s *CStatus, delta float64) {
	// Update perception.
	c.CAiPerception.Update(m, delta)

	// Update states.
	c.CAiStatus.Update(s, delta)

	// Re-evaluate current plans, tasks, or states.
	c.CAiScheduler.Update(delta)

	// Update any path charted.
	c.CAiPath.Update(m, delta)
}
