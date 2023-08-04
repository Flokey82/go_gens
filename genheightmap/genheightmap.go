package genheightmap

import (
	"math"
	"math/rand"

	"github.com/Flokey82/go_gens/vectors"

	opensimplex "github.com/ojrac/opensimplex-go"
)

// Terrain is an interface for a heightmap.
// ... I can't remember if I use this somewhere.
type Terrain interface {
	//ApplyGen(f GenFunc)
	MinMax() (float64, float64)
}

type GenFunc func(x, y float64) float64

// GenSlope returns a generator function that produces a slope in the direction
// of the given vector.
func GenSlope(direction vectors.Vec2) GenFunc {
	return func(x, y float64) float64 {
		return x*direction.X + y*direction.Y
	}
}

// GenFissure returns a generator function that produces a fissure between
// two points on the heightmap, with the given increase in height at the lips
// and the drop in elevation at the base of the fissure.
//
// NOTE: We use a biased random walk to generate the fissure path using the
// given points as the starting and ending points of the walk, and the given
// number of steps as the number of steps in the walk, as well as the amplitude
// of the walk and the width of the fissure.
func GenFissure(p1, p2 vectors.Vec2, steps int, lip, drop, amplitude, width float64) GenFunc {
	// Generate the path of the fissure between the two points.
	// NOTE: We use a biased random walk to generate the fissure turning points using the
	// given points as the starting and ending points of the walk, and the given
	// number of steps as the number of steps in the walk, as well as the amplitude
	// of the walk.
	var path []vectors.Vec2
	for i := 0; i < steps; i++ {
		if i == 0 {
			path = append(path, p1)
			continue
		}
		if i == steps-1 {
			path = append(path, p2)
			continue
		}

		vec12 := p2.Sub(p1)

		// Get the expected point along the path at this step.
		exp := p1.Add(vec12.Mul(float64(i) / float64(steps)))

		// Add a random normal vector with a random fraction of the amplitude (ranging from -0.5 to 0.5 of the amplitude).
		path = append(path, exp.Add(vectors.RandomVec2(1).Mul(amplitude*(rand.Float64()-0.5))))
	}

	// Calculate the total distance between the two points.
	totalDist := p2.Sub(p1).Len()

	// Generate the heightmap.
	return func(x, y float64) float64 {
		p := vectors.Vec2{X: x, Y: y}

		// Calculate the distance to the closest end point on the path.
		var distToEnd float64
		distToEnd = p1.Sub(p).Len()
		if d2 := p2.Sub(p).Len(); d2 < distToEnd {
			distToEnd = d2
		}

		// Calculate fissure decay at this point.
		dec := (2 * distToEnd / totalDist)
		// Calculate fissure width at this point.
		w := width * dec
		// Calculate the lip height at this point.
		l := lip * dec
		// Calculate the drop height at this point.
		d := drop * dec

		// Get the shortest distance from the point to the path segment.
		var minDist float64
		for i := 0; i < len(path)-1; i++ {
			// Get the distance from the point to the path segment.
			dist := vectors.Segment{
				Start: path[i],
				End:   path[i+1],
			}.DistanceToPoint(p)

			// Update the shortest distance.
			if i == 0 || dist < minDist {
				minDist = dist
			}
		}

		// If we are within the width of the fissure, we interpolate the height
		// between the lip and the drop.
		if minDist < w {
			return l - (w-minDist)*(l-d)/w
		}

		// Otherwise, we have an smooth decay of the height.
		return l / (minDist / w)
	}
}

