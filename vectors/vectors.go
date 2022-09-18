// Package vectors implements various vector thingies.
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

// Dot2 returns the dotproduct of two vectors.
func Dot2(v1, v2 Vec2) float64 {
	return v1.Y*v2.Y + v1.X*v2.X
}

// Dist2 returns the eucledian distance between two vectors.
func Dist2(a, b Vec2) float64 {
	return math.Hypot(a.X-b.X, a.Y-b.Y)
}

// Vec3 represents a three dimensional vector.
type Vec3 struct {
	X float64
	Y float64
	Z float64
}

// NewVec3 returns a new vector3.
func NewVec3(x, y, z float64) Vec3 {
	return Vec3{
		X: x,
		Y: y,
		Z: z,
	}
}

// Len returns the vector's length (or magnitude).
//
// NOTE: This is equivalent to the square root of the sum of the squares
// of all elements. E.G. for a Vec2 it's math.Hypot(v[0], v[1]).
func (v1 Vec3) Len() float64 {
	return math.Sqrt(v1.X*v1.X + v1.Y*v1.Y + v1.Z*v1.Z)
}

// Normalize returns the normalized vector (with a length/magnitude of 1).
func (v1 Vec3) Normalize() Vec3 {
	l := 1.0 / v1.Len()
	return Vec3{
		X: v1.X * l,
		Y: v1.Y * l,
		Z: v1.Z * l,
	}
}

// Cross returns the crossproduct of the current and the supplied vector.
func (v1 Vec3) Cross(v2 Vec3) Vec3 {
	return Cross3(v1, v2)
}

// Add adds the current vector with v2 and returns the result.
func (v1 Vec3) Add(v2 Vec3) Vec3 {
	return Add3(v1, v2)
}

// Mul multiplies the current vector with c and returns the result.
func (v1 Vec3) Mul(c float64) Vec3 {
	return Vec3{
		X: v1.X * c,
		Y: v1.Y * c,
		Z: v1.Z * c,
	}
}

// AddToThis adds v2 to the current vector.
func (v1 *Vec3) AddToThis(v2 Vec3) {
	v1.X += v2.X
	v1.Y += v2.Y
	v1.Z += v2.Z
}

// Add3 adds two Vec3 vectors and returns the result.
func Add3(v1, v2 Vec3) Vec3 {
	return Vec3{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
		Z: v1.Z + v2.Z,
	}
}

// Sub3 subtracts v2 from v1 and returns the result.
func Sub3(v1, v2 Vec3) Vec3 {
	return Vec3{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
		Z: v1.Z - v2.Z,
	}
}

// Dist3 returns the eucledian distance between two vectors.
func Dist3(a, b Vec3) float64 {
	xDiff := a.X - b.X
	yDiff := a.Y - b.Y
	zDiff := a.Z - b.Z
	return math.Sqrt(xDiff*xDiff + yDiff*yDiff + zDiff*zDiff)
}

// Dot3 returns the dotproduct of two vectors.
func Dot3(v1, v2 Vec3) float64 {
	return v1.Y*v2.Y + v1.X*v2.X + v1.Z*v2.Z
}

// Cross3 returns the crossproduct of two vectors.
func Cross3(v1, v2 Vec3) Vec3 {
	return Vec3{
		X: v1.Y*v2.Z - v1.Z*v2.Y,
		Y: v1.Z*v2.X - v1.X*v2.Z,
		Z: v1.X*v2.Y - v1.Y*v2.X,
	}
}

// Cross3XYZ returns the cross product of two sets of coordinates.
func Cross3XYZ(v1x, v1y, v1z, v2x, v2y, v2z float64) Vec3 {
	return Vec3{
		X: v1y*v2z - v1z*v2y,
		Y: v1z*v2x - v1x*v2z,
		Z: v1x*v2y - v1y*v2x,
	}
}
