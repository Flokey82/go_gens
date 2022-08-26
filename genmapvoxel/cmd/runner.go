package main

import (
	"log"

	"github.com/Flokey82/go_gens/genmapvoxel"
)

func main() {
	w := genmapvoxel.New(32, 32, 32, 12345)
	if err := w.ExportOBJ("tmp.obj"); err != nil {
		log.Println(err)
	}
}
