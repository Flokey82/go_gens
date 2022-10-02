package gamerogueish

import "math"

// FOV implements a field of view logic.
// TODO: Add raycasting or octants.
// See: http://journal.stuffwithstuff.com/2015/09/07/what-the-hero-sees/
type FOV struct {
	Seen   [][]bool // keeps track of tiles that we have seen and remember
	Radius int      // radius of the FOV
	*World          // world to compute the FOV for
}

// NewFOV returns a new FOV struct.
func NewFOV(w *World, r int) *FOV {
	f := &FOV{
		Radius: r,
		World:  w,
	}
	f.Seen = make([][]bool, f.Height)
	for i := 0; i < f.Height; i++ {
		f.Seen[i] = make([]bool, f.Width)
	}
	return f
}

// Update clears the FOV and recomputes it for the given position.
func (f *FOV) Update(x, y int) {
	f.Clear()
	f.Compute(x, y)
}

// Clear resets all tiles to unseen.
func (f *FOV) Clear() {
	for i := 0; i < f.Height; i++ {
		for j := 0; j < f.Width; j++ {
			f.Seen[i][j] = false
		}
	}
}

// Compute sets all tiles within the FOV radius to seen.
func (f *FOV) Compute(x, y int) {
	for i := 0; i < f.Height; i++ {
		for j := 0; j < f.Width; j++ {
			if f.IsInRadius(x, y, j, i) {
				f.Seen[i][j] = true
			}
		}
	}
}

// IsInRadius returns true if the given coordinates are within the FOV radius.
func (f *FOV) IsInRadius(x, y, x2, y2 int) bool {
	return math.Sqrt(math.Pow(float64(x2-x), 2)+math.Pow(float64(y2-y), 2)) <= float64(f.Radius)
}

/*
func (f *FOV) CanSee(x, y, tgtX, tgtY int) bool {
	if !f.InBounds(x, y) || !f.InBounds(tgtX, tgtY) {
		return false
	}
	if !f.IsInRadius(x, y, tgtX, tgtY) {
		return false
	}
	if !f.IsInLineOfSight(x, y, tgtX, tgtY) {
		return false
	}
	return true
}

func (f *FOV) IsInLineOfSight(x, y, tgtX, tgtY int) bool {
	if x == tgtX && y == tgtY {
		return true
	}
	if x == tgtX {
		return f.IsInLineOfSightVertical(x, y, tgtY)
	}
	if y == tgtY {
		return f.IsInLineOfSightHorizontal(x, y, tgtX)
	}
	return f.IsInLineOfSightDiagonal(x, y, tgtX, tgtY)
}

func (f *FOV) IsInLineOfSightVertical(x, y, tgtY int) bool {
	if y == tgtY {
		return true
	}
	if y > tgtY {
		y, tgtY = tgtY, y
	}
	for i := y + 1; i < tgtY; i++ {
		if f.IsSolid(x, i) {
			return false
		}
	}
	return true
}

func (f *FOV) IsInLineOfSightHorizontal(x, y, tgtX int) bool {
	if x == tgtX {
		return true
	}
	if x > tgtX {
		x, tgtX = tgtX, x
	}
	for i := x + 1; i < tgtX; i++ {
		if f.IsSolid(i, y) {
			return false
		}
	}
	return true
}

func (f *FOV) IsInLineOfSightDiagonal(x, y, tgtX, tgtY int) bool {
	if x == tgtX || y == tgtY {
		return true
	}
	if x > tgtX {
		x, tgtX = tgtX, x
		y, tgtY = tgtY, y
	}
	if y > tgtY {
		y, tgtY = tgtY, y
	}
	// y = mx + b
	m := float64(tgtY-y) / float64(tgtX-x)
	b := float64(y) - m*float64(x)
	for i := x + 1; i < tgtX; i++ {
		j := int(math.Round(m*float64(i) + b))
		if f.IsSolid(i, j) {
			return false
		}
	}
	return true
}*/
