package genarchitecture

var shapeWalls = []string{
	ShapeRectangle,
	ShapeTrapazoid,
	ShapeArch,
}

type WallStyle struct {
	Shape  string // wall shape
	Height int    // wall height
	BaseStyle
}

func generateWallStyle(availableMaterials []string) WallStyle {
	return WallStyle{
		Shape:     randomString(shapeWalls),
		Height:    randomInt(1, 3),
		BaseStyle: generateBaseStyle(availableMaterials),
	}
}

func (s WallStyle) Description() string {
	return s.Shape + " walls of " + s.BaseStyle.Description()
}
