package gencitymap

import (
	"math"

	"github.com/Flokey82/go_gens/vectors"
)

type BasisFieldInterface interface {
	GetTensor(point vectors.Vec2) *Tensor
	GetWeightedTensor(point vectors.Vec2, smooth bool) *Tensor
	GetCentre() vectors.Vec2
}

const (
	// RadialField represents a radial field.
	FTRadialField = iota
	// GridField represents a grid field.
	FTGridField
)

// BasisField represents a basis field used for generating a city map.
type BasisField struct {
	FieldType int
	centre    vectors.Vec2
	size      float64
	decay     float64
}

// NewBasisField creates a new basis field.
func NewBasisField(centre vectors.Vec2, size float64, decay float64, fieldType int) *BasisField {
	return &BasisField{
		centre:    centre,
		size:      size,
		decay:     decay,
		FieldType: fieldType,
	}
}

func (b *BasisField) GetCentre() vectors.Vec2 {
	return b.centre
}

func (b *BasisField) getTensorWeight(point vectors.Vec2, smooth bool) float64 {
	// Interpolates between (0 and 1)^decay
	distanceToCentre := point.Sub(b.centre).Len()
	if distanceToCentre > b.size {
		return 0
	}

	weight := math.Pow(1-distanceToCentre/b.size, b.decay)
	if smooth {
		weight = math.Pow(weight, 0.5)
	}
	return weight
}

type RadialField struct {
	*BasisField
	theta float64
}

// NewRadialField creates a new radial field.
func NewRadialField(centre vectors.Vec2, size float64, decay float64, theta float64) *RadialField {
	return &RadialField{
		BasisField: NewBasisField(centre, size, decay, FTRadialField),
		theta:      theta,
	}
}

func (r *RadialField) GetTensor(point vectors.Vec2) *Tensor {
	t := point.Sub(r.centre)
	t1 := t.Y*t.Y - t.X*t.X
	t2 := -2 * t.X * t.Y
	return newTensor(1, [2]float64{t1, t2})
}

func (r *RadialField) GetWeightedTensor(point vectors.Vec2, smooth bool) *Tensor {
	return r.GetTensor(point).Scale(r.getTensorWeight(point, smooth))
}

type GridField struct {
	*BasisField
	theta float64
}

// NewGridField creates a new grid field.
func NewGridField(centre vectors.Vec2, size float64, decay float64, theta float64) *GridField {
	return &GridField{
		BasisField: NewBasisField(centre, size, decay, FTGridField),
		theta:      theta,
	}
}

func (r *GridField) GetTensor(point vectors.Vec2) *Tensor {
	return newTensor(1, [2]float64{
		math.Cos(2*r.theta) * r.theta,
		math.Sin(2*r.theta) * r.theta,
	})
}

func (r *GridField) GetWeightedTensor(point vectors.Vec2, smooth bool) *Tensor {
	return r.GetTensor(point).Scale(r.getTensorWeight(point, smooth))
}
