package gencitymap

import "github.com/Flokey82/go_gens/vectors"

type FieldIntegrator struct {
	field *TensorField
}

func NewFieldIntegrator(field *TensorField) *FieldIntegrator {
	return &FieldIntegrator{
		field: field,
	}
}

func (f *FieldIntegrator) SampleFieldVector(point vectors.Vec2, major bool) vectors.Vec2 {
	tensor := f.field.samplePoint(point)
	if major {
		return tensor.getMajor()
	}
	return tensor.getMinor()
}

func (f *FieldIntegrator) OnLand(point vectors.Vec2) bool {
	return f.field.onLand(point)
}

type EulerIntegrator struct {
	*FieldIntegrator
	params *StreamlineParams
}

func NewEulerIntegrator(field *TensorField, params *StreamlineParams) *EulerIntegrator {
	return &EulerIntegrator{
		FieldIntegrator: NewFieldIntegrator(field),
		params:          params,
	}
}

func (e *EulerIntegrator) Integrate(point vectors.Vec2, major bool) vectors.Vec2 {
	return e.SampleFieldVector(point, major).Mul(e.params.Dstep)
}

type RK4Integrator struct {
	*FieldIntegrator
	params *StreamlineParams
}

func NewRK4Integrator(field *TensorField, params *StreamlineParams) *RK4Integrator {
	return &RK4Integrator{
		FieldIntegrator: NewFieldIntegrator(field),
		params:          params,
	}
}

func (r *RK4Integrator) Integrate(point vectors.Vec2, major bool) vectors.Vec2 {
	k1 := r.SampleFieldVector(point, major)
	k23 := r.SampleFieldVector(point.Add(k1.Mul(r.params.Dstep/2)), major)
	k4 := r.SampleFieldVector(point.Add(k23.Mul(r.params.Dstep)), major)

	return k1.Add(k23.Mul(4)).Add(k4).Mul(r.params.Dstep / 6)
}

type FieldIntegratorIf interface {
	Integrate(point vectors.Vec2, major bool) vectors.Vec2
	OnLand(point vectors.Vec2) bool
	SampleFieldVector(point vectors.Vec2, major bool) vectors.Vec2
}
