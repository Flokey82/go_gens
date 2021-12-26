package gamecs

import (
	"github.com/Flokey82/go_gens/vectors"
)

type Location struct {
	id          int
	w           *World
	Pos         vectors.Vec2 // Position on the map
	*CInventory              // Location storage.
}

func newLocation(w *World, id int, pos vectors.Vec2) *Location {
	return &Location{
		w:          w,
		id:         id,
		Pos:        pos,
		CInventory: newCInventory(w, id, 100),
	}
}

func (loc *Location) ID() int {
	return loc.id
}
