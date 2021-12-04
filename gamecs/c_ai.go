package gamecs

type CAi struct {
	*CAiPerception
	*CAiScheduler
	*CAiState
	CAiPath

	w *World
}

func newCAi(w *World) *CAi {
	c := &CAi{
		CAiPerception: newCAiPerception(w),
		CAiScheduler:  newCAiScheduler(),
		CAiState:      newCAiState(),
		w:             w,
	}

	c.CAiScheduler.init(&c.CAiPath, c.CAiPerception, c.CAiState)
	c.CAiState.init(c.CAiPerception)
	return c
}

func (c *CAi) Update(m *CMovable, s *CStatus, delta float64) {
	// Update perception.
	c.CAiPerception.Update(m, delta)

	// Update states.
	c.CAiState.Update(m, s, delta)

	// Re-evaluate current plans, tasks, or states.
	c.CAiScheduler.Update(m, delta)

	// Update any path charted.
	c.CAiPath.Update(m, delta)
}
