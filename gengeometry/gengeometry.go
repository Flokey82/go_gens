package gengeometry

import (
	"image"
	"image/color"
	"log"
	"sort"

	"github.com/Flokey82/go_gens/vectors"
)

func DrawLine(img *image.RGBA, a, b vectors.Vec2, color color.RGBA, scale float64) {
	// Bresenham's line algorithm
	x0 := int(a.X * scale)
	y0 := int(a.Y * scale)
	x1 := int(b.X * scale)
	y1 := int(b.Y * scale)

	dx := Abs(float64(x1 - x0))
	dy := Abs(float64(y1 - y0))
	sx := 1
	if x0 >= x1 {
		sx = -1
	}
	sy := 1
	if y0 >= y1 {
		sy = -1
	}
	err := dx - dy

	for {
		img.Set(x0, y0, color)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err = err - dy
			x0 = x0 + sx
		}
		if e2 < dx {
			err = err + dx
			y0 = y0 + sy
		}
	}
}

func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Polygon is a polygon with multiple points.
type Polygon struct {
	Points []vectors.Vec2
	// SubAreas []Polygon TODO: Populate when splitting a polygon.
}

// DrawToImage draws the polygon to an image.
func (p Polygon) DrawToImage(img *image.RGBA, color color.RGBA, scale float64) {
	if len(p.Points) == 0 {
		return
	}
	for i := 0; i < len(p.Points)-1; i++ {
		DrawLine(img, p.Points[i], p.Points[i+1], color, scale)
	}
	DrawLine(img, p.Points[len(p.Points)-1], p.Points[0], color, scale)
}

// Split splits a polygon into multiple polygons by cutting it along a line.
// NOTE: This is buggy, especially when an intersection point is on a vertex.
func (p *Polygon) Split(l vectors.Segment) []*Polygon {
	// https://github.com/gilengel/mapme/blob/main/src/algorithm/geo.rs
	// https://github.com/gpicavet/split-poly/blob/master/index.js
	// https://github.com/xidiq/gitinit/blob/main/lizmap/www/OpenLayers-2.13/lib/OpenLayers/Geometry/Polygon.js
	type interPoint struct {
		is        int          // start segment index
		ie        int          // end segment index
		p         vectors.Vec2 // intersection point
		t         float64      // distance from start to intersection point
		crossback *Polygon     // polygon connected through this intersection
	}

	var interPoints []interPoint
	rayOrig := l.Start
	rayDir := vectors.Vec2{
		X: l.End.X - l.Start.X,
		Y: l.End.Y - l.Start.Y,
	}
	start := p.Points[len(p.Points)-1]
	for ivert := 0; ivert < len(p.Points); ivert++ {
		end := p.Points[ivert]

		// li := Line{start, end}
		// _, intPt := l.Intersects(li)
		// inter = start + s * (end-start)
		edgeDir := vectors.Vec2{
			X: end.X - start.X,
			Y: end.Y - start.Y,
		}
		den := rayDir.X*edgeDir.Y - rayDir.Y*edgeDir.X
		if den != 0 {
			num := rayDir.X*(rayOrig.Y-start.Y) - rayDir.Y*(rayOrig.X-start.X)
			// s = (rayOrig-start) dot (-rayDir.y, rayDir.x) / (end-start) dot (-rayDir.y, rayDir.x)
			if s := num / den; s >= 0 && s <= 1 { // intersection in edge if s>=0 && s<=1
				pt := vectors.Vec2{
					X: start.X + s*edgeDir.X,
					Y: start.Y + s*edgeDir.Y,
				}
				// log.Println("Intersection point:", intPt)
				// log.Println("Intersection point2:", pt)
				// inter = rayOrig + t * rayDir
				// t = det((end-start),(rayOrig-start)) / (end-start) dot (-rayDir.y, rayDir.x)
				num = edgeDir.X*(rayOrig.Y-start.Y) - edgeDir.Y*(rayOrig.X-start.X)
				interPoints = append(interPoints, interPoint{
					is: (ivert + len(p.Points) - 1) % len(p.Points),
					ie: ivert,
					p:  pt,
					t:  num / den,
				})
			}
		}
		start = end
	}

	// Sort inter points by distance from the ray origin.
	sort.Slice(interPoints, func(i, j int) bool {
		return interPoints[i].t > interPoints[j].t
	})

	// Make sure there is an even number of intersection points.
	if len(interPoints)%2 != 0 {
		log.Println("Odd number of intersection points")
		return nil
	}

	// List of new polygons with a first empty one (make it current).
	output := []*Polygon{new(Polygon)}
	curPoly := output[0]

	// Walk through initial polygon points.
	for ivert := 0; ivert < len(p.Points); ivert++ {
		// Append first point to polygon.
		curPoly.Points = append(curPoly.Points, p.Points[ivert])

		// Is there an intersection point?
		inter := -1
		for interTmp := 0; interTmp < len(interPoints); interTmp++ {
			if interPoints[interTmp].is == ivert {
				inter = interTmp
				break
			}
		}

		if inter != -1 {
			// Yes, add the inter point to the current polygon.
			curPoly.Points = append(curPoly.Points, interPoints[inter].p)

			// Set the paired inter point to be the crossback point of this polygon.
			if inter%2 == 0 {
				interPoints[inter+1].crossback = curPoly
			} else {
				interPoints[inter-1].crossback = curPoly
			}

			// Now we have to switch the current poly to a pending one or to a new one.
			if interPoints[inter].crossback != nil {
				curPoly = interPoints[inter].crossback
			} else {
				curPoly = new(Polygon)
				output = append(output, curPoly)
			}

			// Add the inter point to the new current polygon.
			curPoly.Points = append(curPoly.Points, interPoints[inter].p)
		}
	}

	// Return the polygons.
	return output
}
