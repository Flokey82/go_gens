package gamecs

import (
	"fmt"
	"log"

	"github.com/Flokey82/go_gens/vectors"
)

// CMovable is a movable component.
type CMovable struct {
	Pos   vectors.Vec2
	Speed vectors.Vec2
}

func newCMovable(pos vectors.Vec2) *CMovable {
	return &CMovable{
		Pos: pos,
	}
}

func (c *CMovable) Update(delta float64) {
	c.Pos.AddToThis(c.Speed)
	log.Println(fmt.Sprintf("moved to %.2f, %.2f", c.Pos.X, c.Pos.Y))
}
