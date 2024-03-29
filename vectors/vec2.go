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

// SetXY sets the x and y values of the vector.
func (v1 *Vec2) SetXY(x, y float64) {
	v1.X = x
	v1.Y = y
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

// Dot returns the dotproduct of two vectors.
func (v1 Vec2) Dot(v2 Vec2) float64 {
	return v1.Y*v2.Y + v1.X*v2.X
}

// Perpendicular returns the perpendicular vector.
func (v1 Vec2) Perpendicular() Vec2 {
	return Vec2{
		X: -v1.Y,
		Y: v1.X,
	}
}

// RotateAroundPoint rotates the vector around the given point.
func (v Vec2) RotateAroundPoint(p Vec2, angle float64) Vec2 {
	// https://stackoverflow.com/questions/2259476/rotating-a-point-about-another-point-2d
	s := math.Sin(angle)
	c := math.Cos(angle)
	x := v.X - p.X
	y := v.Y - p.Y
	return Vec2{
		X: x*c - y*s + p.X,
		Y: x*s + y*c + p.Y,
	}
}

// Normal returns the normal of the vector.
func (v Vec2) Normal() Vec2 {
	return Vec2{
		X: -v.Y,
		Y: v.X,
	}
}

// Angle returns the angle of the vector.
func (v Vec2) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

// SetLength sets the length of the vector.
func (v *Vec2) SetLength(length float64) Vec2 {
	return v.Mul(length / v.Len())
}

// LengthSquared returns the squared length of the vector.
func (v Vec2) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

// DistanceTo returns the distance to the given vector.
func (v Vec2) DistanceTo(v2 Vec2) float64 {
	return math.Sqrt((v.X-v2.X)*(v.X-v2.X) + (v.Y-v2.Y)*(v.Y-v2.Y))
}

// DistanceToSquared returns the squared distance to the given vector.
func (v Vec2) DistanceToSquared(v2 Vec2) float64 {
	return (v.X-v2.X)*(v.X-v2.X) + (v.Y-v2.Y)*(v.Y-v2.Y)
}

// AngleBetweenTwoVectors returns the angle between two vectors.
func AngleBetweenTwoVectors(v1, v2 Vec2) float64 {
	return math.Atan2(v1.X*v2.Y-v1.Y*v2.X, v1.X*v2.X+v1.Y*v2.Y)
}

// Negate negates the vector.
func (v Vec2) Negate() Vec2 {
	return Vec2{
		X: -v.X,
		Y: -v.Y,
	}
}

// Reflect reflects the vector at the given normal.
func (v Vec2) Reflect(normal Vec2) Vec2 {
	return v.Sub(normal.Mul(v.Dot(normal)))
}

// Normalize returns the normalized vector (with a length/magnitude of 1).
func (v Vec2) Normalize() Vec2 {
	x := v.X
	y := v.Y
	length := (x * x) + (y * y)
	if length > 0 {
		length = 1.0 / math.Sqrt(length)
		return Vec2{
			X: v.X * length,
			Y: v.Y * length,
		}
	}
	return Vec2{}
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

// Sub2 subtracts v2 from v1 and returns the result.
func Sub2(v1, v2 Vec2) Vec2 {
	return Vec2{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
	}
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

// Len returns the length of the segment.
func (s Segment) Len() float64 {
	return s.Start.DistanceTo(s.End)
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

// DistanceToPoint returns the distance to the given point.
func (s1 Segment) DistanceToPoint(p Vec2) float64 {
	var v, w, pb Vec2
	v.X = s1.End.X - s1.Start.X
	v.Y = s1.End.Y - s1.Start.Y
	w.X = p.X - s1.Start.X
	w.Y = p.Y - s1.Start.Y
	c1 := w.Dot(v)
	if c1 <= 0 {
		return p.DistanceTo(s1.Start)
	}
	c2 := v.Dot(v)
	if c2 <= c1 {
		return p.DistanceTo(s1.End)
	}
	b := c1 / c2
	pb.X = s1.Start.X + b*v.X
	pb.Y = s1.Start.Y + b*v.Y
	return p.DistanceTo(pb)
}

// ClosestPoint returns the closest point on the segment to the given point.
func (s1 Segment) ClosestPoint(p Vec2) Vec2 {
	// https://stackoverflow.com/questions/849211/shortest-distance-between-a-point-and-a-line-segment
	line := s1.End.Sub(s1.Start)
	len2 := line.LengthSquared()
	if len2 == 0 {
		return s1.Start
	}
	dot := (p.Sub(s1.Start)).Dot(line)
	t := math.Max(0, math.Min(1, dot/len2))
	return s1.Start.Add(line.Mul(t))
}

// Midpoint returns the midpoint of the segment.
func (s1 Segment) Midpoint() Vec2 {
	return Vec2{
		X: (s1.Start.X + s1.End.X) / 2,
		Y: (s1.Start.Y + s1.End.Y) / 2,
	}
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

func PointInTriangle(p Vec2, triangle []Vec2) bool {
	// https://stackoverflow.com/questions/2049582/how-to-determine-if-a-point-is-in-a-2d-triangle
	// Compute vectors
	v0 := Sub2(triangle[2], triangle[0])
	v1 := Sub2(triangle[1], triangle[0])
	v2 := Sub2(p, triangle[0])

	// Compute dot products
	dot00 := Dot2(v0, v0)
	dot01 := Dot2(v0, v1)
	dot02 := Dot2(v0, v2)
	dot11 := Dot2(v1, v1)
	dot12 := Dot2(v1, v2)

	// Compute barycentric coordinates
	invDenom := 1 / (dot00*dot11 - dot01*dot01)
	u := (dot11*dot02 - dot01*dot12) * invDenom
	v := (dot00*dot12 - dot01*dot02) * invDenom

	// Check if point is in triangle
	return (u >= 0) && (v >= 0) && (u+v < 1)
}

// AngleBetween returns the angle between the three points.
func AngleBetween(p1, p2, p3 Vec2) float64 {
	// https://stackoverflow.com/questions/14066933/direct-way-of-computing-clockwise-angle-between-2-vectors
	v1 := Sub2(p2, p1)
	v2 := Sub2(p3, p2)
	dot := Dot2(v1, v2)
	det := v1.X*v2.Y - v1.Y*v2.X
	return math.Atan2(det, dot)
}

// Ray2 represents a ray in 2D space.
type Ray2 struct {
	Origin Vec2 // Origin of the ray.
	Dir    Vec2 // Direction of the ray.
}

// Intersects returns true if the ray intersects with the line segment.
func (r Ray2) Intersects(s Segment) bool {
	// Determine if the ray (that is infinitely long) intersects with the line segment.
	denom := (s.End.Y-s.Start.Y)*r.Dir.X - (s.End.X-s.Start.X)*r.Dir.Y
	nume_a := (s.End.X-s.Start.X)*(r.Origin.Y-s.Start.Y) - (s.End.Y-s.Start.Y)*(r.Origin.X-s.Start.X)
	nume_b := (r.Dir.X)*(r.Origin.Y-s.Start.Y) - (r.Dir.Y)*(r.Origin.X-s.Start.X)
	if denom == 0 {
		if nume_a == 0 && nume_b == 0 {
			// Collinear
			return true
		}
		// Parallel
		return false
	}
	u_a := nume_a / denom
	u_b := nume_b / denom
	if u_a >= 0 && u_b >= 0 && u_b <= 1 {
		// Intersection
		return true
	}
	// No intersection
	return false
}
