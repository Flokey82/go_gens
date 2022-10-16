package genworldvoronoi

import "math"

// Metal resource flags starting with the cheapest metal.
const (
	ResMetIron = 1 << iota
	ResMetCopper
	ResMetLead
	ResMetTin
	ResMetSilver
	ResMetGold
	ResMetPlatinum
)

// Gemstone resource flags starting with the cheapest gem.
const (
	ResGemAmethyst = 1 << iota
	ResGemTopaz
	ResGemSapphire
	ResGemEmerald
	ResGemRuby
	ResGemDiamond
)

// getFitnessSteepMountains returns a fitness function with high scores for
// steep terrain close to mountains.
func (m *Map) getFitnessSteepMountains() func(int) float64 {
	steepness := m.getRSteepness()
	seedMountains, _, _, _ := m.findCollisions()
	distMountains := m.assignDistanceField(seedMountains, make(map[int]bool))
	return func(r int) float64 {
		if m.r_elevation[r] <= 0 {
			return -1.0
		}
		chance := steepness[r] * math.Sqrt(m.r_elevation[r])
		chance /= (distMountains[r] + 1) / 2
		return chance
	}
}

// getFitnessInlandValleys returns a fitness function with high scores for
// terrain that is not steep and far away from coastlines, mountains, and
// oceans.
func (m *Map) getFitnessInlandValleys() func(int) float64 {
	steepness := m.getRSteepness()
	seedMountains, seedCoastlines, seedOceans, _ := m.findCollisions()

	// Combine all seed points so we can find the spots furthest away from them.
	var seedAll []int
	seedAll = append(seedAll, seedMountains...)
	seedAll = append(seedAll, seedCoastlines...)
	seedAll = append(seedAll, seedOceans...)
	distAll := m.assignDistanceField(seedAll, make(map[int]bool))
	return func(r int) float64 {
		if m.r_elevation[r] <= 0 {
			return -1.0
		}
		chance := 1 - steepness[r]
		chance *= distAll[r]
		return chance
	}
}

func (m *Map) getFitnessArableLand() func(int) float64 {
	// Prefer flat terrain with reasonable precipitation and at
	// lower altitudes.
	steepness := m.getRSteepness()
	_, maxElev := minMax(m.r_elevation)
	return func(r int) float64 {
		temp := m.getRTemperature(r, maxElev)
		if m.r_elevation[r] <= 0 || m.r_rainfall[r] < 0.1 || temp <= 0 {
			return -1.0
		}
		chance := 1 - steepness[r]
		chance *= m.r_rainfall[r]
		chance *= 1 - m.r_elevation[r]*m.r_elevation[r]
		return chance
	}
}

func (m *Map) resourceFitness() []float64 {
	fitness := make([]float64, m.mesh.numRegions)
	f := m.getFitnessSteepMountains()
	for r := range fitness {
		fitness[r] = f(r)
	}
	return fitness
}

