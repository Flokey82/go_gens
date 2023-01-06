package genworldvoronoi

import (
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/Flokey82/go_gens/vectors"
)

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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

func distToSegment2(v, w, p [2]float64) float64 {
	l2 := dist2(v, w)
	if l2 == 0 {
		// If the line segment has a length of 0, we can just return
		// the distance between the point and any of the two line
		// segment points.
		return dist2(p, v)
	}
	t := ((p[0]-v[0])*(w[0]-v[0]) + (p[1]-v[1])*(w[1]-v[1])) / (l2 * l2)
	t = math.Max(0, math.Min(1, t))
	return dist2(p, [2]float64{v[0] + t*(w[0]-v[0]), v[1] + t*(w[1]-v[1])})
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

func minMax64(hm []int64) (int64, int64) {
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

// crossArc calculates the shortest distance between an arc (defined by p1 and p2)
// and a third point, p3. The input is expected in degrees.
// See: https://stackoverflow.com/questions/32771458/distance-from-lat-lng-point-to-minor-arc-segment
func crossArc(lat1, lon1, lat2, lon2, lat3, lon3 float64) float64 {
	// dis Finds the distance between two lat/lon points.
	dis := func(latA, lonA, latB, lonB float64) float64 {
		return math.Acos(math.Sin(latA)*math.Sin(latB) + math.Cos(latA)*math.Cos(latB)*math.Cos(lonB-lonA))
	}

	// bearing Finds the bearing from one lat/lon point to another.
	bearing := func(latA, lonA, latB, lonB float64) float64 {
		return math.Atan2(math.Sin(lonB-lonA)*math.Cos(latB), math.Cos(latA)*math.Sin(latB)-math.Sin(latA)*math.Cos(latB)*math.Cos(lonB-lonA))
	}

	lat1 = degToRad(lat1)
	lat2 = degToRad(lat2)
	lat3 = degToRad(lat3)
	lon1 = degToRad(lon1)
	lon2 = degToRad(lon2)
	lon3 = degToRad(lon3)

	// Prerequisites for the formulas
	bear12 := bearing(lat1, lon1, lat2, lon2)
	bear13 := bearing(lat1, lon1, lat3, lon3)
	dis13 := dis(lat1, lon1, lat3, lon3)

	diff := math.Abs(bear13 - bear12)
	if diff > math.Pi {
		diff = 2*math.Pi - diff
	}
	// Is relative bearing obtuse?
	if diff > math.Pi/2 {
		return dis13
	}
	// Find the cross-track distance.
	dxt := math.Asin(math.Sin(dis13) * math.Sin(bear13-bear12))

	// Is p4 beyond the arc?
	dis12 := dis(lat1, lon1, lat2, lon2)
	dis14 := math.Acos(math.Cos(dis13) / math.Cos(dxt))
	if dis14 > dis12 {
		return dis(lat2, lon2, lat3, lon3)
	}
	return math.Abs(dxt)
}

// heronsTriArea returns the area of a triangle given the three sides.
// See: https://www.mathopenref.com/heronsformula.html
func heronsTriArea(a, b, c float64) float64 {
	p := (a + b + c) / 2
	return math.Sqrt(p * (p - a) * (p - b) * (p - c))
}

func isInIntList(l []int, i int) bool {
	for _, v := range l {
		if v == i {
			return true
		}
	}
	return false
}

// initFloatSlice returns a slice of floats of the given size, initialized to -1.
func initFloatSlice(size int) []float64 {
	res := make([]float64, size)
	for i := range res {
		res[i] = -1
	}
	return res
}

// initRegionSlice returns a slice of ints of the given size, initialized to -1.
func initRegionSlice(size int) []int {
	res := make([]int, size)
	for i := range res {
		res[i] = -1
	}
	return res
}

// initTimeSlice returns a slice of int64s of the given size, initialized to -1.
func initTimeSlice(size int) []int64 {
	res := make([]int64, size)
	for i := range res {
		res[i] = -1
	}
	return res
}

// mergeIndexSegments matches up the ends of the segments (region pairs) and returns
// a slice containing all continuous, connected segments as sequence of connected regions.
func mergeIndexSegments(segs [][2]int) [][]int {
	log.Println("start adj")
	adj := make(map[int][]int)
	var maxSegIdx int
	for i := 0; i < len(segs); i++ {
		seg := segs[i]
		adj[seg[0]] = append(adj[seg[0]], seg[1])
		adj[seg[1]] = append(adj[seg[1]], seg[0])
		if seg[0] > maxSegIdx {
			maxSegIdx = seg[0]
		}
		if seg[1] > maxSegIdx {
			maxSegIdx = seg[1]
		}
	}
	done := make([]bool, maxSegIdx)
	var paths [][]int
	var path []int
	log.Println("start paths")
	for {
		if path == nil {
			for i := 0; i < len(segs); i++ {
				if done[i] {
					continue
				}
				done[i] = true
				path = []int{segs[i][0], segs[i][1]}
				break
			}
			if path == nil {
				break
			}
		}
		var changed bool
		for i := 0; i < len(segs); i++ {
			if done[i] {
				continue
			}
			if len(adj[path[0]]) == 2 {
				if segs[i][0] == path[0] {
					path = unshiftIndexPath(path, segs[i][1])
					done[i] = true
					changed = true
					break
				}
				if segs[i][1] == path[0] {
					path = unshiftIndexPath(path, segs[i][0])
					done[i] = true
					changed = true
					break
				}
			}
			if len(adj[path[len(path)-1]]) == 2 {
				if segs[i][0] == path[len(path)-1] {
					path = append(path, segs[i][1])
					done[i] = true
					changed = true
					break
				}
				if segs[i][1] == path[len(path)-1] {
					path = append(path, segs[i][0])
					done[i] = true
					changed = true
					break
				}
			}
		}
		if !changed {
			//log.Println("done paths", len(paths), "pathlen", len(path))
			paths = append(paths, path)
			path = nil
		}
	}
	return paths
}

func unshiftIndexPath(path []int, p int) []int {
	res := make([]int, len(path)+1)
	res[0] = p
	copy(res[1:], path)
	return res
}

// round value to d decimals
func roundToDecimals(v, d float64) float64 {
	m := math.Pow(10, d)
	return math.Round(v*m) / m
}

// turn weighted array into simple array
func weightedToArray(weighted map[string]int) []string {
	var res []string
	for key, weight := range weighted {
		for j := 0; j < weight; j++ {
			res = append(res, key)
		}
	}
	return res
}

// probability shorthand
func P(probability float64) bool {
	if probability >= 1.0 {
		return true
	}
	if probability <= 0 {
		return false
	}
	return rand.Float64() < probability
}

// queueEntry is a single entry in the priority queue.
type queueEntry struct {
	index       int     // index of the item in the heap.
	score       float64 // priority of the item in the queue.
	origin      int     // origin region / ID
	destination int     // destination region / ID
}

// ascPriorityQueue implements heap.Interface and holds Items.
// Priority is ascending (lowest score first).
type ascPriorityQueue []*queueEntry

func (pq ascPriorityQueue) Len() int { return len(pq) }

func (pq ascPriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	// return pq[i].score > pq[j].score // 3, 2, 1

	// We want Pop to give us the lowest, not highest, priority so we use less than here.
	return pq[i].score < pq[j].score // 1, 2, 3
}

func (pq *ascPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *ascPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*queueEntry)
	item.index = n
	*pq = append(*pq, item)
}

func (pq ascPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

// getDiversionFromRange returns the amount a value diverges from a range.
func getDiversionFromRange(x float64, rng [2]float64) float64 {
	if x < rng[0] {
		return rng[0] - x
	}
	if x > rng[1] {
		return x - rng[1]
	}
	return 0
}

// getRangeFit returns 1 if a value fits within the range, and a value between 0 and 1 if it doesn't.
// 0 means the value is outside the range, 1 means it's at the edge of the range.
// If the value is more than 20% outside the range, -1 is returned.
// If the value is exactly 20% outside the range, 0 is returned.
func getRangeFit(x float64, rng [2]float64) float64 {
	if x < rng[0] {
		if x < rng[0]*0.8 {
			return -1
		}
		return math.Abs(rng[0]-x) / (rng[0] * 0.2)
	}
	if x > rng[1] {
		if x > rng[1]*1.2 {
			return -1
		}
		return math.Abs(x-rng[1]) / (rng[1] * 0.2)
	}
	return 1
}
