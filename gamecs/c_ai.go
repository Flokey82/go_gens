package gamecs

type CAi struct {
	CAiPerception
	*CAiScheduler
	CAiPath
}

func newCAi() *CAi {
	c := &CAi{
		CAiScheduler: newCAiScheduler(),
	}

	c.CAiScheduler.init(&c.CAiPath, &c.CAiPerception)
	return c
}

func (c *CAi) Update(w *World, m *CMovable, delta float64) {
	// Update perception.
	c.CAiPerception.Update(w, m, delta)

	// Re-evaluate current plans, tasks, or states.
	c.CAiScheduler.Update(m, &c.CAiPath, &c.CAiPerception, delta)

	// Update any path charted.
	c.CAiPath.Update(m, delta)
}
