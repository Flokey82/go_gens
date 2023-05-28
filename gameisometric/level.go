// Copyright 2021 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package gameisometric

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Flokey82/go_gens/gendungeon"
)

// Level represents a Game level.
type Level struct {
	seed         int64                         // Seed used to generate the level
	width        int                           // Width of the level in tiles
	height       int                           // Height of the level in tiles
	tiles        [][]*Tile                     // (Y,X) array of tiles
	tileSize     int                           // Size of each tile in pixels
	currentLevel int                           // Current level in the dungeon
	dng          *gendungeon.DungeonMultiLevel // Generated multi-level dungeon
	ss           *SpriteSheet                  // Spritesheet for the level
}

// Tile returns the tile at the provided coordinates, or nil.
func (l *Level) Tile(x, y int) *Tile {
	if x >= 0 && y >= 0 && x < l.width && y < l.height {
		return l.tiles[y][x]
	}
	return nil
}

// Size returns the size of the Level.
func (l *Level) Size() (int, int) {
	return l.width, l.height
}

// LevelUp moves the player up one level.
func (l *Level) LevelUp() {
	l.currentLevel = (l.currentLevel + 1)
	if l.currentLevel >= len(l.dng.Levels) {
		l.currentLevel = len(l.dng.Levels) - 1
	}
	l.LoadLevel()
}

// LevelDown moves the player down one level.
func (l *Level) LevelDown() {
	l.currentLevel = (l.currentLevel - 1) % len(l.dng.Levels)
	if l.currentLevel < 0 {
		l.currentLevel = 0
	}
	l.LoadLevel()
}

// LoadLevel loads the current level from the dungeon.
func (l *Level) LoadLevel() {
	// Seed the random number generator.
	rng := rand.New(rand.NewSource(l.seed))

	// Get the current level.
	dung := l.dng.Levels[l.currentLevel]

	// Update the level size.
	l.width = dung.Width
	l.height = dung.Height

	isGetWallNeighbor := func(x, y int) Direction {
		var dir Direction
		if x > 0 && dung.Tiles[y][x-1].Material == gendungeon.MatWall {
			dir |= DirectionWest
		}
		if x < l.width-1 && dung.Tiles[y][x+1].Material == gendungeon.MatWall {
			dir |= DirectionEast
		}
		if y > 0 && dung.Tiles[y-1][x].Material == gendungeon.MatWall {
			dir |= DirectionNorth
		}
		if y < l.height-1 && dung.Tiles[y+1][x].Material == gendungeon.MatWall {
			dir |= DirectionSouth
		}
		return dir
	}

	// Generate a unique permutation each time.
	l.tiles = make([][]*Tile, l.height)
	for y := 0; y < l.height; y++ {
		l.tiles[y] = make([]*Tile, l.width)
		for x := 0; x < l.width; x++ {
			t := &Tile{}
			isBorderSpace := x == 0 || y == 0 || x == l.width-1 || y == l.height-1
			val := dung.Tiles[y][x]
			switch {
			case isBorderSpace || val.Material == gendungeon.MatWall:
				t.AddSprite(l.ss.Wall)
			case val.Material == gendungeon.MatStairsDown:
				t.AddSprite(l.ss.StairsDown)
			case val.Material == gendungeon.MatStairsUp:
				t.AddSprite(l.ss.Floor)
				t.AddSprite(l.ss.StairsUp)
			case val.Material == gendungeon.MatDoor:
				t.AddSprite(l.ss.Floor)
				// Check which neighbors are walls and pick the appropriate door sprite.
				if x > 0 && dung.Tiles[y][x-1].Material == gendungeon.MatWall || x < l.width-1 && dung.Tiles[y][x+1].Material == gendungeon.MatWall {
					if rng.Intn(2) == 0 {
						t.AddSprite(l.ss.DoorEW)
					} else {
						t.AddSprite(l.ss.DoorEWOpen)
					}
				} else {
					if rng.Intn(2) == 0 {
						t.AddSprite(l.ss.DoorNS)
					} else {
						t.AddSprite(l.ss.DoorNSOpen)
					}
				}
			default:
				t.AddSprite(l.ss.Floor)
				if rng.Intn(100) < 5 {
					var sprite *Sprite
					switch rng.Intn(3) {
					case 0:
						sprite = l.ss.Table
					case 1:
						sprite = l.ss.Bookcase
					case 2:
						sprite = l.ss.DisplayCase
					}
					// Check if any of the neighbors are walls.
					switch isGetWallNeighbor(x, y) {
					case DirectionNorth:
						t.AddSprite(sprite.S)
					case DirectionEast:
						t.AddSprite(sprite.W)
					case DirectionWest:
						t.AddSprite(sprite.E)
					case DirectionSouth:
						t.AddSprite(sprite.N)
					case DirectionNone:
						t.AddSprite(l.ss.Table.E)
					default:
					}
				}
			}
			l.tiles[y][x] = t
		}
	}
}

type Direction int

const (
	DirectionNone  Direction = 0
	DirectionNorth Direction = 1 << iota
	DirectionSouth
	DirectionEast
	DirectionWest
)

// NewLevel returns a new randomly generated Level.
func NewLevel() (*Level, error) {
	width := 108
	height := 108
	tileSize := 64
	nLevels := 3
	// Generate a dungeon.
	cfg := gendungeon.Config{
		Width:        width,
		Height:       height,
		RoomAttempts: 200,
		MinRoomSize:  5,
		MaxRoomSize:  20,
		AllowNonRect: true,
	}
	seed := time.Now().UnixNano()
	dung := gendungeon.GenerateMultiLevelFromConfig(cfg, nLevels, seed)

	// Load embedded SpriteSheet.
	ss, err := LoadSpriteSheet(tileSize)
	if err != nil {
		return nil, fmt.Errorf("failed to load embedded spritesheet: %s", err)
	}

	// Create a 108x108 Level.
	l := &Level{
		seed:     seed,
		width:    width,
		height:   height,
		tileSize: tileSize,
		dng:      dung,
		ss:       ss,
	}
	l.LoadLevel()
	return l, nil
}
