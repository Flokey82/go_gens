// Package gencitymap generates a city map.
package gencitymap

import (
	"container/heap"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"

	"github.com/Flokey82/go_gens/vectors"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/mazznoer/colorgrad"
)

// https://glumbosch.home.blog/2020/01/12/layouts-of-a-village-in-the-middle-ages/
// https://www.albert.io/learn/question/ap-human-geography/nucleated-village-layouts

// Map is a map.
type Map struct {
	rng         *rand.Rand
	queue       PriorityQueue
	allSegments []*Segment // until we add a quadTree, we need to keep track of all segments
	cfg         *MapConfig
}

// NewMap creates a new map.
func NewMap(seed int64, cfg *MapConfig) *Map {
	return &Map{
		rng: rand.New(rand.NewSource(seed)),
		cfg: cfg,
	}
}

// Dimensions returns the dimensions of the map.
func (m *Map) Dimensions() (float64, float64) {
	minX, minY, maxX, maxY := m.GetExtent()

	return maxX - minX, maxY - minY
}

func (m *Map) GetExtent() (minX, minY, maxX, maxY float64) {
	for _, seg := range m.allSegments {
		if seg.Point.X < minX {
			minX = seg.Point.X
		}
		if seg.Point.Y < minY {
			minY = seg.Point.Y
		}
		if seg.Point.X > maxX {
			maxX = seg.Point.X
		}
		if seg.Point.Y > maxY {
			maxY = seg.Point.Y
		}
	}

	return
}

// Origin returns the origin of the map.
func (m *Map) Origin() (float64, float64) {
	minX, minY, _, _ := m.GetExtent()

	return minX, minY
}

// ExportToPNG exports the map to a PNG file.
func (m *Map) ExportToPNG(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// calculate bounds
	minX, minY, maxX, maxY := m.GetExtent()

	// create image
	img := image.NewRGBA(image.Rect(0, 0, int(maxX-minX), int(maxY-minY)))

	// Fill the background with black.
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	// Create a new graphic context
	gc := draw2dimg.NewGraphicContext(img)

	lines := m.Streamlines()
	mgr := NewGraph(lines, 1, false)

	find := NewPolygonFinder(mgr.Nodes, PolygonParams{
		MaxLength:      10,
		MinArea:        50,
		ShrinkSpacing:  5,
		ChanceNoDivide: 0.01,
	}, nil)
	//find.findPolygons()
	find.Shrink(false)
	find.Divide(false)

	originX, originY := m.Origin()
	drawPolygon := func(p []vectors.Vec2, colFill, colStroke color.RGBA) {
		gc.BeginPath()
		gc.SetFillColor(colFill)
		gc.SetStrokeColor(colStroke)
		gc.SetLineWidth(1)
		for i, v := range p {
			if i == 0 {
				gc.MoveTo(v.X-originX, v.Y-originY)
			} else {
				gc.LineTo(v.X-originX, v.Y-originY)
			}
		}
		gc.Close()
		gc.FillStroke()
	}

	// Draw the polygons.
	for _, p := range find.Polygons {
		drawPolygon(p, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 0, 122, 255})
	}
	for _, p := range find.ShrunkPolygons {
		drawPolygon(p, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 220, 222, 255})
	}
	for _, p := range find.DividedPolygons {
		drawPolygon(p, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 42, 255})
	}

	// Set some properties
	gc.SetFillColor(color.RGBA{0, 0, 0, 0})
	gc.SetStrokeColor(color.RGBA{0, 0, 0, 255})
	gc.SetLineWidth(2)

	// Get a new color palette.
	grad := colorgrad.Rainbow()
	cols := grad.Colors(uint(len(m.cfg.Rules) + 1))

	// Draw segments.
	for _, seg := range m.allSegments {
		if seg.Prev == nil {
			continue
		}
		colR, colG, colB, colA := cols[seg.Type].RGBA()
		// Convert to 0-255 range.
		colR = colR >> 8
		colG = colG >> 8
		colB = colB >> 8
		colA = colA >> 8
		gc.SetStrokeColor(color.RGBA{uint8(colR), uint8(colG), uint8(colB), uint8(colA)})
		gc.BeginPath()
		gc.MoveTo(seg.Point.X-minX, seg.Point.Y-minY)
		gc.LineTo(seg.Prev.Point.X-minX, seg.Prev.Point.Y-minY)
		gc.Stroke()
		gc.Close()
	}

	// encode image
	err = png.Encode(f, img)
	if err != nil {
		return err
	}

	return nil
}

