package gencitymap

import (
	"math"

	"github.com/Flokey82/go_gens/vectors"
)

// Tensor represents a tensor used for generating a city map.
// Represent the matrix as a 2 element list
// [ 0, 1 , 1, -0 ]
type Tensor struct {
	matrix   [2]float64
	r        float64
	theta    float64
	oldTheta bool
}

func newTensor(r float64, matrix [2]float64) *Tensor {
	return &Tensor{
		matrix: matrix,
		r:      r,
	}
}

func newZeroTensor() *Tensor {
	return newTensor(0, [2]float64{0, 0})
}

func fromVector(vector vectors.Vec2) *Tensor {
	t1 := vector.X*vector.X - vector.Y*vector.Y
	t2 := 2 * vector.X * vector.Y
	t3 := t1*t1 - t2*t2
	t4 := 2 * t1 * t2
	return newTensor(1, [2]float64{t3, t4})
}

func (t *Tensor) calculateTheta() float64 {
	if t.r == 0 {
		return 0
	}
	return math.Atan2(t.matrix[1]/t.r, t.matrix[0]/t.r) / 2
}

func (t *Tensor) add(tensor *Tensor, smooth bool) *Tensor {
	t.matrix = mapMatrix(t.matrix, func(v float64, i int) float64 {
		return v*t.r + tensor.matrix[i]*tensor.r
	})

	if smooth {
		t.r = hypot(t.matrix[:]...)
		t.matrix = mapMatrix(t.matrix, func(v float64, i int) float64 {
			return v / t.r
		})
	} else {
		t.r = 2
	}

	t.oldTheta = true
	return t
}

func (t *Tensor) Scale(s float64) *Tensor {
	return t.scale(s)
}

// hypot returns the square root of the sum of squares of its arguments.
func hypot(vals ...float64) float64 {
	var sum float64
	for _, v := range vals {
		sum += v * v
	}
	return math.Sqrt(sum)
}

func mapMatrix(matrix [2]float64, f func(float64, int) float64) [2]float64 {
	return [2]float64{
		f(matrix[0], 0),
		f(matrix[1], 1),
	}
}

func (t *Tensor) getTheta() float64 {
	if t.oldTheta {
		t.theta = t.calculateTheta()
		t.oldTheta = false
	}

	return t.theta
}

func (t *Tensor) scale(s float64) *Tensor {
	t.r *= s
	t.oldTheta = true
	return t
}

func (t *Tensor) rotate(theta float64) *Tensor {
	if theta == 0 {
		return t
	}

	newTheta := t.getTheta() + theta
	if newTheta < math.Pi {
		newTheta += math.Pi
	}

	if newTheta >= math.Pi {
		newTheta -= math.Pi
	}

	t.matrix[0] = math.Cos(2*newTheta) * t.r
	t.matrix[1] = math.Sin(2*newTheta) * t.r
	t.theta = newTheta
	return t
}

func (t *Tensor) getMajor() vectors.Vec2 {
	if t.r == 0 {
		return vectors.NewVec2(0, 0)
	}

	return vectors.NewVec2(math.Cos(t.getTheta()), math.Sin(t.getTheta()))
}

func (t *Tensor) getMinor() vectors.Vec2 {
	if t.r == 0 {
		return vectors.NewVec2(0, 0)
	}

	return vectors.NewVec2(math.Cos(t.getTheta()+math.Pi/2), math.Sin(t.getTheta()+math.Pi/2))
}
