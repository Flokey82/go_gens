package simnpcs2

import "github.com/Flokey82/go_gens/vectors"

// CompMoveable is a movable component.
type CompMoveable struct {
	Pos   vectors.Vec2
	Speed vectors.Vec2
}

// newCompMoveable returns a new moveable component.
func newCompMoveable(pos vectors.Vec2) *CompMoveable {
	return &CompMoveable{
		Pos: pos,
	}
}

// Update moves the position in the component by the speed within the
// given time.
func (c *CompMoveable) Update(delta float64) {
	c.Pos.AddToThis(c.Speed)
	// TODO: Check for collisions and provide a way to react to them.
}
