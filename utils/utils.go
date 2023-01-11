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
func (r IntRange) Random() int {
	return r.random(rand.Intn)
}

// random returns a random value between min and max using the supplied rand function.
func (r IntRange) random(rnd RandIntn) int {
	if r[0] == r[1] {
		return r[0]
	}
	if r[1]-r[0] <= 0 {
		log.Println("Warning: IntRange.Random() called with invalid range:", r)
	}
	return rnd(r[1]-r[0]) + r[0]
}

// InRange returns true if the given value is within the range.
func (r IntRange) InRange(value int) bool {
	return value >= r.Min() && value <= r.Max()
}

// Min returns the minimum value of the range.
func (r IntRange) Min() int {
	return r[0]
}

// Max returns the maximum value of the range.
func (r IntRange) Max() int {
	return r[1]
}

// FloatRange represents a minimum and maximum value.
type FloatRange [2]float64

// Random returns a random value between min and max.
func (r FloatRange) Random() float64 {
	return r.random(rand.Float64)
}

// random returns a random value between min and max using the provided
// rand function.
func (r FloatRange) random(rnd RandFloat64) float64 {
	if r[0] == r[1] {
		return r[0]
	}
	return r[0] + rnd()*(r[1]-r[0])
}

// InRange returns true if the given value is within the range.
func (r FloatRange) InRange(value float64) bool {
	return value >= r.Min() && value <= r.Max()
}

// Min returns the minimum value of the range.
func (r FloatRange) Min() float64 {
	return r[0]
}

// Max FloatRange the maximum value of the range.
func (r FloatRange) Max() float64 {
	return r[1]
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

// Number is the interface for some numeric types.
type Number interface {
	int | int64 | float64
}

// MinMax returns the minimum and maximum value of the given slice.
func MinMax[V Number](values []V) (min, max V) {
	if len(values) == 0 {
		return 0, 0
	}
	min = values[0]
	max = values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return
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
