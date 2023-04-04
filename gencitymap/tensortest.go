package gencitymap

import (
	"image"
	"image/color"
	"log"

	"github.com/Flokey82/go_gens/vectors"
	"github.com/llgcode/draw2d/draw2dimg"
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
	SimplifyTolerance: 10,
	CollideEarly:      0.01,
}

func TensorTest() (*StreamlineGenerator, error) {
	// Set up some tensor stuff.
	tensorField := NewTensorField(DefaultNoiseParams)
	// Add some basis fields.
	tensorField.AddGrid(vectors.Vec2{X: 800, Y: 600}, 400, 20.5, 45)
	tensorField.AddGrid(vectors.Vec2{X: 1200, Y: 1000}, 400, 20.5, -45)
	tensorField.AddRadial(vectors.Vec2{X: 1400, Y: 1400}, 300, 20.5)
	tensorField.AddRadial(vectors.Vec2{X: 1200, Y: 1200}, 150, 20.5)
	tensorField.addField(NewOvalField(vectors.Vec2{X: 1600, Y: 1600}, 500, 20.5, -25, 1.5))
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

	for _, s := range gen.allStreamlinesSimple {
		log.Println("streamline", s)
	}

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
	f := NewPolygonFinder(g.nodes, PolygonParams{
		MaxLength:      25,
		MinArea:        30,
		ShrinkSpacing:  4,
		ChanceNoDivide: 0.01,
	}, tensorField)
	//f.findPolygons()
	f.Shrink(false)
	f.Divide(false)

	// New PNG image.
	img := image.NewRGBA(image.Rect(0, 0, int(dimensions.X), int(dimensions.Y)))

	// Fill the background with black.
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}
	// New draw2d graphics context on an RGBA image.
	gc := draw2dimg.NewGraphicContext(img)

	drawPolygon := func(p []vectors.Vec2, colFill, colStroke color.RGBA) {
		gc.SetFillColor(colFill)
		gc.SetStrokeColor(colStroke)
		gc.SetLineWidth(1)
		for i, v := range p {
			if i == 0 {
				gc.MoveTo(v.X-origin.X, v.Y-origin.Y)
			} else {
				gc.LineTo(v.X-origin.X, v.Y-origin.Y)
			}
		}
		gc.Close()
		gc.FillStroke()
	}

	// Draw the polygons.
	for _, p := range f.Polygons {
		log.Println("polygon3", p)
		drawPolygon(p, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 0, 0, 255})
	}
	for _, p := range f.ShrunkPolygons {
		log.Println("polygon2", p)
		drawPolygon(p, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 0, 0, 255})
	}
	for _, p := range f.DividedPolygons {
		log.Println("polygon1", p)
		drawPolygon(p, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 0, 255})
	}

	// write to file
	draw2dimg.SaveToPngFile("test_tensor_polygons.png", img)
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
