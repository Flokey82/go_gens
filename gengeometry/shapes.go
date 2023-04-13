package gengeometry

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/Flokey82/go_gens/vectors"
	"github.com/llgcode/draw2d/draw2dimg"
)

// Shape is an interface for shapes.
type Shape interface {
	ConnectionPoints() []vectors.Vec2 // Returns the connection points of the shape.
	GetPath() []vectors.Vec2          // Returns the path of the shape.
}

/*
// O-Shape:
//  ______
// |  __  |
// | |__| |
// |______|
//
// U-Shape:
//  _    _
// | |__| |
// |______|
//
// L-Shape:
//  _
// | |____
// |______|
//
// J-Shape:
//    ____
//   |__  |
//  ____| |
// |______|
//
// T-Shape:
//  ______
// |_    _|
//   |__|
//
// Plus-Shape:
//    __
//  _|  |_
// |_    _|
//   |__|
//
// Rectangle:
//  ______
// |      |
// |______|
//
// H-Shape:
//  _    _
// | |__| |
// |  __  |
// |_|  |_|
//
// Circle:
//   .--.
// /      \
// \      /
//   '--'
//
// Squircle:
//  . -- .
// |      |
// |      |
//  ' -- '
*/

// HShape is a shape that looks like an H.
type HShape struct {
	Width, Length float64
	WingWidth     float64
}

// ConnectionPoints returns the connection points of the shape.
func (h HShape) ConnectionPoints() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: 0, Y: h.Length / 2},
		{X: h.Width / 2, Y: 0},
		{X: h.Width, Y: h.Length / 2},
		{X: h.Width / 2, Y: h.Length},
	}
}

// GetPath returns the path of the shape.
func (h HShape) GetPath() []vectors.Vec2 {
	widthMargin := h.WingWidth
	lengthMargin := h.WingWidth

	return []vectors.Vec2{
		{X: 0, Y: 0},
		{X: widthMargin, Y: 0},
		{X: widthMargin, Y: (h.Length - lengthMargin) / 2},
		{X: h.Width - widthMargin, Y: (h.Length - lengthMargin) / 2},
		{X: h.Width - widthMargin, Y: 0},
		{X: h.Width, Y: 0},
		{X: h.Width, Y: h.Length},
		{X: h.Width - widthMargin, Y: h.Length},
		{X: h.Width - widthMargin, Y: (h.Length + lengthMargin) / 2},
		{X: widthMargin, Y: (h.Length + lengthMargin) / 2},
		{X: widthMargin, Y: h.Length},
		{X: 0, Y: h.Length},
	}
}

// PlusShape is a plus shape.
type PlusShape struct {
	Width, Length, WingWidth float64
}

// ConnectionPoints returns the connection points of the shape.
func (p PlusShape) ConnectionPoints() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: 0, Y: p.Length / 2},
		{X: p.Width / 2, Y: 0},
		{X: p.Width, Y: p.Length / 2},
		{X: p.Width / 2, Y: p.Length},
	}
}

// GetPath returns the path of the shape.
func (p PlusShape) GetPath() []vectors.Vec2 {
	widthMargin := (p.Width - p.WingWidth) / 2
	lengthMargin := (p.Length - p.WingWidth) / 2

	return []vectors.Vec2{
		{X: widthMargin, Y: 0},
		{X: p.Width - widthMargin, Y: 0},
		{X: p.Width - widthMargin, Y: lengthMargin},
		{X: p.Width, Y: lengthMargin},
		{X: p.Width, Y: p.Length - lengthMargin},
		{X: p.Width - widthMargin, Y: p.Length - lengthMargin},
		{X: p.Width - widthMargin, Y: p.Length},
		{X: widthMargin, Y: p.Length},
		{X: widthMargin, Y: p.Length - lengthMargin},
		{X: 0, Y: p.Length - lengthMargin},
		{X: 0, Y: lengthMargin},
		{X: widthMargin, Y: lengthMargin},
	}
}

// UShape is a shape that looks like a U.
type UShape struct {
	Width, Length, WingWidth float64
}

// ConnectionPoints returns the connection points of the shape.
func (u UShape) ConnectionPoints() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: u.Width / 2, Y: 0},
		{X: u.Width, Y: u.Length / 2},
		{X: u.Width / 2, Y: u.WingWidth},
		{X: 0, Y: u.Length / 2},
		{X: u.WingWidth / 2, Y: u.Length},
		{X: u.Width - (u.WingWidth / 2), Y: u.Length},
	}
}

