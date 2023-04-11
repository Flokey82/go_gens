package simnpcs2

import (
	"math/rand"

	"github.com/Flokey82/go_gens/vectors"
)

// Being represents a being in the world.
type Being struct {
	id            int64
	*CompMoveable // Position and speed.
	*CompStats    // Stats.
	*EventListener
	World *World
}

func NewBeing(world *World) *Being {
	// Find a location that is not an obstacle.
	pos := world.findValidPos()

	// TODO: Return error if no location found.
	return &Being{
		id:            int64(len(world.Beings)),
		CompMoveable:  newCompMoveable(pos),
		CompStats:     newCompStats(100),
		EventListener: newEventListener(),
		World:         world,
	}
}

// AddBeing adds a random being to the world.
func (w *World) AddBeing() {
	w.Beings = append(w.Beings, NewAI(NewBeing(w)))
}

// ID returns the ID of the being.
func (b *Being) ID() int64 {
	return b.id
}

// Pos returns the position of the being.
func (b *Being) Pos() vectors.Vec2 {
	return b.CompMoveable.Pos
}

// Type returns the type of the being.
func (b *Being) Type() EntityType {
	return EntityTypeBeing
}

// Update updates the state of the Being.
func (b *Being) Update(delta float64) {
	// Update Stats.
	b.CompStats.Update(delta)

	// Execute steering behaviors.
	b.CompMoveable.Update(delta)

	// Clear events.
	b.EventListener.Update(delta)
}

// Dead returns true if the being is dead.
func (b *Being) Dead() bool {
	return b.CompStats.Health <= 0
}

// TakeDamage reduces the health of the being.
// TODO: Find a better way to do this. Maybe via an event system?
func (b *Being) TakeDamage(damage float64, attacker Entity) {
	b.CompStats.Health -= damage
}

// InMeleeRange returns true if the being is in melee range of the target.
func (b *Being) InMeleeRange(target Entity) bool {
	return b.Pos().DistanceTo(target.Pos()) < 1.0
}

// randFloat returns a random float between 0 and max.
func randFloat(max float64) float64 {
	return rand.Float64() * max
}
