package gamecs

import (
	"fmt"
	"log"

	"github.com/Flokey82/go_gens/vectors"
)

// CompMovable is a movable component.
type CompMovable struct {
	Pos   vectors.Vec2
	Speed vectors.Vec2
}

// newCompMovable returns a new movable component.
func newCompMovable(pos vectors.Vec2) *CompMovable {
	return &CompMovable{
		Pos: pos,
	}
}

// Update moves the position in the component by the speed within the
// given time.
func (c *CompMovable) Update(delta float64) {
	c.Pos.AddToThis(c.Speed)
	log.Println(fmt.Sprintf("moved to %.2f, %.2f", c.Pos.X, c.Pos.Y))
}
