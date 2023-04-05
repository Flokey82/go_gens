package gencitymap

import (
	"math"
	"sort"

	"github.com/Flokey82/go_gens/utils"
)

// NOTE: This code is based on the PolyK library:
// https://github.com/martywallace/polyk

type Point struct {
	x    float64
	y    float64
	flag bool
}

/**
 * Checks, if polygon is simple. Polygon is simple, when its edges don't cross each other.
 *
 * @param {number[]} polygon [x1, y1, x2, y2...]
 * @returns {boolean} true if Polygon is simple
 */
func IsSimple(polygon []float64) bool {
	p := polygon
	n := len(p) >> 1
	if n < 4 {
		return true
	}
	a1 := Point{}
	a2 := Point{}
	b1 := Point{}
	b2 := Point{}
	c := Point{}

	for i := 0; i < n; i++ {
		a1.x = p[2*i]
		a1.y = p[2*i+1]
		if i == n-1 {
			a2.x = p[0]
			a2.y = p[1]
		} else {
			a2.x = p[2*i+2]
			a2.y = p[2*i+3]
		}

		for j := 0; j < n; j++ {
			if utils.Abs(i-j) < 2 {
				continue
			}
			if j == n-1 && i == 0 {
				continue
			}
			if i == n-1 && j == 0 {
				continue
			}

			b1.x = p[2*j]
			b1.y = p[2*j+1]
			if j == n-1 {
				b2.x = p[0]
				b2.y = p[1]
			} else {
				b2.x = p[2*j+2]
				b2.y = p[2*j+3]
			}

			if _, ok := GetLineIntersection(a1, a2, b1, b2, c); ok {
				return false
			}
		}
	}
	return true
}

/**
 * Checks, if polygon is convex. Polygon is convex, when each inner angle is <= 180°.
 *
 * @param {number[]} polygon [x1, y1, x2, y2...]
 * @returns {boolean}
 */
func IsConvex(polygon []float64) bool {
	p := polygon
	if len(p) < 6 {
		return true
	}
	l := len(p) - 4
	for i := 0; i < l; i += 2 {
		if !convex(p[i], p[i+1], p[i+2], p[i+3], p[i+4], p[i+5]) {
			return false
		}
	}
	if !convex(p[l], p[l+1], p[l+2], p[l+3], p[0], p[1]) {
		return false
	}
	if !convex(p[l+2], p[l+3], p[0], p[1], p[2], p[3]) {
		return false
	}
	return true
}

/**
 * Returns the area of polygon.
 *
 * @param {number[]} polygon [x1, y1, x2, y2...]
 * @returns {number}
 */
func getArea(polygon []float64) float64 {
	if len(polygon) < 6 {
		return 0
	}
	var sum float64 = 0
	l := len(polygon) - 2
	for i := 0; i < l; i += 2 {
		sum += (polygon[i+2] - polygon[i]) * (polygon[i+1] + polygon[i+3])
	}
	sum += (polygon[0] - polygon[l]) * (polygon[l+1] + polygon[1])
	return -sum * 0.5
}

type AABB struct {
	x      float64
	y      float64
	width  float64
	height float64
}

/**
 * Returns the Axis-aligned Bounding Box of polygon
 *
 * @param {number[]} polygon [x1, y1, x2, y2...]
 * @returns {AABB}
 * @example
 * //={x:0, y:0, width:0, height:0}
 */
func GetAABB(polygon []float64) AABB {
	var minx = math.Inf(1)
	var miny = math.Inf(1)
	var maxx = -minx
	var maxy = -miny
	for i := 0; i < len(polygon); i += 2 {
		minx = math.Min(minx, polygon[i])
		maxx = math.Max(maxx, polygon[i])
		miny = math.Min(miny, polygon[i+1])
		maxy = math.Max(maxy, polygon[i+1])
	}
	return AABB{
		x:      minx,
		y:      miny,
		width:  maxx - minx,
		height: maxy - miny,
	}
}

