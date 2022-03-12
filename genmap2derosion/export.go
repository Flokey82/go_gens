package genmap2derosion

import (
	"bufio"
	"fmt"
	"github.com/Flokey82/go_gens/genheightmap"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"log"
	"os"
)

func (w *World) ExportOBJ(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	wr := bufio.NewWriter(f)
	defer wr.Flush()

	// Write the vertex index.
	sx := w.params.Size.X
	sy := w.params.Size.Y
	for i, h := range w.heightmap {
		x := float64(i/int(sy)) / float64(sy)
		y := float64(i%int(sy)) / float64(sx)
		wr.WriteString(fmt.Sprintf("v %f %f %f \n", x, h*0.2, y))
	}

	// Write the triangles.
	for x := 0; x < int(sx)-1; x++ {
		for y := 0; y < int(sy-1); y++ {
			i1 := x*int(sy) + y
			i2 := i1 + 1
			i3 := i2 + int(sy)
			i4 := i1 + int(sy)
			wr.WriteString(fmt.Sprintf("f %d %d %d \n", i1+1, i2+1, i3+1))
			wr.WriteString(fmt.Sprintf("f %d %d %d \n", i4+1, i1+1, i3+1))
		}
	}
	return nil
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

func (w *World) storeGifFrame(h, waterPath, waterPool []float64) {
	width, height := int(w.params.Size.X), int(w.params.Size.Y)
	img := image.NewPaletted(image.Rect(0, 0, width, height), w.palette)
	w.images = append(w.images, img)
	w.delays = append(w.delays, 0)
	min, max := genheightmap.MinMax(h)
	sy := w.params.Size.Y
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := int64(x)*sy + int64(y)
			// NOTE: By averaging things out, stuff might look weird (sediment for example)
			val := (h[i] - min) / (max - min)
			colVal := uint8(val * 255)
			col := w.palette[colVal]
			img.Set(x, y, col)
		}
	}
}

func (w *World) ExportPng(name string, h []float64) {
	width, height := int(w.params.Size.X), int(w.params.Size.Y)

	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	min, max := genheightmap.MinMax(h)
	sy := w.params.Size.Y
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := int64(x)*sy + int64(y)
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
	width, height := int(w.params.Size.X), int(w.params.Size.Y)

	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	min, max := genheightmap.MinMax(heightMap)
	sy := w.params.Size.Y
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := int64(x)*sy + int64(y)
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
