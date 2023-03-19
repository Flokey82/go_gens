package vectors

import (
	"math"
	"math/rand"
)

// IVec2 represents an integer vector of type int64.
type IVec2 struct {
	X int64
	Y int64
}

// NewIVec2 returns a new vector.
func NewIVec2(x, y int64) IVec2 {
	return IVec2{
		X: x,
		Y: y,
	}
}

// RandomVec2 returns a randomized vector.
func RandomVec2(scale float64) Vec2 {
	return NewVec2(scale*rand.Float64(), scale*rand.Float64())
}

// Vec2 represents a two dimensional vector.
type Vec2 struct {
	X float64
	Y float64
}

// NewVec2 returns a new two dimensional vector.
func NewVec2(x, y float64) Vec2 {
	return Vec2{
		X: x,
		Y: y,
	}
}

// Add adds the current vector with v2 and returns the result.
func (v1 Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
	}
}

// Sub subtracts v2 from the current vector and returns the result.
func (v1 Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
	}
}

// AddToThis adds v2 to the current vector.
func (v1 *Vec2) AddToThis(v2 Vec2) {
	v1.X += v2.X
	v1.Y += v2.Y
}

// Mul multiplies the current vector with c and returns the result.
func (v1 Vec2) Mul(c float64) Vec2 {
	return Vec2{
		X: v1.X * c,
		Y: v1.Y * c,
	}
}

// MulWithThis multiplies the current vector with c.
func (v1 *Vec2) MulWithThis(c float64) {
	v1.X *= c
	v1.Y *= c
}

// Len returns the length (or magnitude) of the vector.
func (v1 Vec2) Len() float64 {
	return math.Hypot(v1.X, v1.Y)
}

// Equal returns true if the vectors are equal.
func (v1 Vec2) Equal(v2 Vec2) bool {
	return v1.X == v2.X && v1.Y == v2.Y
}

// Equalish returns true if the vectors are equalish.
func (v1 Vec2) Equalish(v2 Vec2) bool {
	return Equalish(v1.X, v2.X) && Equalish(v1.Y, v2.Y)
}

// Rotate rotates the vector by the given angle (in degrees).
func (v1 Vec2) Rotate(angle float64) Vec2 {
	angleRad := angle * math.Pi / 180
	cos := math.Cos(angleRad)
	sin := math.Sin(angleRad)
	return Vec2{
		X: v1.X*cos - v1.Y*sin,
		Y: v1.X*sin + v1.Y*cos,
	}
}

// Cross returns the crossproduct of two vectors.
func (v1 Vec2) Cross(v2 Vec2) float64 {
	return v1.X*v2.Y - v1.Y*v2.X
}

// Normalize returns the normalized vector (with a length/magnitude of 1).
func Normalize(vec Vec2) (dest Vec2) {
	x := vec.X
	y := vec.Y
	length := (x * x) + (y * y)
	if length > 0 {
		length = 1.0 / math.Sqrt(length)
		dest.X = vec.X * length
		dest.Y = vec.Y * length
	}
	return dest
}

// Cross2 returns the crossproduct of two vectors.
func Cross2(v1, v2 Vec2) float64 {
	return v1.X*v2.Y - v1.Y*v2.X
}

// Dot2 returns the dotproduct of two vectors.
func Dot2(v1, v2 Vec2) float64 {
	return v1.Y*v2.Y + v1.X*v2.X
}

// Dist2 returns the eucledian distance between two vectors.
func Dist2(a, b Vec2) float64 {
	return math.Hypot(a.X-b.X, a.Y-b.Y)
}

// MinDegreesDifference2 returns the minimum difference between two angles in degrees.
func MinDegreesDifference2(a, b float64) float64 {
	bottom := math.Mod(math.Abs(a-b), 180)
	return math.Min(bottom, bottom-180)
}

// EqualVec2 returns true if the two vectors are equal (within a small epsilon).
func EqualVec2(a, b Vec2) bool {
	return math.Abs(a.X-b.X) < 0.00001 && math.Abs(a.Y-b.Y) < 0.00001
}

