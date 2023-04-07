package main

import (
	"time"

	"github.com/Flokey82/gameloop"
	"github.com/Flokey82/go_gens/simnpcs2"
)

func main() {
	w := simnpcs2.NewWorld(100, 100)
	w.AddBeing()
	w.AddBeing()
	w.AddBeing()

	gl := gameloop.New(time.Second/60, func(f float64) {
		w.Update(f)
	})

	gl.Start()
	select {}
}
