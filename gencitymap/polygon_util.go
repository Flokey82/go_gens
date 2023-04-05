package gencitymap

import (
	"math"
	"math/rand"

	"github.com/Flokey82/go_gens/gengeometry"
	"github.com/Flokey82/go_gens/vectors"
)

type PolygonUtil struct {
}

// SliceRectangle slices a rectangle by line, returning the smallest polygon.
func (p *PolygonUtil) SliceRectangle(origin, worldDimensions, p1, p2 vectors.Vec2) []vectors.Vec2 {
	rectangle := []float64{
		origin.X, origin.Y,
		origin.X + worldDimensions.X, origin.Y,
		origin.X + worldDimensions.X, origin.Y + worldDimensions.Y,
		origin.X, origin.Y + worldDimensions.Y,
	}
	sliced := Slice(rectangle, p1.X, p1.Y, p2.X, p2.Y)
	minArea := CalcPolygonArea(convPairsToVec2s(sliced[0]))
	if len(sliced) > 1 && CalcPolygonArea(convPairsToVec2s(sliced[1])) < minArea {
		return convPairsToVec2s(sliced[1])
	}
	return convPairsToVec2s(sliced[0])
}

func convPairsToVec2s(pairs []float64) []vectors.Vec2 {
	var res []vectors.Vec2
	for i := 0; i < len(pairs); i += 2 {
		res = append(res, vectors.NewVec2(pairs[i], pairs[i+1]))
	}
	return res
}

/**
 * Used to create sea polygon
 */
/*
func (p *PolygonUtil) LineRectanglePolygonIntersection(origin, worldDimensions, line vectors.Vec2) []vectors.Vec2 {
	jstsLine := p.lineToJts(line)
	bounds := []vectors.Vec2{
		origin,
		vectors.NewVec2(origin.X+worldDimensions.X, origin.Y),
		vectors.Vec2(origin.X+worldDimensions.X, origin.Y+worldDimensions.Y),
		vectors.Vec2(origin.X, origin.Y+worldDimensions.Y),
	}
	boundingPoly := p.polygonToJts(bounds)
	union := boundingPoly.GetExteriorRing().Union(jstsLine)
	var polygonizer geom.Polygon
	polygonizer.Union(union)
	polygons := polygonizer.GetPolygons()

	smallestArea := math.Inf(1)
	smallestPoly := nil
	for i := polygons.Iterator(); i.HasNext(); {
		polygon := i.Next()
		area := polygon.GetArea()
		if area < smallestArea {
			smallestArea = area
			smallestPoly = polygon
		}
	}

	if smallestPoly == nil {
		return []vectors.Vec2{}
	}
	var res []vectors.Vec2
	for _, c := range smallestPoly.GetCoordinates() {
		res = append(res, vectors.NewVec2(c.X, c.Y))
	}
	return res
}*/

func CalcPolygonArea(polygon []vectors.Vec2) float64 {
	total := 0.0
	for i := 0; i < len(polygon); i++ {
		addX := polygon[i].X
		addY := polygon[(i+1)%len(polygon)].Y
		subX := polygon[(i+1)%len(polygon)].X
		subY := polygon[i].Y

		total += (addX * addY * 0.5)
		total -= (subX * subY * 0.5)
	}
	return math.Abs(total)
}

// SubdividePolygon divides a polygon into smaller polygons until the minArea is reached.
func SubdividePolygon(poly []vectors.Vec2, minArea float64) [][]vectors.Vec2 {
	area := CalcPolygonArea(poly)
	if area < 0.5*minArea {
		return nil
	}

	var divided [][]vectors.Vec2 // Array of polygons

	longestSideLength := 0.0
	longestSide := []vectors.Vec2{poly[0], poly[1]}
	perimeter := 0.0

	for i := 0; i < len(poly); i++ {
		sideLength := poly[i].Sub(poly[(i+1)%len(poly)]).Len()
		perimeter += sideLength
		if sideLength > longestSideLength {
			longestSideLength = sideLength
			longestSide = []vectors.Vec2{poly[i], poly[(i+1)%len(poly)]}
		}
	}

	// Shape index
	// Using rectangle ratio of 1:4 as limit
	// if area/perimeter*perimeter < 0.04 {
	if area/(perimeter*perimeter) < 0.04 {
		return nil
	}

	if area < 2*minArea {
		return [][]vectors.Vec2{poly}
	}

	// Between 0.4 and 0.6
	deviation := (rand.Float64() * 0.2) + 0.4

	averagePoint := longestSide[0].Add(longestSide[1]).Mul(deviation)
	differenceVector := longestSide[0].Sub(longestSide[1])
	perpVector := vectors.Normalize(vectors.Vec2{X: differenceVector.Y, Y: -1 * differenceVector.X}).Mul(100)

	bisect := []vectors.Vec2{averagePoint.Add(perpVector), averagePoint.Sub(perpVector)}

	// Array of polygons
	sliced := Slice(PolygonToPolygonArray(poly), bisect[0].X, bisect[0].Y, bisect[1].X, bisect[1].Y)

	// Recursive call
	for _, s := range sliced {
		divided = append(divided, SubdividePolygon(PolygonArrayToPolygon(s), minArea)...)
	}

	return divided
}

// ResizeGeometry resizes a polygon to a given spacing.
func ResizeGeometry(geometry []vectors.Vec2, spacing float64, isPolygon bool) []vectors.Vec2 {
	return gengeometry.StraightSkeleton(geometry, 0.9, spacing)
}

// AveragePoint returns the average point of a polygon.
func AveragePoint(polygon []vectors.Vec2) vectors.Vec2 {
	if len(polygon) == 0 {
		return vectors.Vec2{}
	}

	sum := vectors.Vec2{}
	for _, v := range polygon {
		sum.Add(v)
	}

	return sum.Mul(1 / float64(len(polygon)))
}

func (p *PolygonUtil) InsidePolygon(point vectors.Vec2, polygon []vectors.Vec2) bool {
	// ray-casting algorithm based on
	// http://www.ecse.rpi.edu/Homepages/wrf/Research/Short_Notes/pnpoly.html
	if len(polygon) == 0 {
		return false
	}

	inside := false
	for i, j := 0, len(polygon)-1; i < len(polygon); j = i {
		xi, yi := polygon[i].X, polygon[i].Y
		xj, yj := polygon[j].X, polygon[j].Y

		intersect := (yi > point.Y) != (yj > point.Y) && point.X < (xj-xi)*(point.Y-yi)/(yj-yi)+xi
		if intersect {
			inside = !inside
		}
	}

	return inside
}

func (p *PolygonUtil) PointInRectangle(point vectors.Vec2, origin vectors.Vec2, dimensions vectors.Vec2) bool {
	return point.X >= origin.X && point.Y >= origin.Y && point.X <= dimensions.X && point.Y <= dimensions.Y
}

/**
 * [ v.x, v.y, v.x, v.y ]...
 */
func PolygonToPolygonArray(polygon []vectors.Vec2) []float64 {
	outP := []float64{}
	for _, v := range polygon {
		outP = append(outP, v.X, v.Y)
	}
	return outP
}

/**
 * [ v.x, v.y, v.x, v.y ]...
 */
func PolygonArrayToPolygon(polygon []float64) []vectors.Vec2 {
	outP := []vectors.Vec2{}
	for i := 0; i < len(polygon)/2; i++ {
		outP = append(outP, vectors.Vec2{X: polygon[2*i], Y: polygon[2*i+1]})
	}
	return outP
}
