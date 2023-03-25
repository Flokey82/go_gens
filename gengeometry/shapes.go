package gengeometry

import (
	"math"

	"github.com/Flokey82/go_gens/vectors"
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
