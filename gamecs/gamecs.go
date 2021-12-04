package gamecs

import (
	"github.com/Flokey82/go_gens/vectors"
	"image"
	"image/color"
	"image/gif"
	"math/rand"
	"os"
)

type World struct {
	c       []*Character
	images  []*image.Paletted // Generated frame used to construct the GIF.
	palette []color.Color     // Default color palette.
	delays  []int             // Delay for each individual frame (0 for now).
	Width   int
	Height  int
}

func New() *World {
	return &World{
		palette: []color.Color{
			color.RGBA{0x00, 0x00, 0x00, 0xff}, color.RGBA{0x00, 0x00, 0xff, 0xff},
			color.RGBA{0x00, 0xff, 0x00, 0xff}, color.RGBA{0x00, 0xff, 0xff, 0xff},
			color.RGBA{0xff, 0x00, 0x00, 0xff}, color.RGBA{0xff, 0x00, 0xff, 0xff},
			color.RGBA{0xff, 0xff, 0x00, 0xff}, color.RGBA{0xff, 0xff, 0xff, 0xff},
		},
		Width:  128,
		Height: 128,
	}
}

func (w *World) Add(c *Character) {
	w.c = append(w.c, c)
}

func (w *World) Update(delta float64) {
	for _, c := range w.c {
		c.Update(delta)
	}
	w.storeGifFrame()
}

func (w *World) storeGifFrame() {
	img := image.NewPaletted(image.Rect(0, 0, w.Width, w.Height), w.palette)
	w.images = append(w.images, img)
	w.delays = append(w.delays, 0)
	for _, c := range w.c {
		img.Set(int(c.Pos.X), int(c.Pos.Y), color.RGBA{0xFF, 0x00, 0x00, 255})
		img.Set(int(c.Target.X), int(c.Target.Y), color.RGBA{0x00, 0xFF, 0x00, 255})
		if c.Waypoints == nil {
			continue
		}
		for _, wp := range c.Waypoints[c.WaypointCurrent:] {
			img.Set(int(wp.X), int(wp.Y), color.RGBA{0xFF, 0xFF, 0x00, 255})
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

type Character struct {
	w *World
	CMovable
	*CAi
}

func (w *World) NewChar() *Character {
	c := &Character{
		w: w,
		CMovable: CMovable{
			Pos: vectors.Vec2{
				X: float64(rand.Intn(w.Height)),
				Y: float64(rand.Intn(w.Width)),
			},
		},
		CAi: newCAi(w),
	}
	w.Add(c)
	return c
}

func (c *Character) Update(delta float64) {
	c.CAi.Update(&c.CMovable, delta)
	c.CMovable.Update(delta)
}
