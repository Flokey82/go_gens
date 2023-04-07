package gengeometry

import (
	"github.com/Flokey82/go_gens/vectors"
)

func getSqSegDist(p vectors.Vec2, p1, p2 vectors.Vec2) float64 {
	x := p1.X
	y := p1.Y
	dx := p2.X - x
	dy := p2.Y - y

	if dx != 0 || dy != 0 {
		t := (p.X-x)*dx + (p.Y-y)*dy/(dx*dx+dy*dy)
		if t > 1 {
			x = p2.X
			y = p2.Y
		} else if t > 0 {
			x += dx * t
			y += dy * t
		}
	}
	dx = p.X - x
	dy = p.Y - y
	return dx*dx + dy*dy
}

func simplifyRadialDist(points []vectors.Vec2, sqTolerance float64) []vectors.Vec2 {
	var prevPoint vectors.Vec2 = points[0]
	var newPoints = []vectors.Vec2{prevPoint}
	var point vectors.Vec2
	for i := 1; i < len(points); i++ {
		point = points[i]
		if point.DistanceToSquared(prevPoint) > sqTolerance {
			newPoints = append(newPoints, point)
			prevPoint = point
		}
	}
	if prevPoint != point {
		newPoints = append(newPoints, point)
	}
	return newPoints
}

func simplifyDPStep(points []vectors.Vec2, first, last int, sqTolerance float64, simplified *[]vectors.Vec2) {
	var maxSqDist float64 = sqTolerance
	var index int
	for i := first + 1; i < last; i++ {
		sqDist := getSqSegDist(points[i], points[first], points[last])
		if sqDist > maxSqDist {
			index = i
			maxSqDist = sqDist
		}
	}
	if maxSqDist > sqTolerance {
		if index-first > 1 {
			simplifyDPStep(points, first, index, sqTolerance, simplified)
		}
		*simplified = append(*simplified, points[index])
		if last-index > 1 {
			simplifyDPStep(points, index, last, sqTolerance, simplified)
		}
	}
}

func simplifyDouglasPeucker(points []vectors.Vec2, sqTolerance float64) []vectors.Vec2 {
	last := len(points) - 1
	var simplified []vectors.Vec2 = []vectors.Vec2{points[0]}
	simplifyDPStep(points, 0, last, sqTolerance, &simplified)
	simplified = append(simplified, points[last])
	return simplified
}

// SimplifyPolyline simplifies a polyline.
// NOTE: This code is based on:
// https://github.com/mourner/simplify-js
func SimplifyPolyline(points []vectors.Vec2, tolerance float64, highestQuality bool) []vectors.Vec2 {
	if len(points) <= 2 {
		return points
	}
	var sqTolerance float64
	if tolerance != 0 {
		sqTolerance = tolerance * tolerance
	}
	if !highestQuality {
		points = simplifyRadialDist(points, sqTolerance)
	}
	points = simplifyDouglasPeucker(points, sqTolerance)
	return points
}
