package main

import (
	"github.com/Flokey82/go_gens/gendungeon"
)

func main() {
	// Generate 2D dungeon.
	cfg := gendungeon.Config{
		Width:        40,
		Height:       40,
		RoomAttempts: gendungeon.RoomAttempts,
		MinRoomSize:  gendungeon.MinRoomSize,
		MaxRoomSize:  gendungeon.MaxRoomSize,
	}
	dng := gendungeon.GenerateFromConfig(cfg, 1234)
	dng.RenderToConsole()

	// Generate multi level dungeon.
	dng3d := gendungeon.GenerateMultiLevelFromConfig(&gendungeon.ConfigMultiLevel{
		Depth:  3,
		Config: cfg,
	}, 1234)
	dng3d.RenderToConsole()
}
