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
// TODO: Based on insight, we should be able to gain an understanding of the
// motives of other beings. This will allow us to make better decisions, like
// whether to attack or flee from a threat.
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

// CanSeeEntity returns true if the AI can see the given entity.
// TODO: Instead we should just take a position and check if we have line of
// sight to that position, plus some modifier based on the entity's size or
// stealth.
func (p *Perception) CanSeeEntity(e Entity) bool {
	for _, e2 := range p.Entities {
		if e.ID() == e2.ID() {
			return true
		}
	}
	return false
}
