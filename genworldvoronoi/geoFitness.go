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
// scores for regions with high rainfall high temperatures, and alternatively high flux.
func (m *Geo) getFitnessClimate() func(int) float64 {
	_, maxRain := minMax(m.Rainfall)
	_, maxElev := minMax(m.Elevation)
	_, maxFlux := minMax(m.Flux)

	return func(r int) float64 {
		rTemp := m.getRTemperature(r, maxElev)
		if rTemp < 0 {
			return 0.1
		}
		scoreTemp := math.Sqrt(rTemp / maxTemp)
		scoreRain := m.Rainfall[r] / maxRain
		scoreFlux := math.Sqrt(m.Flux[r] / maxFlux)
		return 0.1 + 0.9*(scoreTemp*(scoreFlux+scoreRain)/2)
	}
}

// CalcFitnessScore calculates the fitness value for all regions based on the
// given fitness function.
//
// - 'sf' is the fitness function for scoring a region.
// - 'distSeedFunc' returns a number of regions from which we maximize the distance when
// calculating the fitness score.
func (m *Geo) CalcFitnessScore(sf func(int) float64, distSeedFunc func() []int) []float64 {
	score := make([]float64, m.mesh.numRegions)

	// Get distance to other seed regions returned by the distSeedFunc.
	r_distance_c := m.assignDistanceField(distSeedFunc(), make(map[int]bool))

	// Get the max distance for normalizing the distance.
	_, maxDistC := minMax(r_distance_c)

	// Calculate the fitness score for each region
	for i := 0; i < m.mesh.numRegions; i++ {
		score[i] = sf(i)

		// Check if we have a valid score.
		if score[i] == -1.0 {
			continue
		}

		// Penalty for proximity / bonus for higher distance to other seed regions.
		//
		// We multiply the score by the distance to other seed regions, amplifying
		// positive scores.
		//
		// NOTE: Originally this was done with some constant values, which might be better
		// since we are here dependent on the current score we have assigned and cannot
		// recover an initially bad score caused by a low water flux.
		if math.IsInf(r_distance_c[i], 0) {
			continue
		}
		dist := (r_distance_c[i] / maxDistC)
		score[i] *= dist // originally: -= 0.02 / (float64(r_distance_c[i]) + 1e-9)
	}
	return score
}