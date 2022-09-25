package gamecs

import (
	"math/rand"

	"github.com/Flokey82/aifiver"
)

// CompAi is the AI component.
type CompAi struct {
	id int
	w  *World

	*CAiPerception
	*CAiScheduler
	*CAiStatus
	*CAiMemory
	*CAiPath
	aifiver.SmallModel
}

// newCompAi returns a new AI component.
func newCompAi(w *World, id int) *CompAi {
	c := &CompAi{
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

// Conflict returns true if the personality indicates low agreeableness.
func (c *CompAi) Conflict() bool {
	return c.Get(aifiver.FactorAgreeableness) <= 0
}

// Update updates the AI state, performs calculations and magic.
func (c *CompAi) Update(m *CompMovable, s *CompStatus, delta float64) {
	// Update perception.
	c.CAiPerception.Update(m, delta)

	// Update states.
	c.CAiStatus.Update(s, delta)

	// Re-evaluate current plans, tasks, or states.
	c.CAiScheduler.Update(delta)

	// Update any path charted.
	c.CAiPath.Update(m, delta)
}
