package main

import (
	"fmt"
	"github.com/Flokey82/go_gens/gamecs"
	"github.com/Flokey82/go_gens/gameloop"
	"time"
)

func main() {
	w := gamecs.New()
	w.NewChar()
	w.NewChar()

	gl := gameloop.New(200*time.Millisecond, w.Update)
	gl.Start()
	fmt.Scanln()
	fmt.Println("done")
	w.ExportGif("rgb.gif")
}