/**
 * Computes the triangulation. Output array is array of triangles (triangle = 3 indices of polygon vertices).
 *
 * Works with simple polygons only.
 *
 * @param {number[]} polygon [x1, y1, x2, y2...]
 * @returns {number[]} array of triangles (triangle = 3 indices of polygon vertices)
 * @example
 * var ids = PolyK.Triangulate([0, 0, 1, 0, 1, 1, 0, 1]);
 * //=[0, 1, 2, 0, 2, 3]
 */
func Triangulate(polygon []float64) []int {
	p := polygon
	n := len(p) >> 1
	if n < 3 {
		return nil
	}
	var tgs, avl []int
	for i := 0; i < n; i++ {
		avl = append(avl, i)
	}

	i := 0
	al := n
	for al > 3 {
		i0 := avl[(i+0)%al]
		i1 := avl[(i+1)%al]
		i2 := avl[(i+2)%al]

		ax := p[2*i0]
		ay := p[2*i0+1]
		bx := p[2*i1]
		by := p[2*i1+1]
		cx := p[2*i2]
		cy := p[2*i2+1]

		earFound := false
		if convex(ax, ay, bx, by, cx, cy) {
			earFound = true
			for j := 0; j < al; j++ {
				vi := avl[j]
				if vi == i0 || vi == i1 || vi == i2 {
					continue
				}
				if PointInTriangle(p[2*vi], p[2*vi+1], ax, ay, bx, by, cx, cy) {
					earFound = false
					break
				}
			}
		}
		if earFound {
			tgs = append(tgs, i0, i1, i2)
			avl = append(avl[:i+1], avl[i+2:]...)
			al--
			i = 0
		} else if i++; i > 3*al {
			break // no convex angles :(
		}
	}
	return append(tgs, avl[0], avl[1], avl[2])
}

func Slice(polygon []float64, startX, startY, endX, endY float64) [][]float64 {
	p := polygon
	ax := startX
	ay := startY
	bx := endX
	by := endY
	// if ContainsPoint(p, ax, ay) || ContainsPoint(p, bx, by) {
	// 	return []float64{p}
	// }

	a := Point{x: ax, y: ay}
	b := Point{x: bx, y: by}
	var iscs []Point // intersections
	var ps []Point   // points
	for i := 0; i < len(p); i += 2 {
		ps = append(ps, Point{x: p[i], y: p[i+1]})
	}
	for i := 0; i < len(ps); i++ {
		isc := Point{}
		iscp, ok := GetLineIntersection(a, b, ps[i], ps[(i+1)%len(ps)], isc)
		var fisc Point
		var lisc Point
		if len(iscs) > 0 {
			fisc = iscs[0]
			lisc = iscs[len(iscs)-1]
		}
		// && (isc.x!=ps[i].x || isc.y!=ps[i].y) )
		if ok && (fisc == (Point{}) || distance(iscp, fisc) > 1e-10) && (lisc == (Point{}) || distance(iscp, lisc) > 1e-10) {
			iscp.flag = true
			iscs = append(iscs, iscp)
			ps = append(ps, Point{})
			copy(ps[i+2:], ps[i+1:])
			ps[i+1] = iscp
			i++
		}
	}

	if len(iscs) < 2 {
		return [][]float64{p}
	}
	var comp = func(u, v Point) bool {
		return distance(a, u)-distance(a, v) < 0
	}
	sort.Slice(iscs, func(i, j int) bool { return comp(iscs[i], iscs[j]) })

	var pgs [][]Point
	var dir int
	for len(iscs) > 1 {
		// n := len(ps) // is assigned a value but never used. (no-unused-vars)
		i0 := iscs[0]
		i1 := iscs[1]
		// if(i0.x==i1.x && i0.y==i1.y) { iscs.splice(0,2); continue;}
		index0 := indexOfPoint(ps, i0)
		index1 := indexOfPoint(ps, i1)
		solved := false

		if firstWithFlag(ps, index0) == index1 {
			solved = true
		} else {
			i0 = iscs[1]
			i1 = iscs[0]
			index0 = indexOfPoint(ps, i0)
			index1 = indexOfPoint(ps, i1)
			if firstWithFlag(ps, index0) == index1 {
				solved = true
			}
		}
		if solved {
			dir--
			pgn := getPoints(ps, index0, index1)
			pgs = append(pgs, pgn)
			ps = getPoints(ps, index1, index0)
			i0.flag = false
			i1.flag = false
			iscs = iscs[2:]
			if len(iscs) == 0 {
				pgs = append(pgs, ps)
			}
		} else {
			dir++
			iscs = reversePoints(iscs)

		}
		if dir > 1 {
			break
		}
	}
	result := make([][]float64, 0)
	for i := 0; i < len(pgs); i++ {
		pg := pgs[i]
		npg := make([]float64, 0)
		for j := 0; j < len(pg); j++ {
			npg = append(npg, pg[j].x, pg[j].y)
		}
		result = append(result, npg)
	}
	return result
}

