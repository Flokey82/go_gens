package gamehex

import (
	"bytes"
	_ "embed"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten"
)

var (
	//go:embed tiles/HexTilesetv3.png
	Spritesheet_png []byte
)

// SpriteSheet represents a collection of sprite images.
type SpriteSheet struct {
	Grass *ebiten.Image
	Dirt  *ebiten.Image
	Water *ebiten.Image
	Snow  *ebiten.Image
	Trees *ebiten.Image
}

// LoadSpriteSheet loads the embedded SpriteSheet.
func LoadSpriteSheet(tileSize int) (*SpriteSheet, error) {
	tileHeight := tileSize * 2 / 2
	tileWidth := tileSize * 2 / 2
	// diff := (float64(tileSize*2) - float64(tileSize)*math.Sqrt(3)) / 2
	img, _, err := image.Decode(bytes.NewReader(Spritesheet_png))
	if err != nil {
		return nil, err
	}

	sheet := ebiten.NewImageFromImage(img)

	// spriteAt returns a sprite at the provided coordinates.
	spriteAt := func(x, y int) *ebiten.Image {
		return sheet.SubImage(image.Rect(x*tileWidth, y*tileHeight, (x+1)*tileWidth, (y+1)*tileHeight)).(*ebiten.Image)
	}

	// Scale the sprite to the correct size.
	// The spritesheet uses 32x32 tiles, but we use a radius of 64 and proper hex geometry.
	scaleSprite := func(img *ebiten.Image) *ebiten.Image {
		newImage := ebiten.NewImage(tileSize*2, int(float64(tileSize)*math.Sqrt(3)))
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, math.Sqrt(3))
		newImage.DrawImage(img, op)
		return newImage
	}

	// Populate SpriteSheet.
	s := &SpriteSheet{}
	s.Grass = scaleSprite(spriteAt(0, 1))
	s.Dirt = scaleSprite(spriteAt(1, 1))
	s.Water = scaleSprite(spriteAt(2, 1))
	s.Snow = scaleSprite(spriteAt(0, 1))
	s.Trees = scaleSprite(spriteAt(0, 2))

	return s, nil
}
