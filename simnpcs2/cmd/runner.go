package main

import (
	"bufio"
	"os"
	"time"

	"github.com/Flokey82/gameloop"
	"github.com/Flokey82/go_gens/simnpcs2"
)

func main() {
	w := simnpcs2.NewWorld(100, 100, 1234)
	w.AddBeing()
	w.AddBeing()
	w.AddBeing()

	w.AddItem()
	w.AddItem()
	w.AddItem()

	gl := gameloop.New(time.Second/60, func(f float64) {
		w.Update(f)
	})
	gl.Start()

	// Wait for user to press enter.
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	if err := w.ExportWebp("test.webp"); err != nil {
		panic(err)
	}
	if err := w.ExportGif("test.gif"); err != nil {
		panic(err)
	}
}