func reversePoints(points []Point) []Point {
	n := len(points)
	for i := 0; i < n/2; i++ {
		points[i], points[n-i-1] = points[n-i-1], points[i]
	}
	return points
}

func indexOfPoint(points []Point, p Point) int {
	for i := 0; i < len(points); i++ {
		if points[i] == p {
			return i
		}
	}
	return -1
}

/**
 * Checks, if polygon contains [x, y].
 *
 * Works with simple polygons only.
 *
 * @param {number[]} polygon [x1, y1, x2, y2...]
 * @param {number} pointX Coordinate [x]
 * @param {number} pointY Coordinate [y]
 * @returns {boolean} depth
 */
func ContainsPoint(polygon []float64, pointX, pointY float64) bool {
	p := polygon
	px := pointX
	py := pointY
	n := len(p) / 2
	ax := 0.0
	ay := p[2*n-3] - py
	bx := p[2*n-2] - px
	by := p[2*n-1] - py
	lup := by > ay
	for i := 0; i < n; i++ {
		ax = bx
		ay = by
		bx = p[2*i] - px
		by = p[2*i+1] - py
		if ay == by {
			continue
		}
		lup = by > ay
	}
	depth := 0
	for i := 0; i < n; i++ {
		ax = bx
		ay = by
		bx = p[2*i] - px
		by = p[2*i+1] - py
		if ay < 0 && by < 0 {
			continue // both "up" or both "down"
		}
		if ay > 0 && by > 0 {
			continue // both "up" or both "down"
		}
		if ax < 0 && bx < 0 {
			continue // both points on the left
		}
		if ay == by && math.Min(ax, bx) <= 0 {
			return true
		}
		if ay == by {
			continue
		}
		lx := ax + (bx-ax)*(-ay)/(by-ay)
		if lx == 0 {
			return true // point on edge
		}
		if lx > 0 {
			depth++
		}
		if ay == 0 && lup && by > ay {
			depth-- // hit vertex, both up
		}
		if ay == 0 && !lup && by < ay {
			depth-- // hit vertex, both down
		}
		lup = by > ay
	}
	return (depth & 1) == 1
}

/**
 * Finds the closest point of polygon, which lays on ray defined by [x,y] (origin) and [dx,dy] (direction).
 *
 * "dist" is the distance of the polygon point, "edge" is the number of the edge, on which intersection occurs, "norm" is the normal in that place, "refl" is reflected direction.
 *
 * Works with simple polygons only.
 *
 * @param {number[]} polygon [x1, y1, x2, y2...]
 * @param {number} originX Origin [x]
 * @param {number} originY Origin [y]
 * @param {number} directionX Direction [x]
 * @param {number} directionY Direction [y]
 * @returns {Raycast}
 * @example
 * //={dist:0, edge:0, norm:{x:0, y:0}, refl:{x:0, y:0}}
 */