// GenMountainRange returns a generator function that produces a mountain range
// between the two given points, with the given number of peaks (or steps), the
// given radius of the mountains, and the given amplitude of the variation in the
// mountain range and a maximum height of the mountains
//
// NOTE: We use a biased random walk to generate the mountain range using the
// given points as the starting and ending points of the walk, and the given
// number of steps as the number of steps in the walk, as well as the amplitude
// of the walk and the radius of the mountains.
func GenMountainRange(p1, p2 vectors.Vec2, steps int, radius, amplitude, maxHeight float64) GenFunc {
	// Generate the path of the mountain range between the two points.
	// NOTE: We use a biased random walk to generate the mountain range turning points using the
	// given points as the starting and ending points of the walk, and the given
	// number of steps as the number of steps in the walk, as well as the amplitude
	// of the walk.
	var path []vectors.Vec2
	for i := 0; i < steps; i++ {
		if i == 0 {
			path = append(path, p1)
			continue
		}
		if i == steps-1 {
			path = append(path, p2)
			continue
		}

		vec12 := p2.Sub(p1)

		// Get the expected point along the path at this step.
		exp := p1.Add(vec12.Mul(float64(i) / float64(steps)))

		// Add a random normal vector with a random fraction of the amplitude (ranging from -0.5 to 0.5 of the amplitude).
		path = append(path, exp.Add(vectors.RandomVec2(1).Mul(amplitude*(rand.Float64()-0.5))))
	}

	// Calculate the total distance between the two points.
	totalDist := p2.Sub(p1).Len()

	// Generate the heightmap.
	return func(x, y float64) float64 {
		p := vectors.Vec2{X: x, Y: y}

		// Calculate the distance to the closest end point on the path.
		distToEnd := p1.Sub(p).Len()
		if d2 := p2.Sub(p).Len(); d2 < distToEnd {
			distToEnd = d2
		}

		// Get the shortest distance from the point to the closest path segment.
		var (
			closestSegDist float64         // Distance from the point to the closest path segment.
			closestSeg     vectors.Segment // Closest path segment.
		)
		for i := 0; i < len(path)-1; i++ {
			// Get the distance from the point to the path segment.
			seg := vectors.Segment{
				Start: path[i],
				End:   path[i+1],
			}
			dist := seg.DistanceToPoint(p)

			// Update the shortest distance.
			if i == 0 || dist < closestSegDist {
				closestSegDist = dist
				closestSeg = seg
			}
		}
		// Length of the closest path segment.
		closestSegLen := closestSeg.Len()

		// Distance from the closest point on the segment to the center of the segment.
		poinOnLine := closestSeg.ClosestPoint(p)
		distToSegCenter := closestSeg.Midpoint().Sub(poinOnLine).Len()

		// Calculate the height of the mountain at this point, with some additional
		// decay based on the distance to the center of the segment.
		//
		// Depending on the distance to the peak, we sag a bit in the middle of the mountain
		// range to make it look more like a mountain range.
		//   _   __  _
		//  / \_/  \ _  <- 0.4 * sag
		// /        \_  <- 0.6
		sag := 1 - math.Pow(distToSegCenter/(closestSegLen/2), 2)

		// Calculate the height of the mountain at this point.
		h := radius * (1.0 - 0.4*sag)

		// Calculate mountain range decay at this point.
		// The closer a point is to the end of the mountain range, the lower the height.
		dec := (2 * distToEnd / totalDist)

		// Calculate the width of the mountain at this point.
		rDec := radius * dec

		// Calculate the height of the mountain at this point.
		hDec := h * dec

		// If we are within the radius of the mountain, we interpolate the height
		// given the expected height of the mountain at this point.
		if closestSegDist < rDec {
			return hDec * math.Pow((rDec-closestSegDist)/rDec, 2)
		}
		return 0
	}
}

// GenCrater modifies the height values in a circular region to create a crater with a lip.
func GenCrater(center vectors.Vec2, diameter, lip, depth float64) GenFunc {
	radius := diameter / 2

	return func(x, y float64) float64 {
		// Calculate the distance from the center of the crater
		distance := math.Sqrt((x-center.X)*(x-center.X) + (y-center.Y)*(y-center.Y))

		// If the distance is less than the radius of the crater, we are inside the crater.
		if distance < radius {
			// Lower the height value based on distance from the crater center
			// forming a nice bowl shape.
			return lip - depth*(1-(distance/radius)*(distance/radius))
		}
		// Decays the lip height based on distance from the crater lip (the edge of the crater)
		// The farther from the lip, the shallower the lip height
		// The lip height determines the maximum height the crater can have.
		return lip / (distance / float64(radius))
	}
}

// GenCone returns a generator function for a cone at the center of the heightmap.
// TODO: Allow the user to specify the center of the cone.
func GenCone(slope float64) GenFunc {
	return func(x, y float64) float64 {
		return math.Sqrt(x*x+y*y) * slope
	}
}

