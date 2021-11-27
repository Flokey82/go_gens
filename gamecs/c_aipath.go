package gamecs

import (
	"github.com/Flokey82/go_gens/vectors"
	"math"
)

const (
	walkSpeed   = 1.4 // m/s
	runingSpeed = 5.0 // m/s
)

// CAiPath is a path planning component.
type CAiPath struct {
	Target  vectors.Vec2
	active  bool
	running bool
}

func (c *CAiPath) SetTarget(t vectors.Vec2) {
	c.active = true
	c.Target = t
}

func (c *CAiPath) Update(m *CMovable, delta float64) {
	if !c.active {
		return
	}

	dist := calcDist(m.Pos, c.Target)
	if dist < 0.02 {
		c.active = false
		m.Speed = vectors.Vec2{}
		return
	}

	m.Speed = calcNormVec(m.Pos, c.Target)
	mul := delta // magnitude of speed vec.
	if c.running {
		mul *= runingSpeed
	} else {
		mul *= walkSpeed
	}
	if dist < mul {
		mul = dist
	}
	m.Speed.MulWithThis(mul)
}

func calcVec(a, b vectors.Vec2) vectors.Vec2 {
	return vectors.Vec2{
		X: b.X - a.X,
		Y: b.Y - a.Y,
	}
}

func calcNormVec(a, b vectors.Vec2) vectors.Vec2 {
	x := b.X - a.X
	y := b.Y - a.Y
	l := math.Hypot(x, y)
	return vectors.Vec2{
		X: x / l,
		Y: y / l,
	}
}

func calcDist(a, b vectors.Vec2) float64 {
	return math.Hypot(b.X-a.X, b.Y-a.Y)
}
