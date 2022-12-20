package genworldvoronoi

import (
	"math"

	"github.com/ojrac/opensimplex-go"
)

// fbmNoiseCustom returns a function that returns the 'fractal bownian motion'-ish noise value for a given region.
func (m *Geo) fbmNoiseCustom(octaves int, persistence, mx, my, mz, dx, dy, dz float64) func(int) float64 {
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

// genFbmNoise returns the 'fractal bownian motion'-ish noise value for each region.
func (m *Map) genFbmNoise() []float64 {
	fn := m.fbmNoiseCustom(2, 1, 2, 2, 2, 0, 0, 0) // This should be a parameter.
	n := make([]float64, m.mesh.numRegions)
	for r := 0; r < m.mesh.numRegions; r++ {
		n[r] = fn(r)
	}
	return n
}

// getNoiseBandIntersection returns true if a noise value is within a band.
func getNoiseBandIntersection(noisevalue, bandvalue, bandwidth float64) bool {
	return bandvalue-bandwidth/2 <= noisevalue && noisevalue <= bandvalue+bandwidth/2
}

// Noise is a wrapper for opensimplex.Noise, initialized with
// a given seed, persistence, and number of octaves.
type Noise struct {
	Octaves     int
	Persistence float64
	Amplitudes  []float64
	Seed        int64
	OS          opensimplex.Noise
}

// NewNoise returns a new Noise.
func NewNoise(octaves int, persistence float64, seed int64) *Noise {
	n := &Noise{
		Octaves:     octaves,
		Persistence: persistence,
		Amplitudes:  make([]float64, octaves),
		Seed:        seed,
		OS:          opensimplex.NewNormalized(seed),
	}

	// Initialize the amplitudes.
	for i := range n.Amplitudes {
		n.Amplitudes[i] = math.Pow(persistence, float64(i))
	}

	return n
}

// Eval3 returns the noise value at the given point.
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
