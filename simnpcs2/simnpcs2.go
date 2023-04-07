package simnpcs2

import "github.com/Flokey82/go_gens/vectors"

// World represents the game world.
type World struct {
	Beings []Entity
	Items  []Entity
	Width  int
	Height int
}

// NewWorld creates a new world.
func NewWorld(w, h int) *World {
	return &World{
		Width:  w,
		Height: h,
	}
}

// Update updates the state of the world.
func (w *World) Update(delta float64) {
	for _, b := range w.Beings {
		b.Update(delta)
	}
}

// GetEntitiesInRadius returns all entities within a radius of a position.
func (w *World) GetEntitiesInRadius(pos vectors.Vec2, radius float64) []Entity {
	entities := make([]Entity, 0, 10)
	for _, e := range w.Beings {
		if pos.DistanceTo(e.Pos()) < radius {
			entities = append(entities, e)
		}
	}
	return entities
}

// GetItemsInRadius returns all items within a radius of a position.
func (w *World) GetItemsInRadius(pos vectors.Vec2, radius float64) []Entity {
	entities := make([]Entity, 0, 10)
	for _, e := range w.Items {
		if pos.DistanceTo(e.Pos()) < radius {
			entities = append(entities, e)
		}
	}
	return entities
}

type Entity interface {
	ID() int64
	Type() EntityType
	Pos() vectors.Vec2
	Update(delta float64)
}

type EntityType int

const (
	EntityTypeBeing EntityType = iota
	EntityTypeItem
	EntityTypeObstacle
)
