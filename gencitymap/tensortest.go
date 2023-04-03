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

func TensorTest() (*StreamlineGenerator, error) {
	// Set up some tensor stuff.
	tensorField := NewTensorField(DefaultNoiseParams)
	// Add some basis fields.
	tensorField.AddGrid(vectors.Vec2{X: 800, Y: 600}, 400, 20.5, 45)
	tensorField.AddGrid(vectors.Vec2{X: 1200, Y: 1000}, 400, 20.5, -45)
	tensorField.AddRadial(vectors.Vec2{X: 1400, Y: 1400}, 300, 20.5, 0.5)
	tensorField.AddRadial(vectors.Vec2{X: 1200, Y: 1200}, 150, 20.5, 0.5)

	// Add a random sea.
	// tensorField.sea = append(tensorField.sea, vectors.Vec2{X: 500, Y: 500}, vectors.Vec2{X: 1500, Y: 1500}, vectors.Vec2{X: 1100, Y: 1100}, vectors.Vec2{X: 100, Y: 100})
	// tensorField.river = append(tensorField.river, vectors.Vec2{X: 500, Y: 500}, vectors.Vec2{X: 1500, Y: 1500}, vectors.Vec2{X: 1100, Y: 1100}, vectors.Vec2{X: 100, Y: 100})
	// tensorField.river = append(tensorField.river, vectors.Vec2{X: 100, Y: 100}, vectors.Vec2{X: 500, Y: 100}, vectors.Vec2{X: 500, Y: 400}, vectors.Vec2{X: 100, Y: 400})

	streamlineParams := DefaultStreamlineParams
	origin := vectors.Vec2{X: 500, Y: 500}
	dimensions := vectors.Vec2{X: 2000, Y: 2000}
	gen, err := NewStreamlineGenerator(1234, NewRK4Integrator(tensorField, streamlineParams), origin, dimensions, streamlineParams)
	if err != nil {
		return nil, err
	}
	gen.createAllStreamlines(false)

	/*
		g := NewGraph(gen.allStreamlinesSimple, streamlineParams.Dstep, false)
		for _, in := range g.intersections {
			log.Println("intersection", in)
		}
		for _, n := range g.nodes {
			if len(n.neighbors) == 0 {
				println("dangling node")
			} else {
				println("ok node with neighbors")
			}
		}
	*/
	return gen, nil
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