// GetPath returns the path of the shape.
func (u UShape) GetPath() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: 0, Y: 0},
		{X: u.Width, Y: 0},
		{X: u.Width, Y: u.Length},
		{X: u.Width - u.WingWidth, Y: u.Length},
		{X: u.Width - u.WingWidth, Y: u.WingWidth},
		{X: u.WingWidth, Y: u.WingWidth},
		{X: u.WingWidth, Y: u.Length},
		{X: 0, Y: u.Length},
	}
}

// LShape is a shape that looks like an L.
type LShape struct {
	Width, Length, WingWidth float64
}

// ConnectionPoints returns the connection points of the shape.
func (l LShape) ConnectionPoints() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: 0, Y: l.WingWidth / 2},
		{X: l.Width - l.WingWidth/2, Y: l.Length},
		{X: l.WingWidth / 2, Y: l.Length - l.WingWidth},
		{X: l.Width - l.WingWidth/2, Y: 0},
	}
}

// GetPath returns the path of the shape.
func (l LShape) GetPath() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: 0, Y: 0},
		{X: l.Width, Y: 0},
		{X: l.Width, Y: l.Length},
		{X: l.Width - l.WingWidth, Y: l.Length},
		{X: l.Width - l.WingWidth, Y: l.WingWidth},
		{X: 0, Y: l.WingWidth},
	}
}

// JShape is a shape that looks like a J.
type JShape struct {
	Width, Length, WingWidth float64
}

// ConnectionPoints returns the connection points of the shape.
func (j JShape) ConnectionPoints() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: j.WingWidth / 2, Y: 0},
		{X: j.Width, Y: j.Length - j.WingWidth/2},
		{X: j.Width - j.WingWidth, Y: j.Length},
		{X: j.WingWidth / 2, Y: j.WingWidth},
	}
}

// GetPath returns the path of the shape.
func (j JShape) GetPath() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: 0, Y: 0},
		{X: 0, Y: j.WingWidth * 2},
		{X: j.WingWidth, Y: j.WingWidth * 2},
		{X: j.WingWidth, Y: j.WingWidth},
		{X: j.Width - j.WingWidth, Y: j.WingWidth},
		{X: j.Width - j.WingWidth, Y: j.Length},
		{X: j.Width, Y: j.Length},
		{X: j.Width, Y: 0},
	}
}

// RectangleShape is a shape that is a rectangle.
type RectangleShape struct {
	Width, Length float64
}

// ConnectionPoints returns the connection points of the shape.
func (r RectangleShape) ConnectionPoints() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: r.Width / 2, Y: 0},
		{X: r.Width, Y: r.Length / 2},
		{X: r.Width / 2, Y: r.Length},
		{X: 0, Y: r.Length / 2},
	}
}

// GetPath returns the path of the shape.
func (r RectangleShape) GetPath() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: 0, Y: 0},
		{X: r.Width, Y: 0},
		{X: r.Width, Y: r.Length},
		{X: 0, Y: r.Length},
	}
}

// CircleShape is a shape that is a circle.
type CircleShape struct {
	Radius float64 // Radius of the circle
	Steps  int     // Number of points to use to draw the circle
}

// ConnectionPoints returns the connection points of the circle
func (c CircleShape) ConnectionPoints() []vectors.Vec2 {
	// In the middle of each segment
	var res []vectors.Vec2
	angleIncrement := 2 * math.Pi / float64(c.Steps)
	for i := 0; i < c.Steps; i++ {
		angle := float64(i)*angleIncrement + angleIncrement/2
		res = append(res, vectors.Vec2{
			X: c.Radius * math.Cos(angle),
			Y: c.Radius * math.Sin(angle),
		})
	}
	return res
}

// GetPath returns the path of the circle
func (c CircleShape) GetPath() []vectors.Vec2 {
	path := make([]vectors.Vec2, c.Steps)
	angleIncrement := 2 * math.Pi / float64(c.Steps)
	for i := 0; i < c.Steps; i++ {
		angle := float64(i) * angleIncrement
		path[i] = vectors.Vec2{
			X: c.Radius * math.Cos(angle),
			Y: c.Radius * math.Sin(angle),
		}
	}
	return path
}

// TriangleShape is a shape that is a triangle.
type TriangleShape struct {
	Width, Length float64
}

// ConnectionPoints returns the connection points of the triangle
func (t TriangleShape) ConnectionPoints() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: t.Width / 2, Y: 0},
		{X: t.Width, Y: t.Length},
		{X: 0, Y: t.Length},
	}
}

// GetPath returns the path of the triangle
func (t TriangleShape) GetPath() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: 0, Y: 0},
		{X: t.Width, Y: 0},
		{X: t.Width / 2, Y: t.Length},
	}
}

