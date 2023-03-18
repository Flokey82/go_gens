// Package vectors implements various vector thingies.
package vectors

import (
	"math"
)

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

// Rotate rotates the vector around the supplied axis by the supplied angle.
func (v1 Vec3) Rotate(axis Vec3, angle float64) Vec3 {
	// Normalize the axis vector.
	axis = axis.Normalize()

	x, y, z := v1.X, v1.Y, v1.Z
	x2, y2, z2 := axis.X, axis.Y, axis.Z

	cosTheta := math.Cos(angle)
	sinTheta := math.Sin(angle)
	dotProd := Dot3(v1, axis)

	xPrime := x2*dotProd*(1-cosTheta) + x*cosTheta + (-z2*y+y2*z)*sinTheta
	yPrime := y2*dotProd*(1-cosTheta) + y*cosTheta + (z2*x-x2*z)*sinTheta
	zPrime := z2*dotProd*(1-cosTheta) + z*cosTheta + (-y2*x+x2*y)*sinTheta

	return Vec3{
		X: xPrime,
		Y: yPrime,
		Z: zPrime,
	}
}

// Dot returns the dotproduct of the current and the supplied vector.
func (v1 Vec3) Dot(v2 Vec3) float64 {
	return Dot3(v1, v2)
}

// Sub subtracts v2 from the current vector and returns the result.
func (v1 Vec3) Sub(v2 Vec3) Vec3 {
	return Sub3(v1, v2)
}

// Up is a vector pointing up (0, 1, 0).
var Up = Vec3{0, 1, 0}

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

// Column-major order
// 0 3 6
// 1 4 7
// 2 5 8
type Mat3 [3][3]float64

func (m Mat3) MulVec3(v Vec3) Vec3 {
	return Vec3{
		X: m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z,
		Y: m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z,
		Z: m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z,
	}
}