// Equalish returns true if the two values are equalish.
func Equalish(a, b float64) bool {
	return math.Abs(a-b) < 0.00001
}

// Angle2 returns the angle of a vector in degrees.
func Angle2(v Vec2) float64 {
	return math.Atan2(v.Y, v.X) * 180 / math.Pi
}

// Segment represents a line segment.
type Segment struct {
	Start Vec2
	End   Vec2
}

// NewSegment returns a new segment.
func NewSegment(start, end Vec2) Segment {
	return Segment{
		Start: start,
		End:   end,
	}
}

// Intersect returns true if the two segments intersect.
func (s1 Segment) Intersect(s2 Segment) bool {
	return s1.IntersectPoint(s2) != nil
}

// IntersectPoint returns the intersection point of the two segments.
func (s1 Segment) IntersectPoint(s2 Segment) *Vec2 {
	denom := (s2.End.Y-s2.Start.Y)*(s1.End.X-s1.Start.X) - (s2.End.X-s2.Start.X)*(s1.End.Y-s1.Start.Y)
	if denom == 0 {
		return nil
	}
	ua := ((s2.End.X-s2.Start.X)*(s1.Start.Y-s2.Start.Y) - (s2.End.Y-s2.Start.Y)*(s1.Start.X-s2.Start.X)) / denom
	ub := ((s1.End.X-s1.Start.X)*(s1.Start.Y-s2.Start.Y) - (s1.End.Y-s1.Start.Y)*(s1.Start.X-s2.Start.X)) / denom
	if ua >= 0 && ua <= 1 && ub >= 0 && ub <= 1 {
		x := s1.Start.X + ua*(s1.End.X-s1.Start.X)
		y := s1.Start.Y + ua*(s1.End.Y-s1.Start.Y)
		return &Vec2{X: x, Y: y}
	}
	return nil
}

// IsPointOnLine returns true if the point is on the line.
func (s1 Segment) IsPointOnLine(p Vec2) bool {
	if s1.Start.X == s1.End.X {
		// vertical line
		return p.X == s1.Start.X
	}
	if s1.Start.Y == s1.End.Y {
		// horizontal line
		return p.Y == s1.Start.Y
	}
	return (p.X-s1.Start.X)*(s1.End.Y-s1.Start.Y) == (p.Y-s1.Start.Y)*(s1.End.X-s1.Start.X)
}

// Intersects returns true if the line intersects with the other line.
func (s1 Segment) Intersects(l2 Segment) (bool, Vec2) {
	// Check if one of the points is somewhere on the other line.
	if l2.IsPointOnLine(s1.Start) {
		return true, s1.Start
	}
	if l2.IsPointOnLine(s1.End) {
		return true, s1.End
	}
	if s1.IsPointOnLine(l2.Start) {
		return true, l2.Start
	}
	if s1.IsPointOnLine(l2.End) {
		return true, l2.End
	}
	denominator := (l2.End.Y-l2.Start.Y)*(s1.End.X-s1.Start.X) - (l2.End.X-l2.Start.X)*(s1.End.Y-s1.Start.Y)
	if denominator == 0 {
		return false, Vec2{}
	}
	uA := ((l2.End.X-l2.Start.X)*(s1.Start.Y-l2.Start.Y) - (l2.End.Y-l2.Start.Y)*(s1.Start.X-l2.Start.X)) / denominator
	uB := ((s1.End.X-s1.Start.X)*(s1.Start.Y-l2.Start.Y) - (s1.End.Y-s1.Start.Y)*(s1.Start.X-l2.Start.X)) / denominator
	if uA >= 0 && uA <= 1 && uB >= 0 && uB <= 1 {
		return true, Vec2{
			X: s1.Start.X + (uA * (s1.End.X - s1.Start.X)),
			Y: s1.Start.Y + (uA * (s1.End.Y - s1.Start.Y)),
		}
	}
	return false, Vec2{}
}
