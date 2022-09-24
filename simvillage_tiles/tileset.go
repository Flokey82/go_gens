package simvillage_tiles

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten"
)

// TileSet is a convenience wrapper around locating tiles in a tileset.
// TODO: Add support for animations.
type TileSet struct {
	image  *ebiten.Image
	x      int
	xCount int
}

func NewTileSet(imgData []byte) (*TileSet, error) {
	// Decode an image from the image file's byte slice.
	// Now the byte slice is generated with //go:generate for Go 1.15 or older.
	// If you use Go 1.16 or newer, it is strongly recommended to use //go:embed to embed the image file.
	// See https://pkg.go.dev/embed for more details.
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}
	eImg := ebiten.NewImageFromImage(img)
	x, _ := eImg.Size()
	return &TileSet{
		image:  eImg,
		x:      x,
		xCount: x / tileSize,
	}, nil
}

func (t *TileSet) TileImage(idx int) *ebiten.Image {
	return t.image.SubImage(image.Rect((idx%t.xCount)*tileSize, (idx/t.xCount)*tileSize, (idx%t.xCount+1)*tileSize, (idx/t.xCount+1)*tileSize)).(*ebiten.Image)
}
