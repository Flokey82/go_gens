package genheightmap

import (
	"math"
	"math/rand"

	opensimplex "github.com/ojrac/opensimplex-go"
)

type Terrain interface {
	//ApplyGen(f GenFunc)
	MinMax() (float64, float64)
}

type GenFunc func(x, y float64) float64

func GenSlope(direction [2]float64) GenFunc {
	return func(x, y float64) float64 {
		return x*direction[0] + y*direction[1]
	}
}

func GenCone(slope float64) GenFunc {
	return func(x, y float64) float64 {
		return math.Pow(x*x+y*y, 0.5) * slope
	}
}

func GenVolCone(slope float64) GenFunc {
	return func(x, y float64) float64 {
		dist := math.Pow(x*x+y*y, 0.5)
		if dist < 0.1 {
			return -4 * dist * slope
		}
		return dist * slope
	}
}

func GenMountains(maxX, maxY float64, n int, r float64) GenFunc {
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
