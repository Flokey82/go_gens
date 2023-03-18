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
)

// https://glumbosch.home.blog/2020/01/12/layouts-of-a-village-in-the-middle-ages/
// https://www.albert.io/learn/question/ap-human-geography/nucleated-village-layouts

// Map is a map.
type Map struct {
	Root        *Segment
	rng         *rand.Rand
	queue       PriorityQueue
	allSegments []*Segment // until we add a quadTree, we need to keep track of all segments
}

// NewMap creates a new map.
func NewMap(seed int64) *Map {
	return &Map{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// ExportToPNG exports the map to a PNG file.
func (m *Map) ExportToPNG(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// calculate bounds
	var minX, minY, maxX, maxY float64
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

	// create image
	img := image.NewRGBA(image.Rect(0, 0, int(maxX-minX), int(maxY-minY)))

	// Create a new graphic context
	gc := draw2dimg.NewGraphicContext(img)

	// Set some properties
	gc.SetFillColor(color.RGBA{0, 0, 0, 0})
	gc.SetStrokeColor(color.RGBA{0, 0, 0, 255})
	gc.SetLineWidth(2)

	// Draw segments.
	for _, seg := range m.allSegments {
		if seg.Prev == nil {
			continue
		}
		switch seg.Type {
		case Highway:
			gc.SetStrokeColor(color.RGBA{255, 0, 0, 255})
		case Street:
			gc.SetStrokeColor(color.RGBA{0, 255, 0, 255})
		case Footpath:
			gc.SetStrokeColor(color.RGBA{255, 255, 255, 255})
		}
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

// Generate generates the map.
func (m *Map) Generate() {
	// Create root segment.
	m.Root = &Segment{
		Point: vectors.NewVec2(0, 0),
		Type:  Highway,
		Step:  0,
	}

	// Initialize priority queue.
	m.queue = make(PriorityQueue, 0)
	heap.Init(&m.queue)
	heap.Push(&m.queue, m.Root)

	// Also generate a segment into the opposite direction.
	// This will ensure we grow the highway in both directions.
	opposite := &Segment{
		Point: vectors.NewVec2(-1, 0),
		Type:  Highway,
		Step:  0,
		Prev:  m.Root,
	}
	m.Root.Prev = opposite
	heap.Push(&m.queue, opposite)
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
	cfg := getSegTypeConfig(seg.Type)
	if seg.Prev != nil && m.rng.Float64() < cfg.BranchingChance {
		m.newSegment(seg, true)
	}
}

func (m *Map) newSegment(origin *Segment, branch bool) *Segment {
	segType := origin.Type
	if branch && segType < Footpath {
		segType++
	}
	config := getSegTypeConfig(segType)

	// Calculate length.
	dist := config.MinLength + config.MinLength*m.rng.Float64()*config.LengthVariation

	// Calculate angle.
	angle := config.MinAngle + m.rng.Float64()*config.AngleVariance

	// 50% chance to flip angle.
	if config.AngleReversal && m.rng.Float64() < 0.5 {
		angle *= -1
	}

	// If branch, add branching angle.
	if branch {
		originConfig := getSegTypeConfig(origin.Type)

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
	// NOTE: This should be maybe optimized by using a quadtree or something.
	for _, seg := range m.allSegments {
		if seg == origin || seg == newSeg {
			continue
		}

		// Check if there is an intersection. If so, we set the end point of the new segment to the intersection point.
		// And also we won't add the new segment to the queue.
		if ok, ip := seg.Intersects(newSeg); ok {
			if !newSeg.IsPointOnLine(ip) {
				continue
			}
			newSeg.Point = ip
			newSeg.Length = vectors.Dist2(ip, origin.Point)
			newSeg.End = true
			break
		}
	}

	// Add to queue (if not end).
	if !newSeg.End {
		heap.Push(&m.queue, newSeg)
	}

	// add branch to origin
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
