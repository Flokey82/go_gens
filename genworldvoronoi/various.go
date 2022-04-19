package genworldvoronoi

import (
	"math"
	"sort"

	"github.com/Flokey82/go_gens/vectors"
)

// dist2 returns the eucledian distance between two points.
func dist2(a, b [2]float64) float64 {
	xDiff := a[0] - b[0]
	yDiff := a[1] - b[1]
	return float64(math.Sqrt(float64(xDiff*xDiff + yDiff*yDiff)))
}

// min is the int equivalent of math.Min(a, b).
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// minMax returns the smallest and largest value in hm.
func minMax(hm []float64) (float64, float64) {
	if len(hm) == 0 {
		return 0, 0
	}
	min, max := hm[0], hm[0]
	for _, h := range hm {
		if h > max {
			max = h
		}
		if h < min {
			min = h
		}
	}
	return min, max
}

func convToMap(in []int) map[int]bool {
	res := make(map[int]bool)
	for _, v := range in {
		res[v] = true
	}
	return res
}

func convToArray(in map[int]bool) []int {
	var res []int
	for v := range in {
		res = append(res, v)
	}
	sort.Ints(res)
	return res
}

// convToVec3 converts a float slice containing 3 values into a vectors.Vec3.
func convToVec3(xyz []float64) vectors.Vec3 {
	return vectors.Vec3{xyz[0], xyz[1], xyz[2]}
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func radToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}
