package genworldvoronoi

import (
	"log"
	"math"

	"github.com/Flokey82/go_gens/genbiome"
)

// sumResources returns the sum of the resource flag IDs in the byte.
// This is a convenience function for determining the approximate
// value of local resources.
// NOTE: In theory one could just cast the int to a byte and use the
// value like that, but the value would be a power of 2, which might
// be too stark a difference.
func sumResources(r byte) int {
	sum := 0
	for i := 0; i < 8; i++ {
		if r&(1<<i) != 0 {
			sum += i + 1
		}
	}
	return sum
}

const (
	ResourceTypeMetal = iota
	ResourceTypeGem
	ResourceTypeStone
)

// getRegsWithResource returns the regions that have the specified resource.
func (m *Geo) getRegsWithResource(resource byte, resourceType int) []int {
	// Pick the correct resource slice.
	var search []byte
	switch resourceType {
	case ResourceTypeMetal:
		search = m.Metals
	case ResourceTypeGem:
		search = m.Gems
	case ResourceTypeStone:
		search = m.Stones
	}

	// Find the regions that have the specified resource.
	var regions []int
	for r, val := range search {
		if val&resource != 0 {
			regions = append(regions, r)
		}
	}
	return regions
}

// Resources maps regions to natural resources.
type Resources struct {
	Metals  []byte // Metal ores
	Gems    []byte // Gemstones
	Stones  []byte // Rocks or minerals
	Various []byte // Other resources
	Wood    []byte // Wood
}

func newResources(size int) *Resources {
	return &Resources{
		Metals:  make([]byte, size),
		Gems:    make([]byte, size),
		Stones:  make([]byte, size),
		Various: make([]byte, size),
		Wood:    make([]byte, size),
	}
}

func (res *Resources) sumRegion(r int) int {
	return sumResources(res.Metals[r]) + sumResources(res.Gems[r]) + sumResources(res.Stones[r]) + sumResources(res.Various[r]) + sumResources(res.Wood[r])
}

func (m *Geo) resourceFitness() []float64 {
	fitness := make([]float64, m.mesh.numRegions)
	f := m.getFitnessSteepMountains()
	for r := range fitness {
		fitness[r] = f(r)
	}
	return fitness
}

func (m *Geo) placeResources() {
	// NOTE: This currently sucks.
	// TODO: Use fitness function instead or in addition.

	// Place metals.
	// Metals can be found mainly in mountains, so steepness
	// will be an indicator along with the distance from the
	// mountain seed points.
	m.placeMetals()

	// Place gemstones.
	// Gemstones can be found mainly in inland valleys, so
	// distance from the coastlines, mountains, and oceans
	// will be an indicator.
	m.placeGems()

	// Place forests.
	// Forests can be found mainly in valleys, so steepness
	// will be an indicator along with the distance from the
	// valley's center.
	m.placeForests()

	// Place potential quarry sites.
	// Potential quarry sites can be found mainly in mountains,
	m.placeStones()

	// Place energy sources and other resources.
	// Oil, coal, and natural gas, as well as geothermal energy
	// and magical handwavium... and clay, and salt, and stuff.
	m.placeVarious()

	// Place arable land.
	// Arable land can be found mainly in valleys, so steepness
	// will be an indicator along with the distance from the
	// valley's center.
}

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

const ResMaxMetals = 7

func metalToString(metalID int) string {
	switch 1 << metalID {
	case ResMetIron:
		return "Iron"
	case ResMetCopper:
		return "Copper"
	case ResMetLead:
		return "Lead"
	case ResMetTin:
		return "Tin"
	case ResMetSilver:
		return "Silver"
	case ResMetGold:
		return "Gold"
	case ResMetPlatinum:
		return "Platinum"
	default:
		return "Unknown"
	}
}

