package main

import (
	"github.com/Flokey82/go_gens/genlsystem"
)

func main() {
	genlsystem.ExportToPNG("bintree.png", genlsystem.BinTree(8))
	genlsystem.ExportToPNG("plant.png", genlsystem.Plant(7))
	genlsystem.ExportToPNG("tree.png", genlsystem.Tree(9))
	genlsystem.ExportToPNG("hilbert.png", genlsystem.Hilbert(5))
	genlsystem.Hilbert3d("out.obj", 3)
	genlsystem.Plant3d("plant.obj", 4)
	genlsystem.Pyramid3d("pyramid.obj", 5)
}
