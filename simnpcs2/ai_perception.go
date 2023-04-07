package simnpcs2

// Perception represents the perception of an AI of the world.
type Perception struct {
	*AI
	Perceived []Entity
}

// newPerception creates a new Perception system for the given AI.
func newPerception(ai *AI) *Perception {
	return &Perception{
		AI:        ai,
		Perceived: make([]Entity, 0, 10),
	}
}

// Update updates the perception of the world.
func (p *Perception) Update(delta float64) {
	// Update all items / entities we can see.
	p.Perceived = p.Perceived[:0]
	for _, e := range p.World.GetEntitiesInRadius(p.Being.Pos(), perceptionRadius) {
		if e.ID() != p.Being.ID() { // Don't perceive ourselves.
			p.Perceived = append(p.Perceived, e)
		}
	}
	p.Perceived = append(p.Perceived, p.World.GetItemsInRadius(p.Being.Pos(), perceptionRadius)...)

	// TODO: If we see a threat, make sure that the AI knows about it.
}
