package gencitymap

import (
	"github.com/Flokey82/go_gens/vectors"
)

var DefaultNoiseParams = &NoiseParams{
	Seed:             0,
	globalNoise:      true,
	noiseSizePark:    2,
	noiseAnglePark:   9,
	noiseSizeGlobal:  3,
	noiseAngleGlobal: 2,
}

var DefaultStreamlineParams = &StreamlineParams{
	Dsep:              20,
	Dtest:             15,
	Dstep:             1,
	Dcirclejoin:       5,
	Dlookahead:        40,
	Joinangle:         0.1,
	PathIterations:    2000,
	SeedTries:         30,
	SimplifyTolerance: 0.5,
	CollideEarly:      0.01,
}

func Tensstuff() {
	// Set up some tensor stuff.
	tensorField := NewTensorField(DefaultNoiseParams)

	// Add some basis fields.
	tensorField.AddGrid(vectors.Vec2{X: 800, Y: 600}, 400, 20.5, 45)
	tensorField.AddGrid(vectors.Vec2{X: 1200, Y: 1000}, 400, 20.5, -45)
	tensorField.AddRadial(vectors.Vec2{X: 1400, Y: 1400}, 300, 20.5, 0.5)
	tensorField.AddRadial(vectors.Vec2{X: 1200, Y: 1200}, 150, 20.5, 0.5)

	streamlineParams := DefaultStreamlineParams
	origin := vectors.Vec2{X: 500, Y: 500}
	dimensions := vectors.Vec2{X: 2000, Y: 2000}
	gen := NewStreamlineGenerator(1234, NewRK4Integrator(tensorField, streamlineParams), origin, dimensions, streamlineParams)
	gen.createAllStreamlines(false)

	// Create a png image.
	gen.ExportToPNG("test2.png")

	// Create an svg image.
	gen.ExportToSVG("test2.svg")
}

func convToPairs(vectors []vectors.Vec2, origin vectors.Vec2) ([]int, []int) {
	x := make([]int, len(vectors))
	y := make([]int, len(vectors))
	for i, v := range vectors {
		x[i] = int(v.X - origin.X)
		y[i] = int(v.Y - origin.Y)
	}
	return x, y
}
