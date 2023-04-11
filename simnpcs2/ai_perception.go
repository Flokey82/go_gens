package simnpcs2

// Perception represents the perception of an AI of the world.
type Perception struct {
	*AI
	Entities    []Entity // All entities we can see.
	EntitiesNew []Entity // Newly appeared entities.
	Items       []*Item  // All items we can see.
	ItemsNew    []*Item  // Newly appeared items.
}

// newPerception creates a new Perception system for the given AI.
func newPerception(ai *AI) *Perception {
	return &Perception{
		AI:          ai,
		Entities:    make([]Entity, 0, 10),
		EntitiesNew: make([]Entity, 0, 10),
		Items:       make([]*Item, 0, 10),
		ItemsNew:    make([]*Item, 0, 10),
	}
}

// Update updates the perception of the world.
// TODO: Based on insight, we should be able to gain an understanding of the
// motives of other beings. This will allow us to make better decisions, like
// whether to attack or flee from a threat.
func (p *Perception) Update(delta float64) {
	// Keep track of all entities we have already seen.
	seenEntities := make(map[Entity]bool)
	for _, e := range p.Entities {
		seenEntities[e] = true
	}

	// Update all entities we can see.
	p.Entities = p.Entities[:0]
	p.EntitiesNew = p.EntitiesNew[:0]
	for _, e := range p.World.GetEntitiesInRadius(p.Being.Pos(), perceptionRadius) {
		if e.ID() != p.Being.ID() { // Don't perceive ourselves.
			p.Entities = append(p.Entities, e)

			// Newly seen entities.
			if !seenEntities[e] {
				p.EntitiesNew = append(p.EntitiesNew, e)
			}
		}
	}

	// Keep track of all items we have already seen.
	seenItems := make(map[*Item]bool)
	for _, i := range p.Items {
		seenItems[i] = true
	}

	// Update all items we can see.
	p.Items = p.Items[:0]
	p.ItemsNew = p.ItemsNew[:0]
	for _, i := range p.World.GetItemsInRadius(p.Being.Pos(), perceptionRadius) {
		p.Items = append(p.Items, i)

		// Newly seen items.
		if !seenItems[i] {
			p.ItemsNew = append(p.ItemsNew, i)
		}
	}

	// TODO: If we see a threat, make sure that the AI knows about it.
}

// CanSeeEntity returns true if the AI can see the given entity.
// TODO: Instead we should just take a position and check if we have line of
// sight to that position, plus some modifier based on the entity's size or
// stealth.
func (p *Perception) CanSeeEntity(e Entity) bool {
	for _, e2 := range p.Entities {
		if e.Pos().DistanceTo(e2.Pos()) < perceptionRadius {
			return true
		}
	}
	return false
}
