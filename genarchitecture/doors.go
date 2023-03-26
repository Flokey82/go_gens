package genarchitecture

import (
	"math"

	svg "github.com/ajstarks/svgo"
)

type DoorStyle struct {
	Shape string
	Size  int
	BaseStyle
	// TODO: Handles, hinges, locks, etc.
}

func (s DoorStyle) Description() string {
	return "A " + s.Shape + " door " + s.BaseStyle.Description()
}

func (s DoorStyle) DrawToSVG(sv *svg.SVG, x, y, width, height int) {
	switch s.Shape {
	case ShapeRectangle:
		sv.Rect(x, y, width, height)
	case ShapeCircle:
		sv.Circle(x+width/2, y+height/2, height/2)
	case ShapeTrapazoid:
		var points [][2]int
		topWidth := width / 2
		points = append(points, [2]int{x, y + height})
		points = append(points, [2]int{x + width, y + height})
		points = append(points, [2]int{x + topWidth + (width-topWidth)/2, y})
		points = append(points, [2]int{x + (width-topWidth)/2, y})

		var x, y []int
		for _, p := range points {
			x = append(x, p[0])
			y = append(y, p[1])
		}

		sv.Polygon(x, y)
	case ShapeOval:
		sv.Ellipse(x+width/2, y+height/2, width/2, height/2)
	case ShapeTriangle:
		var points [][2]int
		points = append(points, [2]int{x + width/2, y})
		points = append(points, [2]int{x + width, y + height})
		points = append(points, [2]int{x, y + height})

		var x, y []int
		for _, p := range points {
			x = append(x, p[0])
			y = append(y, p[1])
		}

		sv.Polygon(x, y)
	case ShapeHexagon:
		var points [][2]int
		points = append(points, [2]int{x, y})
		points = append(points, [2]int{x + width, y})
		points = append(points, [2]int{x + width + width/2, y + height/2})
		points = append(points, [2]int{x + width, y + height})
		points = append(points, [2]int{x, y + height})
		points = append(points, [2]int{x - width/2, y + height/2})

		var x, y []int
		for _, p := range points {
			x = append(x, p[0])
			y = append(y, p[1])
		}

		sv.Polygon(x, y)
	case ShapeOctagon:
		var points [][2]int

		// Draw the 8 sides of the octagon
		for i := 0; i < 8; i++ {
			points = append(points, [2]int{
				x + width/2 + int(float64(width/2)*math.Cos(float64(i)*math.Pi/4+math.Pi/8)),
				y + height/2 + int(float64(height/2)*math.Sin(float64(i)*math.Pi/4+math.Pi/8)),
			})
		}

		var x, y []int
		for _, p := range points {
			x = append(x, p[0])
			y = append(y, p[1])
		}

		sv.Polygon(x, y)
	case ShapeArch:
		// A simple arch like a door arc (part of an oval)
		sv.Arc(x, y+height, width/2, height, 0, false, true, x+width, y+height)
	}

}

var shapeDoors = []string{
	ShapeRectangle,
	ShapeCircle,
	ShapeTrapazoid,
	ShapeOval,
	ShapeTriangle,
	ShapeHexagon,
	ShapeOctagon,
	ShapeArch,
	ShapeArch,
	ShapeArch,
	ShapeArch,
	ShapeArch,
}

func generateDoorStyle(availableMaterials []string) DoorStyle {
	return DoorStyle{
		Shape:     randomString(shapeDoors),
		Size:      randomInt(1, 3),
		BaseStyle: generateBaseStyle(availableMaterials),
	}
}
