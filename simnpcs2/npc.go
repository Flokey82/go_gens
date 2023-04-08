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
	World         *World
}

func NewBeing(world *World) *Being {
	// Find a location that is not an obstacle.
	var pos vectors.Vec2
	for i := 0; i < 100; i++ {
		idx := rand.Intn(len(world.Cells))
		if world.CheckIdxReachable(idx) != nil {
			continue
		}
		pos = *world.CellIdxToPos(idx)
		break
	}

	// TODO: Return error if no location found.
	return &Being{
		id:           int64(len(world.Beings)),
		CompMoveable: newCompMoveable(pos),
		CompStats:    newCompStats(),
		World:        world,
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
}

// randFloat returns a random float between 0 and max.
func randFloat(max float64) float64 {
	return rand.Float64() * max
}
