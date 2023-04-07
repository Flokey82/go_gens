package gencitymap

import (
	"math"

	"github.com/Flokey82/go_gens/vectors"
)

// Segment is a segment of a road.
type Segment struct {
	Point    vectors.Vec2 // end point of the segment
	Length   float64      // length of the segment
	Next     *Segment     // next segment in the same road
	Prev     *Segment     // previous segment in the same road
	Type     RoadType     // type of the segment
	Branches []*Segment   // branches from the same road
	Step     int          // Iteration when this segment was created
	End      bool         // true if this segment is the end of a road
}

// GetVector returns the vector of the segment.
func (s *Segment) GetVector() vectors.Vec2 {
	if s.Prev == nil {
		return vectors.NewVec2(0, 0)
	}
	return s.Point.Sub(s.Prev.Point)
}

// Intersects returns true if the segment intersects with the given segment and
// the intersection point.
func (s *Segment) Intersects(seg *Segment) (bool, vectors.Vec2) {
	if s.Prev == nil || seg.Prev == nil {
		return false, vectors.Vec2{}
	}

	ok, ip, _ := doLineSegmentsIntersectVec2(s.Prev.Point, s.Point, seg.Prev.Point, seg.Point)
	return ok, ip
}

// IsPointOnLine returns true if the given point is on the line of the segment (within 0.0001)
func (s *Segment) IsPointOnLine(p vectors.Vec2) bool {
	return math.Abs((p.X-s.Prev.Point.X)*(s.Point.Y-s.Prev.Point.Y)-(p.Y-s.Prev.Point.Y)*(s.Point.X-s.Prev.Point.X)) < 0.0001
}

// Split splits the segment at the given point and returns the new segment.
func (s *Segment) Split(p vectors.Vec2) *Segment {
	if s.Prev == nil {
		return nil
	}

	// create new segment
	newSeg := &Segment{
		Point:    p,
		Length:   s.Point.Sub(p).Len(),
		Next:     s,
		Prev:     s.Prev,
		Type:     s.Type,
		Branches: s.Branches,
		Step:     s.Step,
		End:      false,
	}

	// update old segment
	s.Prev.Next = newSeg
	s.Prev = newSeg
	s.Length = s.Point.Sub(newSeg.Point).Len()

	return newSeg
}

type RoadType int

// SegmentTypeConfig is the configuration for a segment type.
type SegmentTypeConfig struct {
	LengthMin            float64 // minimum length of a segment of this type
	LengthVariation      float64 // maximum length variation of a segment of this type (0.1 = 10%)
	AngleMin             float64 // minimum angle of a segment extension in degrees (10.0 = 10°)
	AngleVariance        float64 // maximum angle variation of a segment extension in degrees (10.0 = 10°)
	AngleReversal        bool    // allow reverse angle of the road
	BranchingChance      float64 // chance of branching (0.1 = 10%)
	BranchingAngle       float64 // angle of subbranches in degrees (10.0 = 10°)
	BranchingReversal    bool    // allow reverse branching direction of the road
	BranchSameType       bool    // allow branching to the same type of road
	BranchSameTypeChance float64 // chance of branching to the same type of road (0.1 = 10%)
}

var (
	HighwayConfig = SegmentTypeConfig{
		LengthMin:         200,
		LengthVariation:   0.1,
		AngleMin:          10.0,
		AngleVariance:     1.0,
		AngleReversal:     true,
		BranchingChance:   0.3,
		BranchingAngle:    90.0,
		BranchingReversal: true,
	}
	StreetConfig = SegmentTypeConfig{
		LengthMin:            100,
		LengthVariation:      0.5,
		AngleVariance:        4.0,
		AngleReversal:        true,
		BranchingChance:      0.7,
		BranchingAngle:       90.0,
		BranchingReversal:    true,
		BranchSameType:       true,
		BranchSameTypeChance: 0.2,
	}
	FootpathConfig = SegmentTypeConfig{
		LengthMin:       50,
		LengthVariation: 0.6,
		AngleVariance:   3.0,
		AngleReversal:   true,
		BranchingChance: 0.5,
		BranchingAngle:  90.0,
	}
)

func doLineSegmentsIntersectVec2(p0, p1, p2, p3 vectors.Vec2) (ok bool, res vectors.Vec2, t float64) {
	b := p1.Sub(p0)
	d := p3.Sub(p2)
	f := vectors.Cross2(p2.Sub(p0), b)
	k := vectors.Cross2(b, d)
	if f == 0 && k == 0 || k == 0 {
		return false, vectors.Vec2{}, 0
	}
	f /= k
	e := vectors.Cross2(p2.Sub(p0), d) / k
	intersect := 0.00001 < e && 0.99999 > e && 0.00001 < f && 0.99999 > f
	if intersect {
		return true, vectors.NewVec2(p0.X+e*b.X, p0.Y+e*b.Y), e
	}
	return false, vectors.Vec2{}, 0
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

type MapConfig struct {
	SeedRoots func() []*Segment
	Rules     []*SegmentTypeConfig
}

func (mc *MapConfig) getTypeConfig(segType RoadType) *SegmentTypeConfig {
	if segType < 0 {
		return mc.Rules[0]
	}
	if segType >= RoadType(len(mc.Rules)) {
		return mc.Rules[len(mc.Rules)-1]
	}
	return mc.Rules[segType]
}

var DefaultMapConfig = &MapConfig{
	SeedRoots: func() []*Segment {
		root1 := &Segment{
			Point:  vectors.NewVec2(0, 0),
			Length: 100,
			Type:   0,
		}
		/*
			opposite := &Segment{
				Point: vectors.NewVec2(-1, 0),
				Type:  0,
				Step:  0,
				Prev:  root1,
			}
			root1.Prev = opposite
			return []*Segment{root1, opposite}
		*/
		return []*Segment{root1}
	}, Rules: []*SegmentTypeConfig{
		&HighwayConfig,
		&StreetConfig,
		&FootpathConfig,
	},
}
