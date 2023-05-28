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
	"golang.org/x/image/draw"
)

var (
	//go:embed tiles/64x64.png
	Spritesheet_png []byte

	//go:embed tiles/kenney/longTable_N.png
	TableN_png []byte

	//go:embed tiles/kenney/longTable_E.png
	TableE_png []byte

	//go:embed tiles/kenney/longTable_S.png
	TableS_png []byte

	//go:embed tiles/kenney/longTable_W.png
	TableW_png []byte

	//go:embed tiles/kenney/displayCase_N.png
	DisplayCaseN_png []byte

	//go:embed tiles/kenney/displayCase_E.png
	DisplayCaseE_png []byte

	//go:embed tiles/kenney/displayCase_S.png
	DisplayCaseS_png []byte

	//go:embed tiles/kenney/displayCase_W.png
	DisplayCaseW_png []byte

	//go:embed tiles/kenney/bookcaseBooks_N.png
	BookcaseBooksN_png []byte

	//go:embed tiles/kenney/bookcaseBooks_E.png
	BookcaseBooksE_png []byte

	//go:embed tiles/kenney/bookcaseBooks_S.png
	BookcaseBooksS_png []byte

	//go:embed tiles/kenney/bookcaseBooks_W.png
	BookcaseBooksW_png []byte
)

// SpriteSheet represents a collection of sprite images.
type SpriteSheet struct {
	Floor       *ebiten.Image
	Wall        *ebiten.Image
	StairsUp    *ebiten.Image
	StairsDown  *ebiten.Image
	DoorNS      *ebiten.Image
	DoorEW      *ebiten.Image
	DoorNSOpen  *ebiten.Image
	DoorEWOpen  *ebiten.Image
	Table       *Sprite
	DisplayCase *Sprite
	Bookcase    *Sprite
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
		return sheet.SubImage(image.Rect(x*tileSize, y*tileSize, (x+1)*tileSize, (y+1)*tileSize)).(*ebiten.Image)
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

	s.Table = NewSprite(
		ReadImage(TableN_png),
		ReadImage(TableE_png),
		ReadImage(TableS_png),
		ReadImage(TableW_png),
	)

	s.DisplayCase = NewSprite(
		ReadImage(DisplayCaseN_png),
		ReadImage(DisplayCaseE_png),
		ReadImage(DisplayCaseS_png),
		ReadImage(DisplayCaseW_png),
	)

	s.Bookcase = NewSprite(
		ReadImage(BookcaseBooksN_png),
		ReadImage(BookcaseBooksE_png),
		ReadImage(BookcaseBooksS_png),
		ReadImage(BookcaseBooksW_png),
	)

	return s, nil
}

// ReadImage reads an image from a byte slice and downscales it to the tile size.
func ReadImage(data []byte) *ebiten.Image {
	src, _, _ := image.Decode(bytes.NewReader(data))

	// NOTE that the tile is isometric.
	tileSize := 64

	// Ratio of the tile width to the tile height.
	// ratio := float64(src.Bounds().Dx()) / float64(src.Bounds().Dy())

	// Create a new image with the tile size.
	img := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))

	// Resize:
	draw.NearestNeighbor.Scale(img, image.Rect(16+2, 0-8, tileSize-16+2, tileSize-8), src, src.Bounds(), draw.Over, nil)
	// move the image to the center
	// TODO: Crop it to 64x64:
	ret := ebiten.NewImageFromImage(img)
	return ret
}

type Sprite struct {
	N, E, S, W *ebiten.Image
}

func NewSprite(n, e, s, w *ebiten.Image) *Sprite {
	return &Sprite{
		N: n,
		E: e,
		S: s,
		W: w,
	}
}