/*
func Raycast(polygon []float64, originX, originY, directionX, directionY float64) (isc Raycast) {
	p := polygon
	x := originX
	y := originY
	dx := directionX
	dy := directionY
	l := len(p) - 2
	empty := emptyPoints()
	a1 := empty[0]
	a2 := empty[1]
	b1 := empty[2]
	b2 := empty[3]
	c := empty[4]
	a1.x = x
	a1.y = y
	a2.x = x + dx
	a2.y = y + dy

	if isc == nil {
		isc = Raycast{dist: 0, edge: 0, norm: Point{x: 0, y: 0}, refl: Point{x: 0, y: 0}}
	}
	isc.dist = math.Inf(1)

	var nisc Raycast
	for i := 0; i < l; i += 2 {
		b1.x = p[i]
		b1.y = p[i+1]
		b2.x = p[i+2]
		b2.y = p[i+3]
		nisc = RayLineIntersection(a1, a2, b1, b2, c)
		if nisc != nil {
			isc = updateISC(dx, dy, a1, b1, b2, c, i/2, isc)
		}
	}
	b1.x = b2.x
	b1.y = b2.y
	b2.x = p[0]
	b2.y = p[1]
	nisc = RayLineIntersection(a1, a2, b1, b2, c)
	if nisc != nil {
		isc = updateISC(dx, dy, a1, b1, b2, c, (len(p)/2)-1, isc)
	}

	if isc.dist != math.Inf(1) {
		return isc
	}
	return nil
}*/

/**
 * Finds the point on polygon edges, which is closest to [x,y]. Returns an object in this format
 *
 * "dist" is the distance of the polygon point, "edge" is the number of the closest edge, "point" is the closest point on that edge, "norm" is the normal from "point" to [x,y].
 *
 * @param {number[]} polygon [x1, y1, x2, y2...]
 * @param {number} x Coordinate [x]
 * @param {number} y Coordinate [y]
 * @returns {ClosestEdge}
 * @example
 * //={dist:0, edge:0, point:{x:0, y:0}, norm:{x:0, y:0}}
 */
/*
func ClosestEdge(polygon []float64, x, y float64) (isc ClosestEdge) {
	var p = polygon
	var l = len(p) - 2
	var empty = emptyPoints()
	var a1 = empty[0]
	var b1 = empty[2]
	var b2 = empty[3]
	// var c = tp[4] // is assigned a value but never used.
	a1.x = x
	a1.y = y

	isc.dist = math.Inf(1)

	for i := 0; i < l; i += 2 {
		b1.x = p[i]
		b1.y = p[i+1]
		b2.x = p[i+2]
		b2.y = p[i+3]
		isc = pointLineDist(a1, b1, b2, i>>1, isc)
	}
	b1.x = b2.x
	b1.y = b2.y
	b2.x = p[0]
	b2.y = p[1]
	isc = pointLineDist(a1, b1, b2, l>>1, isc)

	idst := 1 / isc.dist
	isc.norm.x = (x - isc.point.x) * idst
	isc.norm.y = (y - isc.point.y) * idst
	return isc
}
*/
/**
 * Reverse
 *
 * @param {number[]} polygon [x1, y1, x2, y2...]
 */
func Reverse(polygon []float64) []float64 {
	p := polygon
	np := []float64{}
	for j := len(p) - 2; j >= 0; j -= 2 {
		np = append(np, p[j], p[j+1])
	}
	return np
}

type ISC struct {
	dist  float64
	Edge  *Edge
	Point *Point
}

type Edge struct {
	Index int
	A     *Point
	B     *Point
}