func (m *Geo) placeMetals() {
	steepness := m.GetSteepness()
	// distMountains, _, _, _ := m.findCollisions()

	// https://www.reddit.com/r/worldbuilding/comments/kbmnd6/a_guide_to_placing_resources_on_fictional_worlds/
	const (
		chancePlatinum = 0.005
		chanceGold     = chancePlatinum + 0.020
		chanceSilver   = chanceGold + 0.040
		chanceCopper   = chanceSilver + 0.06
		chanceLead     = chanceCopper + 0.07
		chanceTin      = chanceLead + 0.1
		chanceIron     = chanceTin + 0.4
	)
	fn := m.fbmNoiseCustom(2, 1, 2, 2, 2, 0, 0, 0)
	fm := m.getFitnessSteepMountains()

	// NOTE: By encoding the resources as bit flags, we can easily
	// determine the value of a region given the assumption that
	// each resource is twice (or half) as valuable as the previous
	// resource. This will be handy for fitness functions and such.
	//
	// I feel pretty clever about this one, but it's not realistic.
	m.resetRand()
	metals := make([]byte, len(steepness))

	// TODO: Use noise intersection instead of rand.
	for r := 0; r < m.mesh.numRegions; r++ {
		if fm(r) > 0.9 {
			switch rv := math.Abs(m.rand.NormFloat64() * fn(r)); {
			case rv < chancePlatinum:
				metals[r] |= ResMetPlatinum
			case rv < chanceGold:
				metals[r] |= ResMetGold
			case rv < chanceSilver:
				metals[r] |= ResMetSilver
			case rv < chanceCopper:
				metals[r] |= ResMetCopper
			case rv < chanceLead:
				metals[r] |= ResMetLead
			case rv < chanceTin:
				metals[r] |= ResMetTin
			case rv < chanceIron:
				metals[r] |= ResMetIron
			}
		}
	}
	m.Metals = metals

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

// Gemstone resource flags starting with the cheapest gem.
const (
	ResGemAmethyst = 1 << iota
	ResGemTopaz
	ResGemSapphire
	ResGemEmerald
	ResGemRuby
	ResGemDiamond
)

const ResMaxGems = 6

func gemToString(gemsID int) string {
	switch 1 << gemsID {
	case ResGemAmethyst:
		return "Amethyst"
	case ResGemTopaz:
		return "Topaz"
	case ResGemSapphire:
		return "Sapphire"
	case ResGemEmerald:
		return "Emerald"
	case ResGemRuby:
		return "Ruby"
	case ResGemDiamond:
		return "Diamond"
	default:
		return "Unknown"
	}
}

func (m *Geo) placeGems() {
	steepness := m.GetSteepness()
	const (
		chanceDiamond  = 0.005
		chanceRuby     = chanceDiamond + 0.025
		chanceEmerald  = chanceRuby + 0.04
		chanceSapphire = chanceEmerald + 0.05
		chanceTopaz    = chanceSapphire + 0.06
		chanceAmethyst = chanceTopaz + 0.1
		// chanceQuartz   = 0.75 // Usually goes hand in hand with gold?
		// chanceFlint    = 0.9
	)

	gems := make([]byte, len(steepness))
	for r := 0; r < m.mesh.numRegions; r++ {
		if steepness[r] > 0.9 && m.Elevation[r] > 0.5 {
			switch rv := math.Abs(m.rand.NormFloat64()); {
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
	m.Gems = gems
}

// Stone resource flags starting with the most common stone.
// NOTE: Clay?
const (
	ResStoSandstone = 1 << iota
	ResStoLimestone
	ResStoChalk
	ResStoSlate
	ResStoMarble
	ResStoGranite
	ResStoBasalt
	ResStoObsidian
)

const ResMaxStones = 7

func stoneToString(stoneID int) string {
	switch 1 << stoneID {
	case ResStoSandstone:
		return "Sandstone"
	case ResStoLimestone:
		return "Limestone"
	case ResStoChalk:
		return "Chalk"
	case ResStoSlate:
		return "Slate"
	case ResStoMarble:
		return "Marble"
	case ResStoGranite:
		return "Granite"
	case ResStoBasalt:
		return "Basalt"
	case ResStoObsidian:
		return "Obsidian"
	default:
		return "Unknown"
	}
}

func (m *Geo) placeStones() {
	log.Println("placing stones is not implemented")

	// Chalk:
	// Ancient Chalk beds formed on the floor of ancient seas.
	//
	// Limestone:
	// The Chalk later solidifies into Limestone. Can be placed where hill
	// meet grasslands in non wet areas.
	//
	// Flint:
	// Flint (also called Chert) forms as lumps between layers and in cavities
	// left in the sea floor in these Chalk beds.
	//
	// Marble:
	// Marble is formed from Limestone that has been subjected to intense heat
	// and pressure. Marble will be placed near mountain ranges.
	//
	// Obsidian:
	// Obsidian is formed when water flows over volcanic lava to cool it rapidly.
	// Placed near volcanic plate boundaries that no longer have large amounts of
	// water. Water breaks down obsidian over time.
	//
	// Granite:
	// Granite is formed when molten rock is slowly cooled. It forms the bottom
	// layer of all land continents. Placed along two land type convergent boundaries
	// on the uplifted side where it is raised to the surface, making quarrying easy.
	//
	// Sandstone:
	// Sandstone is formed when sand is deposited in large quantities and under goes
	// large amounts of pressure, heat, and drainage causing the sand and other
	// minerals to "cement" together. Placed near ancient drainage basins that deposited
	// sand from deserts or beaches, or alternatively where hills or mountains meet a
	// dry desert.
	//
	// Basalt:
	// Basalt is formed when lava cools quickly. Placed near volcanic plate boundaries
	// that have large amounts of water. Water breaks down basalt over time.
	//
	// Slate:
	// Slate is formed when shale is subjected to intense heat and pressure. Slate
	// will be placed near mountain ranges.

	// Initialize the stone map.
	stones := make([]byte, m.mesh.numRegions)

	biomeFunc := m.getRegWhittakerModBiomeFunc()
	steepness := m.GetSteepness()

	// Generate a distance field for volcanoes, mountains, and faultlines.
	var volcanoes, mountains, faultlines []int
	stopSea := make(map[int]bool)
	isBeach := make(map[int]bool)
	for r := 0; r < m.mesh.numRegions; r++ {
		if m.RegionIsVolcano[r] {
			volcanoes = append(volcanoes, r)
		}
		if m.RegionIsMountain[r] {
			mountains = append(mountains, r)
		}
		if math.Abs(m.RegionCompression[r]) > 0.1 {
			faultlines = append(faultlines, r)
		}
		if m.Elevation[r] <= 0.0 {
			stopSea[r] = true
		} else {
			// Check if the region is a beach.
			for _, n := range m.GetRegNeighbors(r) {
				if m.Elevation[n] <= 0.0 {
					isBeach[r] = true
					break
				}
			}
		}
	}

	distVolcanoes := m.assignDistanceField(volcanoes, stopSea)
	distMountains := m.assignDistanceField(mountains, stopSea)
	distFaultlines := m.assignDistanceField(faultlines, stopSea)

	// Loop through all the regions and place stones based on the region's
	// properties.
	for r := 0; r < m.mesh.numRegions; r++ {
		// Skip water regions.
		if m.Elevation[r] <= 0.0 {
			continue
		}

		// Get the region's biome.
		biome := biomeFunc(r)

		// Check if we have sandstone (beach, or desert).
		if biome == genbiome.WhittakerModBiomeSubtropicalDesert || isBeach[r] {
			stones[r] |= ResStoSandstone
		}

		// Chalk and limestone.
		if biome == genbiome.WhittakerModBiomeTemperateGrassland && steepness[r] > 0.1 {
			// If we are close to mountains, we have marble.
			if distMountains[r] < 2 {
				stones[r] |= ResStoMarble
			} else if !m.isRegRiver(r) && !m.isRegLakeOrWaterBody(r) {
				// Check if we have limestone (dryer, hilly grassland)
				stones[r] |= ResStoLimestone
			} else if m.Rainfall[r] > 0.5 {
				// Check if we have chalk (wetter, hilly grassland)
				stones[r] |= ResStoChalk
			}
		}

		// Obsidian, and basalt.
		// For these stones, we need to check if we are near a volcano or faultline.
		if distVolcanoes[r] < 2 || distFaultlines[r] < 2 {
			// Check if we have obsidian (near a volcano).
			if distVolcanoes[r] < 2 {
				stones[r] |= ResStoObsidian
			}

			// Check if we have basalt (near a faultline).
			if distFaultlines[r] < 2 {
				stones[r] |= ResStoBasalt
			}
		}

		// Check if we have granite (near a mountain and faultline).
		if distMountains[r] < 3 && distFaultlines[r] < 2 {
			stones[r] |= ResStoGranite
		} else if steepness[r] > 0.2 && distMountains[r] > 2 && distMountains[r] < 5 {
			// Slate.
			// For slate, we need to check if we are near a mountain range or if the region
			// is steep.
			stones[r] |= ResStoSlate
		}
	}

	// Assign the stone map.
	m.Stones = stones
}

const (
	ResVarClay = 1 << iota
	ResVarSulfur
	ResVarSalt
	ResVarCoal
	ResVarOil
	ResVarGas
)

const ResMaxVarious = 6

func variousToString(v int) string {
	switch 1 << v {
	case ResVarClay:
		return "clay"
	case ResVarSulfur:
		return "sulfur"
	case ResVarSalt:
		return "salt"
	case ResVarCoal:
		return "coal"
	case ResVarOil:
		return "oil"
	case ResVarGas:
		return "gas"
	default:
		return "unknown"
	}
}

func (m *Geo) placeVarious() {
	varRes := make([]byte, m.mesh.numRegions)
	biomeFunc := m.getRegWhittakerModBiomeFunc()
	steepness := m.GetSteepness()
	for r := 0; r < m.mesh.numRegions; r++ {
		if m.Elevation[r] <= 0.0 {
			continue
		}
		biome := biomeFunc(r)

		if m.RegionIsVolcano[r] {
			varRes[r] |= ResVarSulfur
		}

		if m.RegionIsMountain[r] {
			varRes[r] |= ResVarCoal
		}

		if m.isRegRiver(r) && steepness[r] > 0.1 && steepness[r] < 0.3 {
			varRes[r] |= ResVarClay
		}

		if biome == genbiome.WhittakerModBiomeHotSwamp {
			varRes[r] |= ResVarGas
		}

		// TODO: Salt, oil, coal.
	}
	m.Various = varRes
}

// The 8 most important types of wood.
const (
	ResWoodOak = 1 << iota
	ResWoodBirch
	ResWoodPine
	ResWoodSpruce
	ResWoodCedar
	ResWoodShrub
	ResWoodFir
	ResWoodPalm
)

const ResMaxWoods = 8

func woodToString(v int) string {
	switch 1 << v {
	case ResWoodOak:
		return "oak"
	case ResWoodBirch:
		return "birch"
	case ResWoodPine:
		return "pine"
	case ResWoodSpruce:
		return "spruce"
	case ResWoodCedar:
		return "cedar"
	case ResWoodShrub:
		return "shrub"
	case ResWoodFir:
		return "fir"
	case ResWoodPalm:
		return "palm"
	default:
		return "unknown"
	}
}

func (m *Geo) placeForests() {
	// Get all biomes that are forested.
	// Place trees in those biomes based on the biome's tree type(s).
	// Of course it can't be too steep.
	biomeFunc := m.getRegWhittakerModBiomeFunc()
	//steepness := m.GetSteepness()

	wood := make([]byte, m.mesh.numRegions)
	for r := 0; r < m.mesh.numRegions; r++ {
		if m.Elevation[r] <= 0.0 {
			continue
		}

		// NOTE: This is absolute garbage. It's just a quick hack to get some forests
		// in the world.
		biome := biomeFunc(r)
		if biome == genbiome.WhittakerModBiomeTemperateRainforest {
			wood[r] |= ResWoodOak
		} else if biome == genbiome.WhittakerModBiomeTemperateSeasonalForest {
			wood[r] |= ResWoodOak
			wood[r] |= ResWoodBirch
		} else if biome == genbiome.WhittakerModBiomeTropicalRainforest {
			wood[r] |= ResWoodOak
			wood[r] |= ResWoodPalm
		} else if biome == genbiome.WhittakerModBiomeTropicalSeasonalForest {
			wood[r] |= ResWoodOak
			wood[r] |= ResWoodPalm
			wood[r] |= ResWoodBirch
		} else if biome == genbiome.WhittakerModBiomeBorealForestTaiga {
			wood[r] |= ResWoodPine
			wood[r] |= ResWoodSpruce
			wood[r] |= ResWoodCedar
		} else if biome == genbiome.WhittakerModBiomeTundra {
			wood[r] |= ResWoodSpruce
			wood[r] |= ResWoodCedar
			wood[r] |= ResWoodFir
			wood[r] |= ResWoodShrub
		} else if biome == genbiome.WhittakerModBiomeWetlands {
			wood[r] |= ResWoodShrub
			wood[r] |= ResWoodFir
			wood[r] |= ResWoodCedar
			wood[r] |= ResWoodOak
			wood[r] |= ResWoodBirch
		} else if biome == genbiome.WhittakerModBiomeWoodlandShrubland {
			wood[r] |= ResWoodShrub
			wood[r] |= ResWoodFir
			wood[r] |= ResWoodCedar
			wood[r] |= ResWoodOak
			wood[r] |= ResWoodBirch
		}
	}
	m.Wood = wood
}