func (m *Map) Streamlines() [][]vectors.Vec2 {
	var streamlines [][]vectors.Vec2
	for _, seg := range m.allSegments {
		if seg.Prev == nil {
			if seg.Next != nil {
				streamlines = append(streamlines, []vectors.Vec2{{
					X: seg.Point.X,
					Y: seg.Point.Y,
				}, {
					X: seg.Next.Point.X,
					Y: seg.Next.Point.Y,
				}})
			}
		} else {
			streamlines = append(streamlines, []vectors.Vec2{{
				X: seg.Point.X,
				Y: seg.Point.Y,
			}, {
				X: seg.Prev.Point.X,
				Y: seg.Prev.Point.Y,
			}})
		}
	}
	return streamlines
}

// Generate generates the map.
func (m *Map) Generate() {
	// Initialize priority queue.
	m.queue = make(PriorityQueue, 0)
	heap.Init(&m.queue)
	for _, rootSeg := range m.cfg.SeedRoots() {
		heap.Push(&m.queue, rootSeg)
	}
}

// Step performs one iteration of the map generation.
func (m *Map) Step() {
	// Get next segment from queue.
	seg := m.queue.Pop().(*Segment)

	// Extend the segment.
	if seg.Next == nil {
		seg.Next = m.newSegment(seg, false)
		seg.Next.Prev = seg
	}

	// Add a branch by chance.
	cfg := m.cfg.getTypeConfig(seg.Type)
	if seg.Prev != nil && m.rng.Float64() < cfg.BranchingChance {
		m.newSegment(seg, true)
	}
}

func (m *Map) newSegment(origin *Segment, branch bool) *Segment {
	segType := origin.Type
	config := m.cfg.getTypeConfig(segType)
	if branch && segType < RoadType(len(m.cfg.Rules)-1) {
		if !config.BranchSameType || m.rng.Float64() > config.BranchSameTypeChance {
			segType++
			config = m.cfg.getTypeConfig(segType)
		}
	}

	// Calculate length.
	dist := config.LengthMin + config.LengthMin*m.rng.Float64()*config.LengthVariation

	// Calculate angle.
	angle := config.AngleMin + m.rng.Float64()*config.AngleVariance

	// 50% chance to flip angle.
	if config.AngleReversal && m.rng.Float64() < 0.5 {
		angle *= -1
	}

	// If branch, add branching angle.
	if branch {
		originConfig := m.cfg.getTypeConfig(origin.Type)

		// 50% chance to flip angle.
		// TODO: Maybe also sometimes flip 180 degrees?
		if originConfig.BranchingReversal && m.rng.Float64() < 0.5 {
			angle -= originConfig.BranchingAngle
		} else {
			angle += originConfig.BranchingAngle
		}
	}

	// Use direction of previous segment if it exists.
	segVec := origin.GetVector()

	// Calculate new point.
	angle += vectors.Angle2(segVec)

	pNew := origin.Point
	pNew.X += dist * math.Cos(degToRad(angle))
	pNew.Y += dist * math.Sin(degToRad(angle))

	// Create new segment.
	newSeg := &Segment{
		Length: dist,
		Point:  pNew,
		Type:   segType,
		Step:   origin.Step + 1,
		Prev:   origin,
	}

	// Find if any segments intersect with the new segment.
	// Find closest intersection point and set that as the end point.
	// NOTE: This should be maybe optimized by using a quadtree or something.
	var foundIntersect bool
	var currentDist float64
	var ipClosest vectors.Vec2
	var intersectIdx int
	for idx, seg := range m.allSegments {
		if seg == origin || seg == newSeg {
			continue
		}

		// Check if there is an intersection. If so, we set the end point of the new segment to the intersection point.
		// And also we won't add the new segment to the queue.
		if ok, ip := seg.Intersects(newSeg); ok {
			if !newSeg.IsPointOnLine(ip) {
				continue
			}
			// No intersection or point is closer than the current closest.
			newDist := vectors.Dist2(ip, origin.Point)
			if !foundIntersect || newDist < currentDist {
				ipClosest = ip
				currentDist = newDist
				foundIntersect = true
				intersectIdx = idx
			}
		}
	}

	// If we found an intersection, set the end point of the new segment to the intersection point.
	if foundIntersect {
		newSeg.Point = ipClosest
		newSeg.Length = vectors.Dist2(ipClosest, origin.Point)
		newSeg.End = true

		// Split the segment that was intersected.
		seg := m.allSegments[intersectIdx]
		second := seg.Split(ipClosest)
		m.allSegments = append(m.allSegments, second)
	}

	// Add to queue (if not end).
	if !newSeg.End {
		heap.Push(&m.queue, newSeg)
	}

	// Add branch to origin.
	if branch {
		origin.Branches = append(origin.Branches, newSeg)
	}

	// Add to all segments.
	m.allSegments = append(m.allSegments, newSeg)
	return newSeg
}

type PriorityQueue []*Segment

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Step > pq[j].Step
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Segment)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