// TrapezoidShape is a shape that is a trapezoid.
type TrapezoidShape struct {
	Width, Length, WidthTop float64
}

// ConnectionPoints returns the connection points of the trapezoid
func (t TrapezoidShape) ConnectionPoints() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: t.Width / 2, Y: 0},
		{X: t.WidthTop, Y: t.Length},
		{X: 0, Y: t.Length},
	}
}

// GetPath returns the path of the trapezoid
func (t TrapezoidShape) GetPath() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: 0, Y: 0},
		{X: t.Width, Y: 0},
		{X: (t.Width - t.WidthTop) / 2, Y: t.Length},
		{X: t.WidthTop + (t.Width-t.WidthTop)/2, Y: t.Length},
	}
}

// SquircleShape is a shape that is a squircle.
type SquircleShape struct {
	Width, Length, Radius float64
	Steps                 int
}

// ConnectionPoints returns the connection points of the squircle
func (s SquircleShape) ConnectionPoints() []vectors.Vec2 {
	// In the middle of each side.
	return []vectors.Vec2{
		{X: s.Width / 2, Y: 0},
		{X: s.Width, Y: s.Length / 2},
		{X: s.Width / 2, Y: s.Length},
		{X: 0, Y: s.Length / 2},
	}
}

// GetPath returns the path of the squircle
// For a squircle, we just need to draw quarter circles at each corner, which
// will give us a squricle.
func (s SquircleShape) GetPath() []vectors.Vec2 {
	fullSteps := s.Steps * 4
	angleIncrement := 2 * math.Pi / float64(fullSteps)
	var path []vectors.Vec2
	var i int
	// Start a quarter circle at the bottom right corner minus / minus the radius.
	circleCenter := vectors.Vec2{X: s.Width - s.Radius, Y: s.Length - s.Radius}
	for ; i < s.Steps; i++ {
		angle := float64(fullSteps-i) * angleIncrement
		path = append(path, vectors.Vec2{
			X: circleCenter.X + s.Radius*math.Cos(angle),
			Y: circleCenter.Y - s.Radius*math.Sin(angle),
		})
	}
	// Start a quarter circle at the bottom left corner plus / minus the radius.
	circleCenter = vectors.Vec2{X: s.Radius, Y: s.Length - s.Radius}
	for ; i < 2*s.Steps; i++ {
		angle := float64(fullSteps-i) * angleIncrement
		path = append(path, vectors.Vec2{
			X: circleCenter.X + s.Radius*math.Cos(angle),
			Y: circleCenter.Y - s.Radius*math.Sin(angle),
		})
	}
	// Start a quarter circle at the top left corner plus / plus the radius.
	circleCenter = vectors.Vec2{X: s.Radius, Y: s.Radius}
	for ; i < 3*s.Steps; i++ {
		angle := float64(fullSteps-i) * angleIncrement
		path = append(path, vectors.Vec2{
			X: circleCenter.X + s.Radius*math.Cos(angle),
			Y: circleCenter.Y - s.Radius*math.Sin(angle),
		})
	}
	// Start a quarter circle at the top right corner minus / plus the radius.
	circleCenter = vectors.Vec2{X: s.Width - s.Radius, Y: s.Radius}
	for ; i < fullSteps; i++ {
		angle := float64(fullSteps-i) * angleIncrement
		path = append(path, vectors.Vec2{
			X: circleCenter.X + s.Radius*math.Cos(angle),
			Y: circleCenter.Y - s.Radius*math.Sin(angle),
		})
	}
	return path
}

// SavePathAsPNG saves the path as a PNG image.
func SavePathAsPNG(path []vectors.Vec2, filename string, scale float64) error {
	// Create a new image
	_, _, maxX, maxY := getPathExtent(path)
	img := image.NewRGBA(image.Rect(0, 0, int(maxX*scale), int(maxY*scale)))
	// Draw the path
	dg := draw2dimg.NewGraphicContext(img)
	for i := 0; i < len(path)-1; i++ {
		dg.SetStrokeColor(color.RGBA{0, 0, 0, 255})
		dg.MoveTo(path[i].X*100, path[i].Y*100)
		dg.LineTo(path[i+1].X*100, path[i+1].Y*100)
		dg.SetLineWidth(1)
		dg.Stroke()
	}
	// Save the image
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func getPathExtent(path []vectors.Vec2) (minX, minY, maxX, maxY float64) {
	minX = math.MaxFloat64
	minY = math.MaxFloat64
	maxX = -math.MaxFloat64
	maxY = -math.MaxFloat64
	for _, p := range path {
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	log.Println(minX, minY, maxX, maxY)
	return
}
