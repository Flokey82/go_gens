package gencellular

import (
	"image"
	"image/color"
	"image/gif"
	"os"
)

// Animator is a cell culture that can be animated and exported to a GIF.
type Animator struct {
	*Culture
	images  []*image.Paletted // Generated frame used to construct the GIF.
	palette []color.Color     // Default color palette.
	delays  []int             // Delay for each individual frame (0 for now).
}

// NewAnimator returns a new cell culture animator with the given height and width.
func NewAnimator(height, width int) *Animator {
	return NewAnimatorCustom(height, width, SeedDefault, EvalDefault)
}

// NewAnimatorCustom allows for a custom seed init and state eval function.
func NewAnimatorCustom(height, width int, sf SeedFunc, ef EvalFunc) *Animator {
	c := NewCustom(height, width, sf, ef)
	return &Animator{
		Culture: c,
		palette: []color.Color{
			color.RGBA{0x00, 0x00, 0x00, 0xff}, color.RGBA{0x00, 0x00, 0xff, 0xff},
			color.RGBA{0x00, 0xff, 0x00, 0xff}, color.RGBA{0x00, 0xff, 0xff, 0xff},
			color.RGBA{0xff, 0x00, 0x00, 0xff}, color.RGBA{0xff, 0x00, 0xff, 0xff},
			color.RGBA{0xff, 0xff, 0x00, 0xff}, color.RGBA{0xff, 0xff, 0xff, 0xff},
		},
	}
}

// Reset the simulation and clear all frames.
func (c *Animator) Reset() {
	c.Culture.Reset()
	c.images = nil
	c.delays = nil
}

// Tick advances the culture by one tick and stores the current state in the frame.
func (c *Animator) Tick() {
	c.Culture.Tick()

	// Store the previous frame since we have already advanced by one.
	c.storeGifFrame(c.Cells[(c.Generation-1)%2])
}

func (c *Animator) storeGifFrame(cells [][]bool) {
	img := image.NewPaletted(image.Rect(0, 0, c.Width, c.Height), c.palette)
	c.images = append(c.images, img)
	c.delays = append(c.delays, 0)
	for x := 0; x < c.Width; x++ {
		for y := 0; y < c.Height; y++ {
			if cells[x][y] {
				img.Set(x, y, color.RGBA{0xFF, 0x00, 0x00, 255})
			} else {
				img.Set(x, y, color.RGBA{0x00, 0x00, 0x00, 255})
			}
		}
	}
}

// Export all frames to a GIF under the given path.
func (c *Animator) ExportGif(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	if err := gif.EncodeAll(f, &gif.GIF{
		Image: c.images,
		Delay: c.delays,
	}); err != nil {
		return err
	}
	return f.Close()
}
