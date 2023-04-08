package simnpcs2

import (
	"fmt"
	"math/rand"

	"github.com/Flokey82/go_gens/vectors"
	"github.com/ojrac/opensimplex-go"
)

// World represents the game world.
type World struct {
	*gifExport
	*webpExport
	Beings []Entity
	Items  []*Item
	Width  int
	Height int
	Cells  []bool
}

// NewWorld creates a new world.
func NewWorld(w, h int, seed int64) *World {
	world := &World{
		gifExport:  newGifExport(),
		webpExport: newWebPExport(w, h),
		Width:      w,
		Height:     h,
		Cells:      make([]bool, w*h),
	}
	world.init(seed)
	return world
}

func (w *World) init(seed int64) {
	// New simplex noise.
	noise := opensimplex.New(seed)

	// Set up some obstacles.
	noisethreshold := -0.55
	for x := 0; x < w.Width; x++ {
		for y := 0; y < w.Height; y++ {
			// Get noise value.
			nval := noise.Eval2(float64(x)/10, float64(y)/10)
			// Set cell to obstacle if noise value is below 0.2.
			if nval < noisethreshold {
				w.Cells[x+y*w.Width] = true
			}
		}
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

func (w *World) findValidPos() vectors.Vec2 {
	var pos vectors.Vec2
	for i := 0; i < 100; i++ {
		idx := rand.Intn(len(w.Cells))
		if w.CheckIdxReachable(idx) != nil {
			continue
		}
		pos = *w.CellIdxToPos(idx)
		break
	}
	return pos
}

// Update updates the state of the world.
func (w *World) Update(delta float64) {
	for _, b := range w.Beings {
		b.Update(delta)
	}
	w.storeWebPFrame()
	w.storeGifFrame()
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
func (w *World) GetItemsInRadius(pos vectors.Vec2, radius float64) []*Item {
	entities := make([]*Item, 0, 10)
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
	Dead() bool
	TakeDamage(damage float64, attacker Entity)
}

type EntityType int

const (
	EntityTypeBeing EntityType = iota
	EntityTypeItem
	EntityTypeObstacle
)
