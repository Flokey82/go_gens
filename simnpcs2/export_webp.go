package simnpcs2

import (
	"image"
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
	// Create a colored image of the given width and height.
	img := image.NewNRGBA(image.Rect(0, 0, m.Width, m.Height))
	m.renderFrame(img)
	if err := m.webpExport.anim.AddFrame(img, m.webpExport.timeline, m.webpExport.config); err != nil {
		return err
	}
	m.webpExport.timeline += m.webpExport.timestep
	return nil
}