/*
func pointLineDist(p *Point, a *Point, b *Point, edge *Edge, isc *ISC) *ISC {
	x := p.X
	y := p.Y
	x1 := a.X
	y1 := a.Y
	x2 := b.X
	y2 := b.Y

	A := x - x1
	B := y - y1
	C := x2 - x1
	D := y2 - y1

	dot := A*C + B*D
	lenSq := C*C + D*D
	param := dot / lenSq

	var xx float64
	var yy float64

	if param < 0 || (x1 == x2 && y1 == y2) {
		xx = x1
		yy = y1
	} else if param > 1 {
		xx = x2
		yy = y2
	} else {
		xx = x1 + param*C
		yy = y1 + param*D
	}

	dx := x - xx
	dy := y - yy
	dst := math.Sqrt(dx*dx + dy*dy)
	if dst < isc.Dist {
		isc.Dist = dst
		isc.Edge = edge
		isc.Point.X = xx
		isc.Point.Y = yy
	}
	return isc
}

func updateISC(dx, dy float64, a1, b1, b2, c, edge, isc Point) Point {
	var nrl float64 = distance(a1, c)
	if nrl < isc.dist {
		var ibl float64 = 1 / distance(b1, b2)
		var nx float64 = -(b2.y - b1.y) * ibl
		var ny float64 = (b2.x - b1.x) * ibl
		var ddot float64 = 2 * (dx*nx + dy*ny)
		isc.dist = nrl
		isc.norm.x = nx
		isc.norm.y = ny
		isc.refl.x = -ddot*nx + dx
		isc.refl.y = -ddot*ny + dy
		isc.edge = edge
	}
	return isc
}*/

/**
 * Get Points
 *
 * @private
 * @param {number[]} points
 * @param {number} index0
 * @param {number} index1
 * @returns {number[]} points
 */
func getPoints(points []Point, index0 int, index1 int) []Point {
	var n int = len(points)
	var result []Point
	if index1 < index0 {
		index1 += n
	}
	for i := index0; i <= index1; i++ {
		result = append(result, points[i%n])
	}
	return result
}

/**
 * First With Flag
 *
 * @private
 * @param {Point[]} points
 * @param {number} index
 * @returns {number}
 */
func firstWithFlag(points []Point, index int) int {
	n := len(points)
	for {
		index = (index + 1) % n
		if points[index].flag {
			return index
		}
	}
	return -1
}

/**
 * Point in Triangle
 *
 * @private
 * @param {number} px
 * @param {number} py
 * @param {number} ax
 * @param {number} ay
 * @param {number} bx
 * @param {number} by
 * @param {number} cx
 * @param {number} cy
 * @returns {boolean}
 */
func PointInTriangle(px, py, ax, ay, bx, by, cx, cy float64) bool {
	v0x := cx - ax
	v0y := cy - ay
	v1x := bx - ax
	v1y := by - ay
	v2x := px - ax
	v2y := py - ay

	dot00 := v0x*v0x + v0y*v0y
	dot01 := v0x*v1x + v0y*v1y
	dot02 := v0x*v2x + v0y*v2y
	dot11 := v1x*v1x + v1y*v1y
	dot12 := v1x*v2x + v1y*v2y

	invDenom := 1 / (dot00*dot11 - dot01*dot01)
	u := (dot11*dot02 - dot01*dot12) * invDenom
	v := (dot00*dot12 - dot01*dot02) * invDenom

	// Check if point is in triangle
	return (u >= 0) && (v >= 0) && (u+v < 1)
}

/**
 * RayLine Intersection
 *
 * @private
 * @param {Point} a1
 * @param {Point} a2
 * @param {Point} b1
 * @param {Point} b2
 * @param {Point} c
 */
/*
func rayLineIntersection(a1, a2, b1, b2, c *Point) *Point {
	dax := (a1.x - a2.x)
	dbx := (b1.x - b2.x)
	day := (a1.y - a2.y)
	dby := (b1.y - b2.y)

	Den := dax*dby - day*dbx
	if Den == 0 {
		return nil // parallel
	}

	A := (a1.x*a2.y - a1.y*a2.x)
	B := (b1.x*b2.y - b1.y*b2.x)

	I := c
	iDen := 1 / Den
	I.x = (A*dbx - dax*B) * iDen
	I.y = (A*dby - day*B) * iDen

	if !InRectangle(I, b1, b2) {
		return nil
	}
	if (day > 0 && I.y > a1.y) || (day < 0 && I.y < a1.y) {
		return nil
	}
	if (dax > 0 && I.x > a1.x) || (dax < 0 && I.x < a1.x) {
		return nil
	}
	return I
}
*/
/* RayLineIntersection
 *
 * @private
 * @param {Point} a1
 * @param {Point} a2
 * @param {Point} b1
 * @param {Point} b2
 * @param {Point} c
 */
