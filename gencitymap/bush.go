package gencitymap

import (
	"log"

	"github.com/Flokey82/go_gens/vectors"
)

type BushOptions struct {
	OnFound func(p vectors.Vec2, interior []vectors.Segment) bool
}

/**
 * https://github.com/anvaka/isect/blob/master/bush.js
 *
 * This implementation is inspired by discussion here
 * https://twitter.com/mourner/status/1049325199617921024 and
 * here https://github.com/anvaka/isect/issues/1
 *
 * It builds an index of all segments using static spatial index
 * and then for each segment it queries overlapping rectangles.
 */
func bush(lines []vectors.Segment, options *BushOptions) *Bush {
	var results []intersection
	var asyncState *BushAsyncState
	index := NewFlatbush(len(lines), 16)
	for _, line := range lines {
		index.Add(line.Start.X, line.Start.Y, line.End.X, line.End.Y)
	}
	index.Finish()
	return &Bush{
		lines:      lines,
		index:      index,
		results:    results,
		asyncState: asyncState,
	}
}

func (b *Bush) Run() []intersection {
	for i, line := range b.lines {
		if b.checkIntersection(line, i) {
			return b.results
		}
	}
	return b.results
}

func (b *Bush) checkIntersection(currentSegment vectors.Segment, currentId int) bool {
	minX := currentSegment.Start.X
	maxX := currentSegment.End.X
	minY := currentSegment.Start.Y
	maxY := currentSegment.End.Y
	if minX > maxX {
		minX, maxX = maxX, minX
	}
	if minY > maxY {
		minY, maxY = maxY, minY
	}
	ids := b.index.Query(minX, minY, maxX, maxY, nil)
	for _, segmentIndex := range ids {
		if segmentIndex == currentId {
			continue
		}
		otherSegment := b.lines[segmentIndex]
		if ok, intersection := currentSegment.Intersects(otherSegment); ok {
			log.Println("Intersection found", intersection)
			if b.reportIntersection(intersection, []vectors.Segment{currentSegment, otherSegment}) {
				log.Println("returning true")
				return true
			}
		} else {
			log.Println("No intersection found")
		}
	}
	log.Println("returning false")
	return false
}

func (b *Bush) Step() bool {
	if b.asyncState == nil {
		b.asyncState = &BushAsyncState{I: 0}
	}
	test := b.lines[b.asyncState.I]
	b.checkIntersection(test, b.asyncState.I)
	b.asyncState.I += 1
	return b.asyncState.I < len(b.lines)
}

type BushAsyncState struct {
	I int
}

type Bush struct {
	lines      []vectors.Segment
	index      *Flatbush
	results    []intersection
	asyncState *BushAsyncState
}

func (b *Bush) addToIndex(line vectors.Segment) {
	minX := line.Start.X
	maxX := line.End.X
	minY := line.Start.Y
	maxY := line.End.Y
	if minX > maxX {
		minX, maxX = maxX, minX
	}
	if minY > maxY {
		minY, maxY = maxY, minY
	}
	b.index.Add(minX, minY, maxX, maxY)
}

func (b *Bush) reportIntersection(p vectors.Vec2, interior []vectors.Segment) bool {
	b.results = append(b.results, intersection{Point: p, Segments: interior})
	log.Println("number of results", len(b.results))
	return false
}
