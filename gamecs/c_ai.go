package gamecs

type CAi struct {
	*CAiPerception
	*CAiScheduler
	CAiPath

	w *World
}

func newCAi(w *World) *CAi {
	c := &CAi{
		CAiPerception: newCAiPerception(w),
		CAiScheduler:  newCAiScheduler(),
		w:             w,
	}

	c.CAiScheduler.init(&c.CAiPath, c.CAiPerception)
	return c
}

func (c *CAi) Update(m *CMovable, delta float64) {
	// Update perception.
	c.CAiPerception.Update(m, delta)

	// Re-evaluate current plans, tasks, or states.
	c.CAiScheduler.Update(m, &c.CAiPath, c.CAiPerception, delta)

	// Update any path charted.
	c.CAiPath.Update(m, delta)
}
