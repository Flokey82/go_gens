package simnpcs2

// Perception represents the perception of an AI of the world.
type Perception struct {
	*AI
	Entities []Entity
	Items    []*Item
}

// newPerception creates a new Perception system for the given AI.
func newPerception(ai *AI) *Perception {
	return &Perception{
		AI:       ai,
		Entities: make([]Entity, 0, 10),
		Items:    make([]*Item, 0, 10),
	}
}

// Update updates the perception of the world.
func (p *Perception) Update(delta float64) {
	// Update all entities we can see.
	p.Entities = p.Entities[:0]
	for _, e := range p.World.GetEntitiesInRadius(p.Being.Pos(), perceptionRadius) {
		if e.ID() != p.Being.ID() { // Don't perceive ourselves.
			p.Entities = append(p.Entities, e)
		}
	}

	// Update all items we can see.
	p.Items = p.Items[:0]
	for _, i := range p.World.GetItemsInRadius(p.Being.Pos(), perceptionRadius) {
		p.Items = append(p.Items, i)
	}

	// TODO: If we see a threat, make sure that the AI knows about it.
}
