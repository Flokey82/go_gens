package genworldvoronoi

import (
	"math"
	"sort"

	"github.com/Flokey82/go_gens/vectors"
)

// getCentroidOfTriangle returns the centroid of a triangle defined by
// the xyz coordinates a, b, c as a vectors.Vec3.
func getCentroidOfTriangle(a, b, c []float64) vectors.Vec3 {
	return vectors.Vec3{
		X: (a[0] + b[0] + c[0]) / 3,
		Y: (a[1] + b[1] + c[1]) / 3,
		Z: (a[2] + b[2] + c[2]) / 3,
	}.Normalize()
}

// dist2 returns the eucledian distance between two points.
func dist2(a, b [2]float64) float64 {
	xDiff := a[0] - b[0]
	yDiff := a[1] - b[1]
	return math.Sqrt(xDiff*xDiff + yDiff*yDiff)
}

func dot2(a, b [2]float64) float64 {
	return a[0]*b[0] + a[1]*b[1]
}

func len2(a [2]float64) float64 {
	return math.Sqrt(a[0]*a[0] + a[1]*a[1])
}

func normal2(v [2]float64) [2]float64 {
	l := 1.0 / len2(v)
	return [2]float64{
		v[0] * l,
		v[1] * l,
	}
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
	return vectors.Vec3{
		X: xyz[0],
		Y: xyz[1],
		Z: xyz[2],
	}
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func radToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

func calcVecFromLatLong(lat1, lon1, lat2, lon2 float64) [2]float64 {
	// convert to radians
	lat1 = degToRad(lat1)
	lon1 = degToRad(lon1)
	lat2 = degToRad(lat2)
	lon2 = degToRad(lon2)
	return [2]float64{
		math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(lon2-lon1), // X
		math.Sin(lon2-lon1) * math.Cos(lat2),                                              // Y
	}
}

// haversine returns the great arc distance between two lat/long pairs.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	// distance between latitudes and longitudes
	dLat := degToRad(lat2 - lat1)
	dLon := degToRad(lon2 - lon1)

	// convert to radians
	lat1 = degToRad(lat1)
	lat2 = degToRad(lat2)

	// apply formula
	a := math.Pow(math.Sin(dLat/2), 2) + math.Pow(math.Sin(dLon/2), 2)*math.Cos(lat1)*math.Cos(lat2)
	return 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// heronsTriArea returns the area of a triangle given the three sides.
// See: https://www.mathopenref.com/heronsformula.html
func heronsTriArea(a, b, c float64) float64 {
	p := (a + b + c) / 2
	return math.Sqrt(p * (p - a) * (p - b) * (p - c))
}
