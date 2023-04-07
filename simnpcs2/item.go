package simnpcs2

import "github.com/Flokey82/go_gens/vectors"

// Item represents a static item in the world.
type Item struct {
	id       int64
	Position vectors.Vec2
}

// AddItem adds a random item to the world.
func (w *World) AddItem() {
	w.Items = append(w.Items, &Item{
		id: int64(len(w.Items)),
		Position: vectors.Vec2{
			X: randFloat(float64(w.Width)),
			Y: randFloat(float64(w.Height)),
		},
	})
}

// ID returns the ID of the item.
func (i *Item) ID() int64 {
	return i.id
}

// Type returns the type of the item.
func (i *Item) Type() EntityType {
	return EntityTypeItem
}

// Pos returns the position of the item.
func (i *Item) Pos() vectors.Vec2 {
	return i.Position
}

// Update updates the state of the item.
func (i *Item) Update(delta float64) {
	// Do nothing.
}