/*
func RayLineIntersection(a1, a2, b1, b2, c *Point) *Point {
	dax := (a1.x - a2.x)
	dbx := (b1.x - b2.x)
	day := (a1.y - a2.y)
	dby := (b1.y - b2.y)

	Den := dax*dby - day*dbx
	if Den == 0 {
		return nil // parallel
	}

	A := (a1.x*a2.y - a1.y*a2.x)
	B := (b1.x*b2.y - b1.y*b2.x)

	I := c
	iDen := 1 / Den
	I.x = (A*dbx - dax*B) * iDen
	I.y = (A*dby - day*B) * iDen

	if !InRectangle(I, b1, b2) {
		return nil
	}
	if (day > 0 && I.y > a1.y) || (day < 0 && I.y < a1.y) {
		return nil
	}
	if (dax > 0 && I.x > a1.x) || (dax < 0 && I.x < a1.x) {
		return nil
	}
	return I
}*/

// Get Line Intersection
//
// @private
// @param {Point} a1
// @param {Point} a2
// @param {Point} b1
// @param {Point} b2
// @param {Point} c
// @returns {Point}
func GetLineIntersection(a1, a2, b1, b2, c Point) (Point, bool) {
	dax := a1.x - a2.x
	dbx := b1.x - b2.x
	day := a1.y - a2.y
	dby := b1.y - b2.y

	Den := dax*dby - day*dbx

	if Den == 0 {
		return Point{}, false // parallel
	}

	A := a1.x*a2.y - a1.y*a2.x
	B := b1.x*b2.y - b1.y*b2.x

	I := c
	I.x = (A*dbx - dax*B) / Den
	I.y = (A*dby - day*B) / Den

	if InRectangle(I, a1, a2) && InRectangle(I, b1, b2) {
		return I, true
	}
	return Point{}, false
}

/**
 * In Rectangle
 *
 * @private
 * @param {Point} a
 * @param {Point} b
 * @param {Point} c
 * @return {boolean}
 */
func InRectangle(a, b, c Point) bool {
	minx := math.Min(b.x, c.x)
	maxx := math.Max(b.x, c.x)
	miny := math.Min(b.y, c.y)
	maxy := math.Max(b.y, c.y)

	if minx == maxx {
		return (miny <= a.y && a.y <= maxy)
	}
	if miny == maxy {
		return (minx <= a.x && a.x <= maxx)
	}

	// return (minx <= a.x && a.x <= maxx && miny <= a.y && a.y <= maxy)
	return (minx <= a.x+1e-10 && a.x-1e-10 <= maxx && miny <= a.y+1e-10 && a.y-1e-10 <= maxy)
}

/**
 * Convex
 *
 * @private
 * @param {Point} ax
 * @param {Point} ay
 * @param {Point} bx
 * @param {Point} by
 * @param {Point} cx
 * @param {Point} cy
 * @returns {boolean}
 */
func convex(ax, ay, bx, by, cx, cy float64) bool {
	return (ay-by)*(cx-bx)+(bx-ax)*(cy-by) >= 0
}

/**
 * Distance
 *
 * @private
 * @param {Point} a
 * @param {Point} b
 * @returns {number}
 */
func distance(a, b Point) float64 {
	var dx = b.x - a.x
	var dy = b.y - a.y
	return math.Sqrt(dx*dx + dy*dy)
}

// Creates an array of empty Points
//
// @private
// @param {number} [num=10] Number of points
// @returns {Point[]}
/*
func emptyPoints(num int) []Point {
	if num == 0 {
		num = 10
	}
	var container []Point
	for i := 0; i < num; i++ {
		container = append(container, Point{0, 0})
	}
	return container
}
*/
