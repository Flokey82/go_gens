package main

import (
	"github.com/Flokey82/go_gens/gendungeon"
)

func main() {
	dng := gendungeon.Generate(40, 40, gendungeon.RoomAttempts, gendungeon.MinRoomSize, gendungeon.MaxRoomSize, 1234)
	dng.RenderToConsole()
}
