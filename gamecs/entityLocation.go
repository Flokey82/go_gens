package gamecs

import (
	"github.com/Flokey82/go_gens/vectors"
)

// Location represents a destination with storage space.
// (e.g. building, hut, hiding place, ...)
type Location struct {
	id             int
	w              *World
	Pos            vectors.Vec2 // Position on the map
	*CompInventory              // Location storage.
}

// newLocation returns a new location with the given ID and position.
func newLocation(w *World, id int, pos vectors.Vec2) *Location {
	return &Location{
		w:             w,
		id:            id,
		Pos:           pos,
		CompInventory: newCompInventory(w, id, 100),
	}
}

// ID returns the unique identifier for this location.
func (loc *Location) ID() int {
	return loc.id
}
