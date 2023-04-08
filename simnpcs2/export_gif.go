package simnpcs2

import (
	"image"
	"image/color"
	"image/gif"
	"os"
)

type gifExport struct {
	images  []*image.Paletted // Generated frame used to construct the GIF.
	palette []color.Color     // Default color palette.
	delays  []int             // Delay for each individual frame (0 for now).
}

func newGifExport() *gifExport {
	return &gifExport{
		palette: []color.Color{
			color.RGBA{0x00, 0x00, 0x00, 0xff}, color.RGBA{0x00, 0x00, 0xff, 0xff},
			color.RGBA{0x00, 0xff, 0x00, 0xff}, color.RGBA{0x00, 0xff, 0xff, 0xff},
			color.RGBA{0xff, 0x00, 0x00, 0xff}, color.RGBA{0xff, 0x00, 0xff, 0xff},
			color.RGBA{0xff, 0xff, 0x00, 0xff}, color.RGBA{0xff, 0xff, 0xff, 0xff},
		},
	}
}

// Export all frames to a GIF under the given path.
func (w *gifExport) ExportGif(path string) error {
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

func (w *World) storeGifFrame() {
	img := image.NewPaletted(image.Rect(0, 0, w.Width, w.Height), w.palette)
	w.images = append(w.images, img)
	w.delays = append(w.delays, 0)

	// Draw all entities and their paths.
	for _, c := range w.Beings {
		ai := c.(*AI)
		if ai.Pathfinding.Waypoints != nil {
			for _, wpIdx := range ai.Pathfinding.Waypoints[ai.Pathfinding.WaypointIdx:] {
				wp := w.CellIdxToPos(wpIdx)
				img.Set(int(wp.X), int(wp.Y), color.RGBA{0xFF, 0xFF, 0x00, 255})
			}
		}

		pos := c.Pos()
		img.Set(int(pos.X), int(pos.Y), color.RGBA{0xFF, 0x00, 0x00, 255})

		if dst := ai.Destination; dst != nil {
			img.Set(int(dst.X), int(dst.Y), color.RGBA{0x00, 0xFF, 0x00, 255})
		}
	}

	// Draw all obstacles.
	for i, c := range w.Cells {
		if c {
			pos := w.CellIdxToPos(i)
			img.Set(int(pos.X), int(pos.Y), color.RGBA{0x00, 0x00, 0xFF, 255})
		}
	}

	// Draw all items that are visible.
	for _, c := range w.Items {
		pos := c.Pos()
		img.Set(int(pos.X), int(pos.Y), color.RGBA{0xff, 0x00, 0xff, 255})
	}
}
