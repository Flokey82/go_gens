// Package genmap2derosion is based on "Simple Hydrology"
// See: https://nickmcd.me/2020/04/15/procedural-hydrology/
package genmap2derosion

import (
	"fmt"
	"github.com/Flokey82/go_gens/genheightmap"
	"github.com/Flokey82/go_gens/vectors"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"log"
	"math/rand"
	"os"
)

const worldsize = 256

type World struct {
	images  []*image.Paletted // Generated frame used to construct the GIF.
	palette []color.Color     // Default color palette.
	delays  []int             // Delay for each individual frame (0 for now).
	seed    int64
	r       *rand.Rand
	dim     vectors.IVec2 // Size of the heightmap array
	fdim    vectors.Vec2

	drainage    float64   // Drainage factor from pools
	scale       float64   // "Physical" Height scaling of the map
	heightmap   []float64 // Flat Array containing height information
	sediment    []float64 // Sediment storage (stores information on eroded / deposited sediment)
	waterpath   []float64 // Water Path Storage (Rivers)
	waterpool   []float64 // Water Pool Storage (Lakes / Ponds)
	waterdrains []int     // Water Pool drainage points

	// Flux related information (experimental)
	fluxwaterpool []float64 // (TEMP Flux) Water Pool Storage (Lakes / Ponds)

	// Trees
	// trees []Plant
	// plantdensity [worldsize * worldsize]float64 //Density for Plants

	// Sediment erosion resistance
	// sedimentToughness [worldsize * worldsize]float64 // Sediment toughness storage
}

func NewWorld() *World {
	dimX := int64(worldsize)
	dimY := int64(worldsize)
	idxSize := dimX * dimY
	w := &World{
		drainage:      0.01,
		scale:         40.0,
		dim:           vectors.NewIVec2(dimX, dimY),
		fdim:          vectors.NewVec2(float64(dimX), float64(dimY)),
		heightmap:     make([]float64, idxSize),
		sediment:      make([]float64, idxSize),
		waterpath:     make([]float64, idxSize),
		waterpool:     make([]float64, idxSize),
		waterdrains:   make([]int, idxSize),
		fluxwaterpool: make([]float64, idxSize),
	}

	for i := 0; i <= 255; i++ {
		w.palette = append(w.palette, color.RGBA{uint8(i), uint8(i), uint8(i), 0xff})
	}

	// Seed the Random Generator
	if w.seed == 0 {
		w.seed = 12356
		w.r = rand.New(rand.NewSource(w.seed))
	}

	// Initialize all water drains to -1 (unset)
	for i := range w.waterdrains {
		w.waterdrains[i] = -1
	}

	// Generate basic heightmap.
	w.generate()
	w.generateClimate()

	fmt.Println("generated")

	// Erode a few times.
	for j := 0; j < 5; j++ {
		w.erode(50, 300*5)
		//w2 = w.generateClimate()
		//w2.generate(w.heightmap)
		//w.erodeRain(1, w2.climate.AvgRainMap)

		// Export hydrology data generated by the original algorithm.
		w.export(fmt.Sprintf("b_image%d.png", j), w.heightmap)
		w.export(fmt.Sprintf("b_image%d_wp.png", j), w.waterpath)
		w.export(fmt.Sprintf("b_image%d_wpo.png", j), w.waterpool)
		w.export(fmt.Sprintf("b_image%d_sed.png", j), w.sediment)
		w.exportCombined(fmt.Sprintf("b_image%d_combo.png", j), w.heightmap, w.waterpath, w.waterpool)

		// Optional stuff...

		// Call downhill, which will generate flux information.
		// for i, dhs := range w.getDownhill() {
		// 	if len(dhs) == 0 {
		// 		fmt.Println(fmt.Sprintf("sink: %d", i))
		// 	}
		// }

		// Export alternative hydrology data generated through flux.
		// w.export(fmt.Sprintf("b_image%d_flux_raw.png", j), w.getRawFlux())
		// w.export(fmt.Sprintf("b_image%d_flux.png", j), w.getFlux())
		// w.export(fmt.Sprintf("b_image%d_flux_wpo.png", j), w.fluxwaterpool[:])
	}

	if err := w.ExportGif("anim.gif"); err != nil {
		panic(err)
	}

	if err := w.ExportOBJ("tmp.obj"); err != nil {
		panic(err)
	}
	return w
}

func (w *World) storeGifFrame(h, waterPath, waterPool []float64) {
	width, height := int(w.dim.X), int(w.dim.Y)
	img := image.NewPaletted(image.Rect(0, 0, width, height), w.palette)
	w.images = append(w.images, img)
	w.delays = append(w.delays, 0)
	min, max := genheightmap.MinMax(h)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := int64(x)*int64(w.dim.Y) + int64(y)
			// NOTE: By averaging things out, stuff might look weird (sediment for example)
			val := (h[i] - min) / (max - min)
			colVal := uint8(val * 255)
			col := w.palette[colVal]
			img.Set(x, y, col)
		}
	}
}

// Export all frames to a GIF under the given path.
func (w *World) ExportGif(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	if err := gif.EncodeAll(f, &gif.GIF{
		Image: w.images,
		Delay: w.delays,
	}); err != nil {
		return err
	}
	return f.Close()
}

func (w *World) export(name string, h []float64) {
	width, height := int(w.dim.X), int(w.dim.Y)

	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	min, max := genheightmap.MinMax(h)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := int64(x)*int64(w.dim.Y) + int64(y)
			// NOTE: By averaging things out, stuff might look weird (sediment for example)
			val := (h[i] - min) / (max - min)
			colVal := uint8(val * 255)
			img.Set(x, y, color.NRGBA{
				R: colVal,
				G: colVal,
				B: colVal,
				A: 255,
			})
		}
	}

	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func (w *World) exportCombined(name string, heightMap, waterPath, waterPool []float64) {
	width, height := int(w.dim.X), int(w.dim.Y)

	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	min, max := genheightmap.MinMax(heightMap)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := int64(x)*int64(w.dim.Y) + int64(y)
			val := (heightMap[i] - min) / (max - min)
			colVal := uint8(val * 255)
			col := color.NRGBA{
				R: colVal,
				G: colVal,
				B: colVal,
				A: 255,
			}

			// TODO: Shade by water level.
			if waterPath[i] > 0.01 || waterPool[i] > 0.0 {
				col = color.NRGBA{67, 162, 202, 255}
			}
			img.Set(x, y, col)
		}
	}

	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
