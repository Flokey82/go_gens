package main

import (
	"github.com/Flokey82/go_gens/gameraycast"
	"github.com/Flokey82/go_gens/gendungeon"
)

func main() {
	dng := gendungeon.Generate(40, 40, gendungeon.RoomAttempts, gendungeon.MinRoomSize, gendungeon.MaxRoomSize, 1234)
	m := convertDungeonToMap(dng)
	g := gameraycast.NewGame(m)
	g.Run()
}

func convertDungeonToMap(dng *gendungeon.Dungeon) *gameraycast.Map {
	m := gameraycast.NewMap()
	m.X = dng.Width
	m.Y = dng.Height

	m.Array = make([]int, m.X*m.Y)
	for y := 0; y < m.Y; y++ {
		for x := 0; x < m.X; x++ {
			if dng.Tiles[y][x].Material == gendungeon.MatWall {
				m.Array[y*m.X+x] = 1
			} else {
				m.Array[y*m.X+x] = 0
			}
		}
	}
	return m
}
