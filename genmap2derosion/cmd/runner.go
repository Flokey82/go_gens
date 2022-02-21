package main

import (
	"github.com/Flokey82/go_gens/genmap2derosion"
	"log"
	"os"
	"runtime/pprof"
)

func main() {

	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	genmap2derosion.NewWorld()
}