func (m *Map) placeResources() {
	// NOTE: This currently sucks.
	// TODO: Use fitness function instead or in addition.

	// Place metals.
	// Metals can be found mainly in mountains, so steepness
	// will be an indicator along with the distance from the
	// mountain seed points.

	// Place gemstones.
	// Gemstones can be found mainly in inland valleys, so
	// distance from the coastlines, mountains, and oceans
	// will be an indicator.

	// Place forests.
	// Forests can be found mainly in valleys, so steepness
	// will be an indicator along with the distance from the
	// valley's center.

	// Place arable land.
	// Arable land can be found mainly in valleys, so steepness
	// will be an indicator along with the distance from the
	// valley's center.

	// Place potential quarry sites.
	// Potential quarry sites can be found mainly in mountains,

	steepness := m.getRSteepness()
	// distMountains, _, _, _ := m.findCollisions()

	// https://www.reddit.com/r/worldbuilding/comments/kbmnd6/a_guide_to_placing_resources_on_fictional_worlds/
	const (
		chancePlatinum = 0.005
		chanceGold     = 0.025
		chanceSilver   = 0.04
		chanceTin      = 0.1
		chanceCopper   = 0.06
		chanceLed      = 0.07
		chanceIron     = 0.75
		chanceCoal     = 0.9
	)

	// NOTE: By encoding the resources as bit flags, we can easily
	// determine the value of a region given the assumption that
	// each resource is twice (or half) as valuable as the previous
	// resource. This will be handy for fitness functions and such.
	//
	// I feel pretty clever about this one, but it's not realistic.
	metals := make([]byte, len(steepness))
	for r := 0; r < m.mesh.numRegions; r++ {
		if steepness[r] > 0.9 && m.r_elevation[r] > 0.5 {
			switch rv := m.rand.Float64(); {
			case rv < chancePlatinum:
				metals[r] |= ResMetPlatinum
			case rv < chanceGold:
				metals[r] |= ResMetGold
			case rv < chanceSilver:
				metals[r] |= ResMetSilver
			case rv < chanceTin:
				metals[r] |= ResMetTin
			case rv < chanceCopper:
				metals[r] |= ResMetCopper
			case rv < chanceLed:
				metals[r] |= ResMetLead
			case rv < chanceIron:
				metals[r] |= ResMetIron
			}
		}
	}
	m.res_metals_r = metals

	const (
		chanceDiamond  = 0.005
		chanceRuby     = 0.025
		chanceEmerald  = 0.04
		chanceSapphire = 0.1
		chanceTopaz    = 0.06
		chanceAmethyst = 0.07
		// chanceQuartz   = 0.75 // Usually goes hand in hand with gold?
		// chanceFlint    = 0.9
	)

	gems := make([]byte, len(steepness))
	for r := 0; r < m.mesh.numRegions; r++ {
		if steepness[r] > 0.9 && m.r_elevation[r] > 0.5 {
			switch rv := m.rand.Float64(); {
			case rv < chanceDiamond:
				gems[r] |= ResGemDiamond
			case rv < chanceRuby:
				gems[r] |= ResGemRuby
			case rv < chanceEmerald:
				gems[r] |= ResGemEmerald
			case rv < chanceSapphire:
				gems[r] |= ResGemSapphire
			case rv < chanceTopaz:
				gems[r] |= ResGemTopaz
			case rv < chanceAmethyst:
				gems[r] |= ResGemAmethyst
				// case rv < chanceQuartz:
				//	gems[r] |= ResGemQuartz
				// case rv < chanceFlint:
				//	gems[r] |= ResGemFlint
			}
		}
	}
	m.res_gems_r = metals

	// This attempts some weird variation of:
	// https://www.redblobgames.com/x/1736-resource-placement/
	/*
		nA := m.fbm_noise2(5, 0.5, 5, 5, 5, 0, 0, 0)
		nB := m.fbm_noise2(7, 0.5, 5, 5, 5, 0, 0, 0)
		resources := make([]byte, len(steepness))
		for r := range steepness {
			noiseVal := (nA(r) + nB(r) + m.r_elevation[r]) / 3
			if m.getIntersection(noiseVal, 0.75, 0.01) {
				resources[r] |= ResMetPlatinum
			}
			//chance /= float64(distMountains[r])
		}

		nC := m.fbm_noise2(2, 0.5, 5, 5, 5, 0, 0, 0)
		nD := m.fbm_noise2(7, 0.5, 5, 5, 5, 0, 0, 0)
		for r := range steepness {
			noiseVal := (nC(r) + nD(r) + m.r_elevation[r]) / 3
			if m.getIntersection(noiseVal, 0.75, 0.02) {
				resources[r] |= ResMetGold
			}
			//chance /= float64(distMountains[r])
		}

		nC = m.fbm_noise2(2, 0.5, 1, 1, 1, 0, 0, 0)
		nD = m.fbm_noise2(5, 0.1, 1, 1, 1, 0, 0, 0)
		for r := range steepness {
			noiseVal := (-1*(nC(r)+nD(r)) + m.r_elevation[r]) / 3
			if m.getIntersection(noiseVal, 0.52, 0.07) {
				resources[r] |= ResMetIron
			}
			//chance /= float64(distMountains[r])
		}
	*/

	//m.r_metals = resources
}

func (m *Map) getIntersection(noisevalue, bandvalue, bandwidth float64) bool {
	return bandvalue-bandwidth/2 <= noisevalue && noisevalue <= bandvalue+bandwidth/2
}

func (m *Map) fbm_noise2(octaves int, persistence, mx, my, mz, dx, dy, dz float64) func(int) float64 {
	return func(r int) float64 {
		nx, ny, nz := m.r_xyz[3*r]*mx+dx, m.r_xyz[3*r+1]*my+dy, m.r_xyz[3*r+2]*mz+dz
		var sum float64
		var sumOfAmplitudes float64
		amplitude := 1.0
		for octave := 0; octave < octaves; octave++ {
			frequency := 1 << octave
			fFreq := float64(frequency)
			sum += amplitude * m.noise.Eval3(nx*fFreq, ny*fFreq, nz*fFreq) * float64(octave)
			sumOfAmplitudes += amplitude
			amplitude *= persistence
		}
		return ((sum / sumOfAmplitudes) + 1) / 2
	}
}
