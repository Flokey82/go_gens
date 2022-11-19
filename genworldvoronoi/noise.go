package genworldvoronoi

import (
	"math"
)

// Initialize the noise amplitudes for use in our heightmap.
var amplitudes []float64

func init() {
	const persistence = 2.0 / 3.0
	amplitudes = make([]float64, 5)
	for i := range amplitudes {
		amplitudes[i] = math.Pow(persistence, float64(i))
	}
}

// fbm_noise returns a noise value for the given xyz coordinate.
func (m *Geo) fbm_noise(nx, ny, nz float64) float64 {
	sum := 0.0
	sumOfAmplitudes := 0.0
	for octave := 0; octave < len(amplitudes); octave++ {
		frequency := 1 << octave
		sum += amplitudes[octave] * m.noise.Eval3(nx*float64(frequency), ny*float64(frequency), nz*float64(frequency))
		sumOfAmplitudes += amplitudes[octave]
	}
	return sum / sumOfAmplitudes
}

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
			sum += amplitude * m.noise.Eval3(nx*fFreq, ny*fFreq, nz*fFreq) * float64(octave)
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
