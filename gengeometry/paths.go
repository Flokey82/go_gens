package gengeometry

import (
	"math"

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

// StraightSkeleton calculates the straight skeleton of a polygon.
// The straight skeleton is a graph of edges that can be used to generate a roof.
// See: https://github.com/feldhaus/coding-2d-cookbook/blob/main/js/geometry/polygon-straight-skeleton.js
func StraightSkeleton(points []vectors.Vec2, shrink, spacing float64) []vectors.Vec2 {
	order := polygonOrder(points)
	spacing *= order / math.Abs(order)

	resultingPath := make([]vectors.Vec2, len(points))
	l := len(points)
	for i := 0; i < l; i++ {
		p0 := points[i%l]
		p1 := points[(i+1)%l]
		p2 := points[(i+2)%l]

		a0 := vectors.Vec2{X: p1.X - p0.X, Y: p1.Y - p0.Y}
		a1 := vectors.Vec2{X: p2.X - p1.X, Y: p2.Y - p1.Y}

		mi0 := a0.Y / a0.X
		mi1 := a1.Y / a1.X

		li0 := a0.Len()
		li1 := a1.Len()

		ri0 := p0.X + (spacing*a0.Y)/li0
		ri1 := p1.X + (spacing*a1.Y)/li1

		si0 := p0.Y - (spacing*a0.X)/li0
		si1 := p1.Y - (spacing*a1.X)/li1

		point := vectors.Vec2{
			X: (mi1*ri1 - mi0*ri0 + si0 - si1) / (mi1 - mi0),
			Y: (mi0*mi1*(ri1-ri0) + mi1*si0 - mi0*si1) / (mi1 - mi0),
		}

		if a0.X == 0 {
			point.X = p1.X + (spacing*a0.Y)/math.Abs(a0.Y)
			point.Y = mi1*point.X - mi1*ri1 + si1
		}

		if a1.X == 0 {
			point.X = p2.X + (spacing*a1.Y)/math.Abs(a1.Y)
			point.Y = mi0*point.X - mi0*ri0 + si0
		}

		resultingPath[(i+1)%l] = point
	}

	return resultingPath
}

func polygonOrder(points []vectors.Vec2) float64 {
	signedArea := 0.0
	for i := 0; i < len(points); i++ {
		p0 := points[i]
		p1 := points[(i+1)%len(points)]
		signedArea += p0.X*p1.Y - p1.X*p0.Y
	}
	return signedArea
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
