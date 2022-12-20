package genworldvoronoi

import (
	"math"

	"github.com/ojrac/opensimplex-go"
)

func (m *Geo) fbm_noise2(octaves int, persistence, mx, my, mz, dx, dy, dz float64) func(int) float64 {
	// https://thebookofshaders.com/13/
	return func(r int) float64 {
		nx, ny, nz := m.XYZ[3*r]*mx+dx, m.XYZ[3*r+1]*my+dy, m.XYZ[3*r+2]*mz+dz
		var sum float64
		var sumOfAmplitudes float64
		amplitude := 1.0
		for octave := 0; octave < octaves; octave++ {
			frequency := 1 << octave
			fFreq := float64(frequency)
			sum += amplitude * m.noise.OS.Eval3(nx*fFreq, ny*fFreq, nz*fFreq) * float64(octave)
			sumOfAmplitudes += amplitude * float64(octave)
			amplitude *= persistence
		}
		return (sum / sumOfAmplitudes)
	}
}

func (m *Map) genNoise() []float64 {
	fn := m.fbm_noise2(2, 1, 2, 2, 2, 0, 0, 0)
	n := make([]float64, m.mesh.numRegions)
	for r := 0; r < m.mesh.numRegions; r++ {
		n[r] = fn(r)
	}
	return n
}

func getIntersection(noisevalue, bandvalue, bandwidth float64) bool {
	return bandvalue-bandwidth/2 <= noisevalue && noisevalue <= bandvalue+bandwidth/2
}

type Noise struct {
	Octaves     int
	Persistence float64
	Amplitudes  []float64
	Seed        int64
	OS          opensimplex.Noise
}

func NewNoise(octaves int, persistence float64, seed int64) *Noise {
	n := &Noise{
		Octaves:     octaves,
		Persistence: persistence,
		Amplitudes:  make([]float64, octaves),
		Seed:        seed,
		OS:          opensimplex.NewNormalized(seed),
	}
	for i := range n.Amplitudes {
		n.Amplitudes[i] = math.Pow(persistence, float64(i))
	}
	return n
}

func (n *Noise) Eval3(x, y, z float64) float64 {
	var sum float64
	var sumOfAmplitudes float64
	for octave := 0; octave < n.Octaves; octave++ {
		frequency := 1 << octave
		fFreq := float64(frequency)
		sum += n.Amplitudes[octave] * n.OS.Eval3(x*fFreq, y*fFreq, z*fFreq)
		sumOfAmplitudes += n.Amplitudes[octave]
	}
	return sum / sumOfAmplitudes
}

// PlusOneOctave returns a new Noise with one more octave.
func (n *Noise) PlusOneOctave() *Noise {
	return NewNoise(n.Octaves+1, n.Persistence, n.Seed)
}
