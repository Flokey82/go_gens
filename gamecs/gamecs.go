package gamecs

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math/rand"
	"os"

	"github.com/Flokey82/go_gens/vectors"
)

type World struct {
	images  []*image.Paletted // Generated frame used to construct the GIF.
	palette []color.Color     // Default color palette.
	delays  []int             // Delay for each individual frame (0 for now).
	Width   int
	Height  int
	mgr     *Manager
}

func New() *World {
	w := &World{
		palette: []color.Color{
			color.RGBA{0x00, 0x00, 0x00, 0xff}, color.RGBA{0x00, 0x00, 0xff, 0xff},
			color.RGBA{0x00, 0xff, 0x00, 0xff}, color.RGBA{0x00, 0xff, 0xff, 0xff},
			color.RGBA{0xff, 0x00, 0x00, 0xff}, color.RGBA{0xff, 0x00, 0xff, 0xff},
			color.RGBA{0xff, 0xff, 0x00, 0xff}, color.RGBA{0xff, 0xff, 0xff, 0xff},
		},
		Width:  128,
		Height: 128,
	}
	w.mgr = newManager()
	w.placeFood()
	return w
}

func (w *World) placeFood() {
	itFood := NewItemType("goulash", "food")
	for i := 0; i < 200; i++ {
		w.mgr.RegisterItem(itFood.New(w, vectors.Vec2{
			X: float64(rand.Intn(w.Height)),
			Y: float64(rand.Intn(w.Width)),
		}))
	}
}

func (w *World) Update(delta float64) {
	for _, c := range w.mgr.Entities() {
		if c.Dead() {
			w.mgr.RemoveEntity(c)
			continue
		}
		c.Update(delta)
	}
	w.storeGifFrame()
}

func (w *World) storeGifFrame() {
	img := image.NewPaletted(image.Rect(0, 0, w.Width, w.Height), w.palette)
	w.images = append(w.images, img)
	w.delays = append(w.delays, 0)

	// Draw all entities and their paths.
	for _, c := range w.mgr.Entities() {
		img.Set(int(c.Pos.X), int(c.Pos.Y), color.RGBA{0xFF, 0x00, 0x00, 255})
		img.Set(int(c.Target.X), int(c.Target.Y), color.RGBA{0x00, 0xFF, 0x00, 255})
		if c.Waypoints == nil {
			continue
		}
		for _, wp := range c.Waypoints[c.WaypointCurrent:] {
			img.Set(int(wp.X), int(wp.Y), color.RGBA{0xFF, 0xFF, 0x00, 255})
		}
	}

	// Draw all items that are visible.
	for _, c := range w.mgr.items {
		if c.Location != LocWorld {
			continue
		}
		img.Set(int(c.Pos.X), int(c.Pos.Y), color.RGBA{0xff, 0x00, 0xff, 255})
	}

	// Draw all locations / homes.
	for _, loc := range w.mgr.locations {
		img.Set(int(loc.Pos.X), int(loc.Pos.Y), color.RGBA{0x00, 0x00, 0xff, 255})
	}
}

// Export all frames to a GIF under the given path.
func (w *World) ExportGif(path string) error {
	for _, loc := range w.mgr.locations {
		log.Println(fmt.Sprintf("loc %d in storage", len(loc.CInventory.Slots)))
	}
	for _, loc := range w.mgr.entities {
		log.Println(fmt.Sprintf("age %d in storage", len(loc.CInventory.Slots)))
	}
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
