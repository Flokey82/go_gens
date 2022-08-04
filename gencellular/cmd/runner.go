package main

import (
	"github.com/Flokey82/go_gens/gencellular"
)

func main() {
	c := gencellular.NewAnimator(128, 128)
	for i := 0; i < 100; i++ {
		c.Tick()
	}
	c.ExportGif("rgb.gif")
}
