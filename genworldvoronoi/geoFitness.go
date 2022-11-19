package genworldvoronoi

import "math"

// getFitnessSteepMountains returns a fitness function with high scores for
// steep terrain close to mountains.
func (m *Geo) getFitnessSteepMountains() func(int) float64 {
	steepness := m.GetSteepness()
	seedMountains, _, _, _ := m.findCollisions()
	distMountains := m.assignDistanceField(seedMountains, make(map[int]bool))
	return func(r int) float64 {
		if m.Elevation[r] <= 0 {
			return -1.0
		}
		chance := steepness[r] * math.Sqrt(m.Elevation[r])
		chance /= (distMountains[r] + 1) / 2
		return chance
	}
}

// getFitnessInlandValleys returns a fitness function with high scores for
// terrain that is not steep and far away from coastlines, mountains, and
// oceans.
func (m *Geo) getFitnessInlandValleys() func(int) float64 {
	steepness := m.GetSteepness()
	seedMountains, seedCoastlines, seedOceans, _ := m.findCollisions()

	// Combine all seed points so we can find the spots furthest away from them.
	var seedAll []int
	seedAll = append(seedAll, seedMountains...)
	seedAll = append(seedAll, seedCoastlines...)
	seedAll = append(seedAll, seedOceans...)
	distAll := m.assignDistanceField(seedAll, make(map[int]bool))
	return func(r int) float64 {
		if m.Elevation[r] <= 0 {
			return -1.0
		}
		chance := 1 - steepness[r]
		chance *= distAll[r]
		return chance
	}
}

func (m *Geo) getFitnessArableLand() func(int) float64 {
	// Prefer flat terrain with reasonable precipitation and at
	// lower altitudes.
	steepness := m.GetSteepness()
	_, maxElev := minMax(m.Elevation)
	return func(r int) float64 {
		temp := m.getRTemperature(r, maxElev)
		if m.Elevation[r] <= 0 || m.Rainfall[r] < 0.1 || temp <= 0 {
			return -1.0
		}
		chance := 1 - steepness[r]
		chance *= m.Rainfall[r]
		chance *= 1 - m.Elevation[r]*m.Elevation[r]
		return chance
	}
}

// getFitnessClimate returns a fitness function that returns high
// scores for regions with high rainfall and high temperatures.
func (m *Geo) getFitnessClimate() func(int) float64 {
	_, maxRain := minMax(m.Rainfall)
	_, maxElev := minMax(m.Elevation)

	return func(r int) float64 {
		rTemp := m.getRTemperature(r, maxElev)
		if rTemp < 0 {
			return 0.1
		}
		return 0.1 + 0.9*(m.Rainfall[r]/maxRain)*math.Sqrt(rTemp/maxTemp)
	}
}
