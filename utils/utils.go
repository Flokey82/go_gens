// Package utils provides some commonly used functionality.
package utils

import (
	"log"
	"math/rand"
)

// IncrementalAvrg computes the incremental average given the old value, the new value and the
// previous number of samples.
func IncrementalAvrg(oldVal, newVal float64, count int) float64 {
	return oldVal + (newVal-oldVal)/(float64(count))
}

// WeightedAvrg computes the weighted average given the old value, the new value and the
// weight factor.
func WeightedAvrg(oldVal, newVal, weightFactor float64) float64 {
	return oldVal + weightFactor*(newVal-oldVal)
}

// IntRange represents a minimum and maximum value.
type IntRange [2]int

// Random returns a random value between min and max.
func (ir IntRange) Random() int {
	return ir.random(rand.Intn)
}

// random returns a random value between min and max using the supplied rand function.
func (ir IntRange) random(rnd RandIntn) int {
	if ir[0] == ir[1] {
		return ir[0]
	}
	if ir[1]-ir[0] <= 0 {
		log.Println("Warning: IntRange.Random() called with invalid range:", ir)
	}
	return rnd(ir[1]-ir[0]) + ir[0]
}

// InRange returns true if the given value is within the range.
func (ir IntRange) InRange(value int) bool {
	return value >= ir.Min() && value <= ir.Max()
}

// Min returns the minimum value of the range.
func (m IntRange) Min() int {
	return m[0]
}

// Max returns the maximum value of the range.
func (m IntRange) Max() int {
	return m[1]
}

// FloatRange represents a minimum and maximum value.
type FloatRange [2]float64

// Random returns a random value between min and max.
func (m FloatRange) Random() float64 {
	return m.random(rand.Float64)
}

// random returns a random value between min and max using the provided
// rand function.
func (m FloatRange) random(rnd RandFloat64) float64 {
	if m[0] == m[1] {
		return m[0]
	}
	return m[0] + rnd()*(m[1]-m[0])
}

// InRange returns true if the given value is within the range.
func (m FloatRange) InRange(value float64) bool {
	return value >= m.Min() && value <= m.Max()
}

// Min returns the minimum value of the range.
func (m FloatRange) Min() float64 {
	return m[0]
}

// Max FloatRange the maximum value of the range.
func (m FloatRange) Max() float64 {
	return m[1]
}

// Min is the int equivalent of math.Min(a, b).
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max is the int equivalent of math.Max(a, b).
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Abs is the int equivalent of math.Abs(a).
func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// Gauss returns a normally distributed random number with the given mean and standard deviation.
// NOTE: mean: mu, stdDev: sigma
// See: https://en.wikipedia.org/wiki/Normal_distribution
func Gauss(mean, stdDev float64) float64 {
	return gauss(mean, stdDev, rand.NormFloat64)
}

// gauss returns a normally distributed random number with the given mean and standard deviation
// using the provided rand function (rand.NormFloat64).
func gauss(mean, stdDev float64, rnd RandFloat64) float64 {
	return rnd()*stdDev + mean
}

// Type alias for rand.Int63().
type Rand63 func() int64

// Type alias for rand.Float64().
type RandFloat64 func() float64

// Type alias for rand.Intn().
type RandIntn func(int) int
