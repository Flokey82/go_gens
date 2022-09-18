package gamecs

import (
	"math"

	"github.com/Flokey82/go_gens/gameconstants"
	"github.com/Flokey82/go_gens/vectors"
)

// CAiPath is a path planning component.
type CAiPath struct {
	ai              *CAi
	Waypoints       []vectors.Vec2 // Current list of waypoints.
	WaypointCurrent int            // Current index in the waypoints array.
	Target          vectors.Vec2   // Our current target.
	active          bool           // We are actively moving towards a target.
	running         bool           // We move at running speed.
	planned         bool           // We have currently waypoints planned.
}

func newCAiPath() *CAiPath {
	return new(CAiPath)
}

func (c *CAiPath) init(ai *CAi) {
	c.ai = ai
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
	// TODO: Replace with A* or similar stuff.
	curWay := m.Pos
	for vectors.Dist2(curWay, c.Target) > 5.0 {
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
	if vectors.Dist2(curWay, m.Pos) >= 0.2 {
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
	dist := vectors.Dist2(m.Pos, c.Target)
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
		magnitude *= gameconstants.RunningSpeed
	} else {
		magnitude *= gameconstants.WalkingSpeed
	}

	// If we would overshoot the target, we limit the speed vector
	// length to the current distance.
	if dist < magnitude {
		magnitude = dist
	}
	m.Speed.MulWithThis(magnitude)
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
