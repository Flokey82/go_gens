package simnpcs2

import (
	"fmt"

	"github.com/Flokey82/go_gens/vectors"
)

// World represents the game world.
type World struct {
	Beings []Entity
	Items  []Entity
	Width  int
	Height int
	Cells  []bool
}

// NewWorld creates a new world.
func NewWorld(w, h int) *World {
	world := &World{
		Width:  w,
		Height: h,
		Cells:  make([]bool, w*h),
	}
	world.init()
	return world
}

func (w *World) init() {
	// Set up some obstacles.
	// TODO: Use noise or something to generate more interesting obstacles.
	for i := 0; i < 100; i++ {
		w.Cells[int(randFloat(float64(w.Width))*float64(w.Height))] = true
	}
}

// PosToCellIdx returns the cell index for a position.
func (w *World) PosToCellIdx(pos *vectors.Vec2) int {
	return int(pos.Y)*w.Width + int(pos.X)
}

// CellIdxToPos returns the position for a cell index.
func (w *World) CellIdxToPos(idx int) *vectors.Vec2 {
	return &vectors.Vec2{
		X: float64(idx % w.Width),
		Y: float64(idx / w.Width),
	}
}

// CheckIdxReachable checks if a cell is reachable.
func (w *World) CheckIdxReachable(idx int) error {
	if idx < 0 || idx >= len(w.Cells) {
		return fmt.Errorf("index %d is out of bounds", idx)
	}
	if w.Cells[idx] {
		return fmt.Errorf("cell %d is not walkable", idx)
	}
	return nil
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
