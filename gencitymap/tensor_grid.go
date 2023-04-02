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
	FTRadialField = iota // RadialField represents a radial field.
	FTGridField          // GridField represents a grid field.
)

// BasisField represents a basis field used for generating a city map.
type BasisField struct {
	Centre    vectors.Vec2
	FieldType int
	Size      float64
	Decay     float64
}

// NewBasisField creates a new basis field.
func NewBasisField(centre vectors.Vec2, size float64, decay float64, fieldType int) *BasisField {
	return &BasisField{
		Centre:    centre,
		FieldType: fieldType,
		Size:      size,
		Decay:     decay,
	}
}

func (b *BasisField) GetCentre() vectors.Vec2 {
	return b.Centre
}

func (b *BasisField) getTensorWeight(point vectors.Vec2, smooth bool) float64 {
	// Interpolates between (0 and 1)^decay
	distanceToCentre := point.Sub(b.Centre).Len()
	if distanceToCentre > b.Size {
		return 0
	}

	weight := math.Pow(1-distanceToCentre/b.Size, b.Decay)
	if smooth {
		weight = math.Pow(weight, 0.5)
	}
	return weight
}

type RadialField struct {
	*BasisField
	Theta float64
}

// NewRadialField creates a new radial field.
func NewRadialField(centre vectors.Vec2, size float64, decay float64, theta float64) *RadialField {
	return &RadialField{
		BasisField: NewBasisField(centre, size, decay, FTRadialField),
		Theta:      theta,
	}
}

func (r *RadialField) GetTensor(point vectors.Vec2) *Tensor {
	t := point.Sub(r.Centre)
	t1 := t.Y*t.Y - t.X*t.X
	t2 := -2 * t.X * t.Y
	return newTensor(1, [2]float64{t1, t2})
}

func (r *RadialField) GetWeightedTensor(point vectors.Vec2, smooth bool) *Tensor {
	return r.GetTensor(point).Scale(r.getTensorWeight(point, smooth))
}

type GridField struct {
	*BasisField
	Theta float64
}

// NewGridField creates a new grid field.
func NewGridField(centre vectors.Vec2, size float64, decay float64, theta float64) *GridField {
	return &GridField{
		BasisField: NewBasisField(centre, size, decay, FTGridField),
		Theta:      theta,
	}
}

func (r *GridField) GetTensor(point vectors.Vec2) *Tensor {
	return newTensor(1, [2]float64{
		math.Cos(2*r.Theta) * r.Theta,
		math.Sin(2*r.Theta) * r.Theta,
	})
}

func (r *GridField) GetWeightedTensor(point vectors.Vec2, smooth bool) *Tensor {
	return r.GetTensor(point).Scale(r.getTensorWeight(point, smooth))
}
