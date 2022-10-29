// Package simwar provides a simple war simulation based on Lanchester's Laws.
// See: https://en.wikipedia.org/wiki/Lanchester%27s_laws
package simwar

// Army represents an army with a strength and a lethality.
type Army struct {
	Strength  float64 // Strength of the army in soldiers
	Firepower float64 // (0.0 - 1.0) Chance to incapacitate an enemy soldier
	Defense   float64 // (0.0 - 1.0) Chance to survive a successful enemy attack
}

// SimulateSquare runs the simulation for n steps with each step
// being the duration of each step.
// The function returns the time it took for one army to
// defeat the other.
func SimulateSquare(a, b *Army, n int, timestep float64) float64 {
	i := 0
	for ; i < n && a.Strength > 0 && b.Strength > 0; i++ {
		casB := a.Strength * a.Firepower * (1 - b.Defense)
		casA := b.Strength * b.Firepower * (1 - a.Defense)

		a.Strength -= casA * timestep
		b.Strength -= casB * timestep
	}
	return float64(i) * timestep
}

// SimulateLinear runs the simulation for n steps with each step
// being the duration of each step.
// The linear equation takes in account that each soldier can only
// engage with one enemy soldier at a time.
func SimulateLinear(a, b *Army, n int, timestep float64) float64 {
	i := 0
	for ; i < n && a.Strength > 0 && b.Strength > 0; i++ {
		casB := a.Strength * a.Firepower * b.Strength * (1 - b.Defense)
		casA := b.Strength * b.Firepower * a.Strength * (1 - a.Defense)

		a.Strength -= casA * timestep
		b.Strength -= casB * timestep
	}
	return float64(i) * timestep
}
