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
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed tiles/64x64.png
	Spritesheet_png []byte
)

// SpriteSheet represents a collection of sprite images.
type SpriteSheet struct {
	Floor      *ebiten.Image
	Wall       *ebiten.Image
	StairsUp   *ebiten.Image
	StairsDown *ebiten.Image
	DoorNS     *ebiten.Image
	DoorEW     *ebiten.Image
	DoorNSOpen *ebiten.Image
	DoorEWOpen *ebiten.Image
}

// https://opengameart.org/content/dungeon-crawl-32x32-tiles

/*
Attribution Instructions:
Please include the following text in your distribution package, either as a separate file or incorpolated into your license notice document.
The fourth line is optional: Delete it if you have not modified the tiles.
If you have modified some of the tiles, replace "YOURNAME" by your name:
"Part of (or All) the graphic tiles used in this program is the public domain roguelike tileset 'RLTiles'.
Some of the tiles have been modified by YOURNAME.
You can find the original tileset at: http://rltiles.sf.net"
*/

// LoadSpriteSheet loads the embedded SpriteSheet.
func LoadSpriteSheet(tileSize int) (*SpriteSheet, error) {
	img, _, err := image.Decode(bytes.NewReader(Spritesheet_png))
	if err != nil {
		return nil, err
	}

	sheet := ebiten.NewImageFromImage(img)

	// spriteAt returns a sprite at the provided coordinates.
	spriteAt := func(x, y int) *ebiten.Image {
		return sheet.SubImage(image.Rect(x*tileSize, (y+1)*tileSize, (x+1)*tileSize, y*tileSize)).(*ebiten.Image)
	}

	// Populate SpriteSheet.
	s := &SpriteSheet{}
	s.Floor = spriteAt(10, 4)
	s.Wall = spriteAt(2, 3)
	s.StairsUp = spriteAt(8, 5)
	s.StairsDown = spriteAt(2, 5)
	s.DoorNS = spriteAt(1, 2)
	s.DoorEW = spriteAt(2, 2)
	s.DoorNSOpen = spriteAt(11, 4)
	s.DoorEWOpen = spriteAt(12, 4)

	return s, nil
}
