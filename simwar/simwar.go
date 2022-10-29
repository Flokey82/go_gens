// Package simwar provides a simple war simulation based on Lanchester's Laws.
// See: https://en.wikipedia.org/wiki/Lanchester%27s_laws
package simwar

// Army represents an army with a strength and a lethality.
type Army struct {
	Strength  float64
	Firepower float64
}

// SimulateSquare runs the simulation for n steps with each step
// being the duration of timestep until one army is reduced
// to a strength of 0.
// The function returns the time it took for one army to
// defeat the other.
func SimulateSquare(a, b *Army, n int, timestep float64) float64 {
	i := 0
	for ; i < n && a.Strength > 0 && b.Strength > 0; i++ {
		casB := a.Strength * a.Firepower
		casA := b.Strength * b.Firepower

		a.Strength -= casA * timestep
		b.Strength -= casB * timestep
	}
	return float64(i) * timestep
}
