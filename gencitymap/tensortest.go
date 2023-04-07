package gencitymap

import (
	"image"
	"image/color"

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

type TotalTensorThing struct {
	tensorField *TensorField
	streamline  *StreamlineGenerator
}

func (t *TotalTensorThing) ExportToPNG(filename string) error {
	sg := t.streamline
	img := image.NewRGBA(image.Rect(0, 0, int(sg.worldDimensions.X), int(sg.worldDimensions.Y)))

	// Fill the background with black.
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	// New draw2d graphics context on an RGBA image.
	gc := draw2dimg.NewGraphicContext(img)

	gen := t.streamline
	g := NewGraph(gen.allStreamlinesSimple, gen.params.Dstep, false)
	f := NewPolygonFinder(g.Nodes, PolygonParams{
		MaxLength:      40,
		MinArea:        30,
		ShrinkSpacing:  1,
		ChanceNoDivide: 0.01,
	}, t.tensorField)
	//f.findPolygons()
	f.Shrink(false)
	f.Divide(false)

	origin := gen.origin

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
		drawPolygon(p, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 0, 0, 255})
	}
	for _, p := range f.ShrunkPolygons {
		drawPolygon(p, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 0, 0, 255})
	}
	for _, p := range f.DividedPolygons {
		drawPolygon(p, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 0, 255})
	}

	// Draw minor streamlines.
	gc.SetStrokeColor(color.RGBA{255, 255, 255, 255})
	gc.SetLineWidth(4.0)
	for _, v := range sg.streamlinesMajor {
		// Draw a path.
		gc.MoveTo(v[0].X-sg.origin.X, v[0].Y-sg.origin.Y)
		for _, p := range v[1:] {
			gc.LineTo(p.X-sg.origin.X, p.Y-sg.origin.Y)
		}
		gc.Stroke()
	}

	gc.SetStrokeColor(color.RGBA{255, 255, 255, 255})
	gc.SetLineWidth(2.0)
	for _, v := range sg.streamlinesMinor {
		// Draw a path.
		gc.BeginPath()
		gc.MoveTo(v[0].X-sg.origin.X, v[0].Y-sg.origin.Y)
		for _, p := range v[1:] {
			gc.LineTo(p.X-sg.origin.X, p.Y-sg.origin.Y)
		}
		gc.Stroke()
	}

	// Save to file.
	return draw2dimg.SaveToPngFile(filename, img)
}

func TensorTest() (*TotalTensorThing, error) {
	// Set up some tensor stuff.
	tt := &TotalTensorThing{
		tensorField: NewTensorField(DefaultNoiseParams),
	}
	// Add some basis fields.
	tt.tensorField.AddGrid(vectors.Vec2{X: 800, Y: 600}, 400, 20.5, 45)
	tt.tensorField.AddGrid(vectors.Vec2{X: 1200, Y: 1000}, 400, 20.5, -45)
	tt.tensorField.AddRadial(vectors.Vec2{X: 1400, Y: 1400}, 300, 20.5)
	tt.tensorField.AddRadial(vectors.Vec2{X: 1200, Y: 1200}, 150, 20.5)
	tt.tensorField.addField(NewOvalField(vectors.Vec2{X: 1600, Y: 1600}, 500, 20.5, -25, 1.5))
	// Add a random sea.
	// tensorField.sea = append(tensorField.sea, vectors.Vec2{X: 500, Y: 500}, vectors.Vec2{X: 1500, Y: 1500}, vectors.Vec2{X: 1100, Y: 1100}, vectors.Vec2{X: 100, Y: 100})
	// tensorField.river = append(tensorField.river, vectors.Vec2{X: 500, Y: 500}, vectors.Vec2{X: 1500, Y: 1500}, vectors.Vec2{X: 1100, Y: 1100}, vectors.Vec2{X: 100, Y: 100})
	// tensorField.river = append(tensorField.river, vectors.Vec2{X: 100, Y: 100}, vectors.Vec2{X: 500, Y: 100}, vectors.Vec2{X: 500, Y: 400}, vectors.Vec2{X: 100, Y: 400})

	streamlineParams := DefaultStreamlineParams
	origin := vectors.Vec2{X: 500, Y: 500}
	dimensions := vectors.Vec2{X: 2000, Y: 2000}
	gen, err := NewStreamlineGenerator(1234, NewRK4Integrator(tt.tensorField, streamlineParams), origin, dimensions, streamlineParams)
	if err != nil {
		return nil, err
	}
	gen.createAllStreamlines(false)
	tt.streamline = gen
	return tt, nil
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
