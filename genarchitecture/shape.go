package genarchitecture

import "math"

const (
	ShapeRectangle = "rectangle"
	ShapeCircle    = "circle"
	ShapeTrapazoid = "trapazoid"
	ShapeOval      = "oval"
	ShapeTriangle  = "triangle"
	ShapeHexagon   = "hexagon"
	ShapeOctagon   = "octagon"
	ShapeArch      = "arch"
)

func genHexagon(x, y, width, height int) [][2]int {
	var points [][2]int
	points = append(points, [2]int{x + width/3, y})
	points = append(points, [2]int{x + width - width/3, y})
	points = append(points, [2]int{x + width, y + height/2})
	points = append(points, [2]int{x + width - width/3, y + height})
	points = append(points, [2]int{x + width/3, y + height})
	points = append(points, [2]int{x, y + height/2})
	return points
}

func genTriangle(x, y, width, height int) [][2]int {
	var points [][2]int
	points = append(points, [2]int{x + width/2, y})
	points = append(points, [2]int{x + width, y + height})
	points = append(points, [2]int{x, y + height})
	return points
}

func genTrapezoid(x, y, width, height int) [][2]int {
	var points [][2]int
	topWidth := width / 2
	points = append(points, [2]int{x, y + height})
	points = append(points, [2]int{x + width, y + height})
	points = append(points, [2]int{x + topWidth + (width-topWidth)/2, y})
	points = append(points, [2]int{x + (width-topWidth)/2, y})
	return points
}

func genOctagon(x, y, width, height int) [][2]int {
	var points [][2]int
	offsetAngle := math.Pi / 8

	// Draw the 8 sides of the octagon
	for i := 0; i < 8; i++ {
		points = append(points, [2]int{
			x + width/2 + int(float64(width/2)*math.Cos(float64(i)*math.Pi/4+offsetAngle)),
			y + height/2 + int(float64(height/2)*math.Sin(float64(i)*math.Pi/4+offsetAngle)),
		})
	}

	return points
}

func genCircle(x, y, width, height, nSteps int) [][2]int {
	var points [][2]int
	for i := 0; i < nSteps; i++ {
		points = append(points, [2]int{
			x + width/2 + int(float64(width/2)*math.Cos(float64(i)*math.Pi/float64(nSteps))),
			y + height/2 + int(float64(height/2)*math.Sin(float64(i)*math.Pi/float64(nSteps))),
		})
	}
	return points
}

func convertToPairSlices(points [][2]int) ([]int, []int) {
	var x, y []int
	for _, p := range points {
		x = append(x, p[0])
		y = append(y, p[1])
	}
	return x, y
}
