package gamecs

type CAi struct {
	*CAiPerception
	*CAiScheduler
	*CAiState
	CAiPath
}

func newCAi(w *World) *CAi {
	c := &CAi{
		CAiPerception: newCAiPerception(w),
		CAiScheduler:  newCAiScheduler(),
		CAiState:      newCAiState(),
	}
	c.CAiScheduler.init(c)
	c.CAiState.init(c)
	return c
}

func (c *CAi) Update(m *CMovable, s *CStatus, delta float64) {
	// Update perception.
	c.CAiPerception.Update(m, delta)

	// Update states.
	c.CAiState.Update(s, delta)

	// Re-evaluate current plans, tasks, or states.
	c.CAiScheduler.Update(delta)

	// Update any path charted.
	c.CAiPath.Update(m, delta)
}
