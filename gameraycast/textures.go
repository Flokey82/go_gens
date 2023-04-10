package gameraycast

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
)

//go:embed img/wall_brick.png
var wall_brick []byte

//go:embed img/id_door.png
var id_door []byte

//go:embed img/id_stones.png
var id_stones []byte

var (
	texWallBrick TexInterface
	texIDDoor    TexInterface
	texIDStones  TexInterface
)

func init() {
	var err error
	texWallBrick, err = NewImgTextureFromBytes(wall_brick)
	if err != nil {
		panic(err)
	}
	texIDDoor, err = NewImgTextureFromBytes(id_door)
	if err != nil {
		panic(err)
	}
	texIDStones, err = NewImgTextureFromBytes(id_stones)
	if err != nil {
		panic(err)
	}
}

// TexInterface is an interface for textures.
type TexInterface interface {
	Width() int              // Width returns the width of the texture.
	Height() int             // Height returns the height of the texture.
	At(x, y int) color.Color // At returns the color at the given coordinates.
}

// ImgTexture is a texture that is defined by an image.
type ImgTexture struct {
	img image.Image
}

// NewImgTextureFromBytes returns a new texture that is defined by an image.
func NewImgTextureFromBytes(p []byte) (TexInterface, error) {
	img, _, err := image.Decode(bytes.NewReader(p))
	if err != nil {
		return nil, err
	}
	return NewImgTexture(img), nil
}

// NewImgTexture returns a new texture that is defined by an image.
func NewImgTexture(img image.Image) *ImgTexture {
	return &ImgTexture{img: img}
}

// Width returns the width of the texture.
func (t *ImgTexture) Width() int {
	return t.img.Bounds().Dx()
}

// Height returns the height of the texture.
func (t *ImgTexture) Height() int {
	return t.img.Bounds().Dy()
}

// At returns the color at the given coordinates.
func (t *ImgTexture) At(x, y int) color.Color {
	return t.img.At(x, y)
}

// BitTexture is a texture that is defined by a bitmap and a color palette.
type BitTexture struct {
	W      int
	H      int
	Bitmap [][]byte
	Colors []color.RGBA
}

// Width returns the width of the texture.
func (t *BitTexture) Width() int {
	return t.W
}

// Height returns the height of the texture.
func (t *BitTexture) Height() int {
	return t.H
}

// At returns the color at the given coordinates.
func (t *BitTexture) At(x, y int) color.Color {
	return t.Colors[t.Bitmap[y][x]]
}

// Example of a bit texture.
var texWallStone = BitTexture{
	W: 8,
	H: 8,
	Bitmap: [][]byte{
		{1, 1, 1, 1, 1, 1, 1, 1},
		{0, 0, 0, 1, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1},
		{0, 1, 0, 0, 0, 1, 0, 0},
		{1, 1, 1, 1, 1, 1, 1, 1},
		{0, 0, 0, 1, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1},
		{0, 1, 0, 0, 0, 1, 0, 0},
	},
	Colors: []color.RGBA{
		{255, 241, 232, 255},
		{194, 195, 199, 255},
	},
}

var floorColor = color.NRGBA{0x33, 0x33, 0x33, 0xff}

func getWallTexture(t int) TexInterface {
	if t == 2 {
		return texIDDoor
	}
	return texIDStones
}

func getWallColor(t int) color.RGBA {
	if t <= 0 && t >= len(wallTypeColors) {
		return wallTypeColors[0]
	}
	return wallTypeColors[t]
}

var wallTypeColors = []color.RGBA{
	{0, 0, 0, 255},
	{255, 0, 0, 255},
	{0, 255, 0, 255},
}

func darkenColor(c color.RGBA, amount float64) color.RGBA {
	return color.RGBA{
		uint8(float64(c.R) * amount),
		uint8(float64(c.G) * amount),
		uint8(float64(c.B) * amount),
		c.A,
	}
}
