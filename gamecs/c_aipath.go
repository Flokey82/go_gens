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
	Waypoints       []vectors.Vec2 // Current list of waypoints.
	WaypointCurrent int            // Current index in the waypoints array.
	Target          vectors.Vec2   // Our current target.
	active          bool           // We are actively moving towards a target.
	running         bool           // We move at running speed.
	planned         bool           // We have currently waypoints planned.
}

// resetWaypoints resets all waypoints that we have planned.
func (c *CAiPath) resetWaypoints() {
	c.Waypoints = nil
	c.WaypointCurrent = 0
	c.planned = false
}

// planWaypoints plans a number of waypoints for reaching the current target.
func (c *CAiPath) planWaypoints(m *CMovable) {
	// This implementation is poop.
	curWay := m.Pos
	for calcDist(curWay, c.Target) > 5.0 {
		dir := calcNormVec(curWay, c.Target)
		dir.MulWithThis(5.0)
		curWay.AddToThis(dir)
		c.Waypoints = append(c.Waypoints, curWay)
	}
	c.planned = true
}

// currentWaypoint returns the next position to move towards in order to reach the target.
func (c *CAiPath) currentWaypoint(m *CMovable) vectors.Vec2 {
	if c.WaypointCurrent >= len(c.Waypoints) {
		return c.Target
	}
	curWay := c.Waypoints[c.WaypointCurrent]
	if calcDist(curWay, m.Pos) >= 0.2 {
		return curWay
	}
	c.WaypointCurrent++
	if c.WaypointCurrent >= len(c.Waypoints) {
		return c.Target
	}
	return c.Waypoints[c.WaypointCurrent]
}

// SetTarget resets the current waypoints and sets the new target to move towards.
func (c *CAiPath) SetTarget(t vectors.Vec2) {
	c.resetWaypoints()
	c.active = true
	c.Target = t
}

// Update ticks the AI path planner by delta.
func (c *CAiPath) Update(m *CMovable, delta float64) {
	// Check if we're still actively moving.
	if !c.active {
		return
	}

	// Check if we have already reached the target.
	dist := calcDist(m.Pos, c.Target)
	if dist < 0.02 {
		c.resetWaypoints()
		c.active = false
		m.Speed = vectors.Vec2{}
		return
	}

	// Calculate waypoints.
	if !c.planned {
		c.planWaypoints(m)
	}

	// Calculate speed vector to the current waypoint.
	m.Speed = calcNormVec(m.Pos, c.currentWaypoint(m))

	// Calculate new length of vector based on movement speed,
	// time elapsed, and distance to the target.
	magnitude := delta // Magnitude of speed vector.
	if c.running {
		magnitude *= runingSpeed
	} else {
		magnitude *= walkSpeed
	}

	// If we would overshoot the target, we limit the speed vector
	// length to the current distance.
	if dist < magnitude {
		magnitude = dist
	}
	m.Speed.MulWithThis(magnitude)
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