// GenVolCone returns a generator function for a volcanic cone
// at the center of the heightmap.
// TODO: Allow the user to specify the center of the cone.
func GenVolCone(slope float64) GenFunc {
	return func(x, y float64) float64 {
		dist := math.Sqrt(x*x + y*y)
		if dist < slope/4 {
			return math.Pow(dist/(slope/4), 2) * slope / 2
		}
		// Let the cone decay to the edge of the heightmap.
		return slope / math.Pow(dist/(slope/4), 2)
	}
}

// GenMountains returns a generator function that will return the height of a
// point on the heightmap given the point's coordinates, which will produce a
// number of mountains.
// TODO: The seed should be passed into the function as parameter.
//
// 'maxX', 'maxY' are the dimensions of the heightmap.
// 'n' is the number of mountains.
// 'r' is the radius of the mountains.
func GenMountains(maxX, maxY float64, n int, r float64) GenFunc {
	rand.Seed(1234)
	var mounts [][2]float64
	for i := 0; i < n; i++ {
		mounts = append(mounts, [2]float64{maxX * (rand.Float64() - 0.5), maxY * (rand.Float64() - 0.5)})
	}
	return func(x, y float64) float64 {
		var val float64
		for j := 0; j < n; j++ {
			m := mounts[j]
			val += math.Pow(math.Exp(-((x-m[0])*(x-m[0])+(y-m[1])*(y-m[1]))/(2*r*r)), 2)
		}
		return val
	}
}

// GenNoise returns a function that returns the noise/height value of a given point
// on the heightmap. Not sure what the slope parameter was supposed to do.
func GenNoise(seed int64, slope float64) GenFunc {
	perlin := opensimplex.New(seed)

	mult := 15.0
	pow := 1.0
	return func(x, y float64) float64 {
		x *= mult
		y *= mult
		e := 1 * math.Abs(perlin.Eval2(x, y))
		e += 0.5 * math.Abs(perlin.Eval2(x*2, y*2))
		e += 0.25 * perlin.Eval2(x*4, y*4)
		e /= (1 + 0.5 + 0.25)
		return math.Pow(e, pow)
	}
}

// CalcMean calculates the mean of a slice of floats.
func CalcMean(nums []float64) float64 {
	total := 0.0
	for _, v := range nums {
		total += v
	}
	return total / float64(len(nums))
}

// MinMax returns the min and max values of the heightmap.
func MinMax(hm []float64) (float64, float64) {
	if len(hm) == 0 {
		return 0, 0
	}
	min := hm[0]
	max := hm[0]
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

// Modify is a function that modifies a value in a heightmap.
type Modify func(val float64) float64

// ModNormalize normalizes the heightmap to the range [0, 1] given
// the min and max values (the range of heightmap values).
func ModNormalize(min, max float64) Modify {
	return func(val float64) float64 {
		return (val - min) / (max - min)
	}
}

// ModPeaky returns the function applied to a point on a heightmap
// in order to exaggerate the peaks of the map.
func ModPeaky() Modify {
	return math.Sqrt
}

// ModSeaLevel shifts the origin point to the sea level, resulting in
// all points below sea level being negative.
func ModSeaLevel(min, max, q float64) Modify {
	delta := min + (max-min)*0.1
	//delta := quantile(h, q)
	return func(val float64) float64 {
		return val - delta
	}
}

// ModifyWithIndex is a function that modifies a value in a heightmap given
// its index and current value.
type ModifyWithIndex func(idx int, val float64) float64

// GetNeighbors returns all neighbor indices of an index on the heightmap.
type GetNeighbors func(idx int) []int

// GetHeight returns the height of a point on the heightmap given its index.
type GetHeight func(idx int) float64

// ModRelax applies a relaxation algorithm to the heightmap.
func ModRelax(n GetNeighbors, h GetHeight) ModifyWithIndex {
	return func(idx int, val float64) float64 {
		vals := []float64{val}
		for _, nb := range n(idx) {
			vals = append(vals, h(nb))
		}
		return CalcMean(vals)
	}
}
