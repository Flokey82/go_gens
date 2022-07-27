package main

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/Flokey82/go_gens/genmap2derosion"
)

func main() {

	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	w := genmap2derosion.NewWorld(genmap2derosion.DefaultParams)
	if err := w.ExportGif("anim.gif"); err != nil {
		log.Fatal(err)
	}

	if err := w.ExportOBJ("tmp.obj"); err != nil {
		log.Fatal(err)
	}
}
