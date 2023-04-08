package simnpcs2

import (
	"image"
	"image/color"
	"os"

	"github.com/sizeofint/webpanimation"
)

type webpExport struct {
	anim     *webpanimation.WebpAnimation
	config   webpanimation.WebPConfig
	timeline int
	timestep int
}

func newWebPExport(width, height int) *webpExport {
	anim := webpanimation.NewWebpAnimation(width, height, 0)
	anim.WebPAnimEncoderOptions.SetKmin(9)
	anim.WebPAnimEncoderOptions.SetKmax(17)

	config := webpanimation.NewWebpConfig()
	config.SetLossless(1)
	return &webpExport{
		anim:     anim,
		config:   config,
		timeline: 0,
		timestep: 50,
	}
}

func (m *webpExport) ExportWebp(name string) error {
	// Write the final frame.
	m.timeline += m.timestep
	if err := m.anim.AddFrame(nil, m.timeline, m.config); err != nil {
		return err
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}

	// Encode animation and write result bytes in buffer.
	if err = m.anim.Encode(f); err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	m.anim.ReleaseMemory() // TODO: This doesn't really prevent crashes?

	return nil
}

func (m *World) storeWebPFrame() error {
	// Write the current map to the animation.
	if err := m.webpExport.anim.AddFrame(m.getImage(false, true), m.webpExport.timeline, m.webpExport.config); err != nil {
		return err
	}
	m.webpExport.timeline += m.webpExport.timestep
	return nil
}

func (w *World) getImage(drawTerritories, drawSeasonalBiome bool) image.Image {
	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, w.Width, w.Height))

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
	return img
}
