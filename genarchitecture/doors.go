package genarchitecture

import (
	svg "github.com/ajstarks/svgo"
)

type DoorStyle struct {
	Shape string
	Size  int
	BaseStyle
	// TODO: Handles, hinges, locks, frame etc.
}

func (s DoorStyle) Description() string {
	return s.Shape + " shaped " + s.BaseStyle.Description()
}

const doorOuterStyle = "fill:white;stroke:black;stroke-width:2"

func (s DoorStyle) DrawToSVG(sv *svg.SVG, x, y, width, height int) {
	marginFrame := 16
	drawShape(sv, x, y, width, height, randomString(shapeDoors))
	drawShape(sv, x, y, width, height, s.Shape)
	drawShape(sv, x+marginFrame/2, y+marginFrame/2, width-marginFrame, height-marginFrame, s.Shape)
}

func drawShape(sv *svg.SVG, x, y, width, height int, shape string) {
	switch shape {
	case ShapeRectangle:
		sv.Rect(x, y, width, height, doorOuterStyle)
	case ShapeTrapazoid:
		points := genTrapezoid(x, y, width, height)
		x1, y1 := convertToPairSlices(points)
		sv.Polygon(x1, y1, doorOuterStyle)
	case ShapeOval:
		sv.Ellipse(x+width/2, y+height/2, width/2, height/2, doorOuterStyle)
	case ShapeTriangle:
		points := genTriangle(x, y, width, height)
		x1, y1 := convertToPairSlices(points)
		sv.Polygon(x1, y1, doorOuterStyle)
	case ShapeHexagon:
		points := genHexagon(x, y, width, height)
		x1, y1 := convertToPairSlices(points)
		sv.Polygon(x1, y1, doorOuterStyle)
	case ShapeOctagon:
		points := genOctagon(x, y, width, height)
		x1, y1 := convertToPairSlices(points)
		sv.Polygon(x1, y1, doorOuterStyle)
	case ShapeArch:
		// A simple arch like a door arc (part of an oval)
		sv.Arc(x, y+height, width/2, height, 0, false, true, x+width, y+height, doorOuterStyle)
	}
}

var shapeDoors = []string{
	ShapeRectangle,
	ShapeTrapazoid,
	ShapeOval,
	ShapeTriangle,
	ShapeHexagon,
	ShapeOctagon,
	ShapeArch,
}

func generateDoorStyle(availableMaterials []string) DoorStyle {
	return DoorStyle{
		Shape:     randomString(shapeDoors),
		Size:      randomInt(1, 3),
		BaseStyle: generateBaseStyle(availableMaterials),
	}
}
