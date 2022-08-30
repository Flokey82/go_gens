package genflora

import (
	"log"

	"math/rand"
)

type IntRange [2]int

// Random returns a random value between min and max.
func (ir IntRange) Random() int {
	if ir[0] == ir[1] {
		return ir[0]
	}
	if ir[1]-ir[0] <= 0 {
		log.Println("Warning: IntRange.Random() called with invalid range:", ir)
	}
	return rand.Intn(ir[1]-ir[0]) + ir[0]
}

// Min returns the minimum value of the range.
func (m IntRange) Min() int {
	return m[0]
}

// Max returns the maximum value of the range.
func (m IntRange) Max() int {
	return m[1]
}

// Range represents a minimum and maximum value.
type Range [2]float64

// Random returns a random value between min and max.
func (m Range) Random() float64 {
	if m[0] == m[1] {
		return m[0]
	}
	return m[0] + rand.Float64()*(m[1]-m[0])
}

// Min returns the minimum value of the range.
func (m Range) Min() float64 {
	return m[0]
}

// Max returns the maximum value of the range.
func (m Range) Max() float64 {
	return m[1]
}
