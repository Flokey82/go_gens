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
	dng3d := gendungeon.GenerateMultiLevelFromConfig(cfg, 3, 1234)
	dng3d.RenderToConsole()

	// Generate multi level dungeon with custom room sizes.
	dng3d = gendungeon.NewDungeonMultiLevel()

	// Create 3 levels.
	cfg = gendungeon.Config{
		Width:        10,
		Height:       10,
		RoomAttempts: gendungeon.RoomAttempts,
		MinRoomSize:  2,
		MaxRoomSize:  4,
	}
	dng = gendungeon.GenerateFromConfig(cfg, 1234)
	dng3d.AddLevel(dng)

	cfg = gendungeon.Config{
		Width:        20,
		Height:       20,
		RoomAttempts: gendungeon.RoomAttempts,
		MinRoomSize:  3,
		MaxRoomSize:  5,
	}
	dng = gendungeon.GenerateFromConfig(cfg, 1234)
	dng3d.AddLevel(dng)

	cfg = gendungeon.Config{
		Width:        40,
		Height:       40,
		RoomAttempts: gendungeon.RoomAttempts,
		MinRoomSize:  5,
		MaxRoomSize:  10,
	}
	dng = gendungeon.GenerateFromConfig(cfg, 1234)
	dng3d.AddLevel(dng)

	// Finalize, create stairs.
	dng3d.CreateStairs()

	// Render to console.
	dng3d.RenderToConsole()
}
