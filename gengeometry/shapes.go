package gengeometry

import "github.com/Flokey82/go_gens/vectors"

type Shape interface {
	ConnectionPoints() []vectors.Vec2 // Returns the connection points of the shape.
	GetPath() []vectors.Vec2          // Returns the path of the shape.
}

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
// |______|

type PlusShape struct {
	Width, Length, WingWidth float64
}

func (p PlusShape) ConnectionPoints() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: 0, Y: p.Length / 2},
		{X: p.Width / 2, Y: 0},
		{X: p.Width, Y: p.Length / 2},
		{X: p.Width / 2, Y: p.Length},
	}
}

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
		{X: widthMargin, Y: 0},
	}
}

type UShape struct {
	Width, Length, WingWidth float64
}

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
		{X: 0, Y: 0},
	}
}

type LShape struct {
	Width, Length, WingWidth float64
}

func (l LShape) ConnectionPoints() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: 0, Y: l.WingWidth / 2},
		{X: l.Width - l.WingWidth/2, Y: l.Length},
		{X: l.WingWidth / 2, Y: l.Length - l.WingWidth},
		{X: l.Width - l.WingWidth/2, Y: 0},
	}
}

func (l LShape) GetPath() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: 0, Y: 0},
		{X: l.Width, Y: 0},
		{X: l.Width, Y: l.Length},
		{X: l.Width - l.WingWidth, Y: l.Length},
		{X: l.Width - l.WingWidth, Y: l.WingWidth},
		{X: 0, Y: l.WingWidth},
		{X: 0, Y: 0},
	}
}

type RectangleShape struct {
	Width, Length float64
}

func (r RectangleShape) ConnectionPoints() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: r.Width / 2, Y: 0},
		{X: r.Width, Y: r.Length / 2},
		{X: r.Width / 2, Y: r.Length},
		{X: 0, Y: r.Length / 2},
	}
}

func (r RectangleShape) GetPath() []vectors.Vec2 {
	return []vectors.Vec2{
		{X: 0, Y: 0},
		{X: r.Width, Y: 0},
		{X: r.Width, Y: r.Length},
		{X: 0, Y: r.Length},
		{X: 0, Y: 0},
	}
}
