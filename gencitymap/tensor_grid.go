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
func NewBasisField(centre vectors.Vec2, size, decay float64, fieldType int) *BasisField {
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
}

// NewRadialField creates a new radial field.
func NewRadialField(centre vectors.Vec2, size, decay float64) *RadialField {
	return &RadialField{
		BasisField: NewBasisField(centre, size, decay, FTRadialField),
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
func NewGridField(centre vectors.Vec2, size, decay, theta float64) *GridField {
	return &GridField{
		BasisField: NewBasisField(centre, size, decay, FTGridField),
		Theta:      theta * math.Pi / 180,
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

type OvalField struct {
	*BasisField
	Theta float64 // Direction of the major axis
	Ratio float64 // Ratio of the major and minor axis
}

// NewOvalField creates a new oval field.
// NOTE: This tends to create spirals, so use with caution.
func NewOvalField(centre vectors.Vec2, size, decay, theta, ratio float64) *OvalField {
	return &OvalField{
		BasisField: NewBasisField(centre, size, decay, FTGridField),
		Theta:      theta,
		Ratio:      ratio,
	}
}

func (r *OvalField) GetTensor(point vectors.Vec2) *Tensor {
	t := point.Sub(r.Centre)
	// Adjust the point to be relative to the centre of the oval
	t = t.Rotate(r.Theta)
	t.X *= r.Ratio

	// Calculate the x/y ratios based on the angle of the major axis
	t1 := t.Y*t.Y - t.X*t.X
	t2 := -2 * t.X * t.Y
	return newTensor(1, [2]float64{t1, t2})
}

func (r *OvalField) GetWeightedTensor(point vectors.Vec2, smooth bool) *Tensor {
	return r.GetTensor(point).Scale(r.getTensorWeight(point, smooth))
}

func (r *OvalField) getTensorWeight(point vectors.Vec2, smooth bool) float64 {
	diff := point.Sub(r.Centre)
	// Adjust the point to be relative to the centre of the oval
	diff = diff.Rotate(r.Theta)
	diff.X *= r.Ratio

	// Interpolates between (0 and 1)^decay
	distanceToCentre := diff.Len()
	if distanceToCentre > r.Size {
		return 0
	}

	weight := math.Pow(1-distanceToCentre/r.Size, r.Decay)
	if smooth {
		weight = math.Pow(weight, 0.5)
	}
	return weight
}
