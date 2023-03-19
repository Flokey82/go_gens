package gengeometry

import (
	"github.com/Flokey82/go_gens/vectors"
)

// ShrinkPath shrinks a polygon by a given factor around its center.
func ShrinkPath(points []vectors.Vec2, shrink float64) []vectors.Vec2 {
	// Calculate the center of the polygon and shrink the polygon
	// around the center.
	if len(points) < 3 {
		return points
	}
	center := CenterOfPath(points)
	newPoints := make([]vectors.Vec2, len(points))
	for i, point := range points {
		newPoints[i] = vectors.Vec2{
			X: center.X + (point.X-center.X)*shrink,
			Y: center.Y + (point.Y-center.Y)*shrink,
		}
	}
	return newPoints
}

// CenterOfPath calculates the center of a polygon (by averaging all points).
func CenterOfPath(points []vectors.Vec2) vectors.Vec2 {
	center := vectors.Vec2{}
	for _, point := range points {
		center.X += point.X
		center.Y += point.Y
	}
	center.X /= float64(len(points))
	center.Y /= float64(len(points))
	return center
}

// TranslatePath translates a path by the given offset.
func TranslatePath(path []vectors.Vec2, offset vectors.Vec2) []vectors.Vec2 {
	newPath := make([]vectors.Vec2, len(path))
	for i, v := range path {
		newPath[i] = v.Add(offset)
	}
	return newPath
}
