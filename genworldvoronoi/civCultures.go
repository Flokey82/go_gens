package genworldvoronoi

import (
	"log"
	"math"
	"math/rand"

	"github.com/Flokey82/go_gens/genbiome"
)

// GetCulture returns the culture of the given region (if any).
func (m *Civ) GetCulture(r int) *Culture {
	// NOTE: This sucks. This should be done better.
	if m.RegionToCulture[r] <= 0 {
		return nil
	}
	for _, c := range m.Cultures {
		if c.ID == m.RegionToCulture[r] {
			return c
		}
	}
	return nil
}

// PlaceNCultures places n cultures on the map.
// This code is based on:
// https://github.com/Azgaar/Fantasy-Map-Generator/blob/master/modules/cultures-generator.js
func (m *Civ) PlaceNCultures(n int) {
	m.resetRand()
	m.placeNCultures(n)
	m.ExpandCultures()
}

// placeNCultures places n cultures on the map.
func (m *Civ) placeNCultures(n int) {
	// The fitness function, returning a score from
	// 0.0 to 1.0 for a given region.
	var scoreFunc func(int) float64

	// The distance seed point function, returning
	// seed points/regions that we want to be far
	// away from.
	var distSeedFunc func() []int

	regCultureFunc := m.getRegionCultureTypeFunc()
	climateFitness := m.getFitnessClimate()
	scoreFunc = func(r int) float64 {
		if m.Elevation[r] <= 0 {
			return 0
		}
		return math.Sqrt((climateFitness(r) + 3.0) / 4.0)
	}

	// For now we maximize the distance to other cultures.
	distSeedFunc = func() []int {
		var cultureSeeds []int
		for _, c := range m.Cultures {
			cultureSeeds = append(cultureSeeds, c.ID)
		}
		return cultureSeeds
	}
	// Place n cities of the given type.
	for i := 0; i < n; i++ {
		c := m.PlaceCulture(regCultureFunc, scoreFunc, distSeedFunc)
		log.Printf("placing culture %d: %s", i, c.Name)
	}
}

// ExpandCultures expands the cultures on the map based on their expansionism,
// terrain preference, and distance to other cultures.
func (m *Civ) ExpandCultures() {
	// The cultural centers will be the seed points for the expansion.
	var seeds []int
	originToCulture := make(map[int]*Culture)
	for _, c := range m.Cultures {
		seeds = append(seeds, c.ID)
		originToCulture[c.ID] = c
	}

	rCellType := m.getRegCellTypes()
	_, maxElev := minMax(m.Elevation)
	territoryWeightFunc := m.getTerritoryWeightFunc()
	biomeWeight := m.getTerritoryBiomeWeightFunc()
	m.RegionToCulture = m.regPlaceNTerritoriesCustom(seeds, func(o, u, v int) float64 {
		c := originToCulture[o]

		// Get the cost to expand to this biome.
		gotBiome := m.getAzgaarRegionBiome(v, m.Elevation[v]/maxElev, maxElev)
		biomePenalty := biomeWeight(o, u, v) * float64(genbiome.AzgaarBiomeMovementCost[gotBiome]) / 100

		// Check if we have a non-native biome, if so we apply an additional penalty.
		// NOTE: This check has been disabled for now.
		// if m.getAzgaarRegionBiome(o, m.Elevation[o]/maxElev, maxElev) != gotBiome {
		biomePenalty *= c.Type.BiomeCost(gotBiome)
		// }

		cellTypePenalty := c.Type.CellTypeCost(rCellType[v])
		return biomePenalty + cellTypePenalty*territoryWeightFunc(o, u, v)/c.Expansionism
	})

	// TODO: There are small islands that do not have a culture...
	// We should (or could) fix that.
	for _, c := range m.Cultures {
		c.Regions = nil
		// Collect all regions that are part of the current culture.
		for r, cu := range m.RegionToCulture {
			if cu == c.ID {
				c.Regions = append(c.Regions, r)
			}
		}
		c.Stats = m.getStats(c.Regions)
	}
}

// Culture represents a culture.
//
// Also see: https://ck3.paradoxwikis.com/Culture
//
// TODO:
//
// # VALUES
//
// The type of culture will also influence their values.
// For example, a nomadic or a hunting culture will have
// higher regart for martial skills and lower regard for
// sophistication.
//
// I would propose to have a point pool for these attributes
// and then randomly assign them to the various skills with
// a certain distribution based on the culture type.
//
// # CRAFTS AND SKILLS
//
// The culture will also have a set of crafts, arts, and skills
// that they are good at. This will be based on their values and
// the environment they live in.
//
// For example, a nomadic culture will be good at hunting and
// tracking. They might be exceptional at crafting leather wares
// (like saddles), and bows.
//
// A river culture will have a ready supply of clay, and will
// be good at pottery, but not at crafin leather wares like
// saddles.
//
// A mountain culture will be good at mining, prospecting, and
// stone carving. If they mine precious metals, they will be
// good at jewelry making.
//
// A naval culture will be good at ship building, and sailing.
// Since they might have sea shells as a ready supply, they
// might be good at jewelry making.
//
// # ARTS
//
// Arts will be based on the environment and the values of the
// culture. For example, a river culture might focus on water
// and rivers, as well as the flora and fauna related to rivers.
//
// A mountain culture might use iconography to represent the
// harshness of the mountains, and the gifts of the mines.
type Culture struct {
	ID           int         // Region where the culture originates
	Name         string      // Name of the culture
	Type         CultureType // Type of the culture
	Expansionism float64     // Expansionism of the culture
	Martialism   float64     // Martial skills of the culture
	// Sophistication float64
	// Extremism      float64 ?
	// Openness       float64 ?
	// Parent    *Culture
	// Children  []*Culture
	// Extinct   bool
	Language *Language // Language of the culture
	Religion *Religion // Religion of the culture

	// TODO: DO NOT CACHE THIS!
	Regions []int
	*Stats
}

func (c *Culture) Log() {
	log.Printf("The Folk of %s (%s): %d regions", c.Name, c.Type.String(), len(c.Regions))
	log.Printf("Followers of %s (%s)", c.Religion.Name, c.Religion.Type)
	c.Stats.Log()
}

func (m *Civ) newCulture(r int, cultureType CultureType) *Culture {
	lang := GenLanguage(m.Seed + int64(r))
	c := &Culture{
		ID:           r,
		Name:         lang.MakeName(),
		Type:         cultureType,
		Expansionism: cultureType.Expansionism(),
		Martialism:   cultureType.Martialism(),
		Language:     lang,
	}
	c.Religion = m.genFolkReligion(c)
	return c
}

// PlaceCulture places another culture on the map at the region with the highest fitness score.
func (m *Civ) PlaceCulture(regCultureFunc func(int) CultureType, scoreFunc func(int) float64, distSeedFunc func() []int) *Culture {
	// Score all regions, pick highest score.
	var newculture int
	lastMax := math.Inf(-1)
	for i, val := range m.CalcCityScore(scoreFunc, distSeedFunc) {
		if val > lastMax {
			newculture = i
			lastMax = val
		}
	}
	c := m.newCulture(newculture, regCultureFunc(newculture))
	m.Cultures = append(m.Cultures, c)
	return c
}

// PlaceCultureAt places a culture at the given region.
// TODO: Allow specifying the culture type?
func (m *Civ) PlaceCultureAt(r int) *Culture {
	c := m.newCulture(r, m.getRegionCultureTypeFunc()(r))
	c.Regions = []int{r}
	c.Stats = m.getStats(c.Regions)
	m.Cultures = append(m.Cultures, c)
	m.RegionToCulture[r] = r
	// TODO: Grow / Expand this culture.
	// NOTE: This might be quite expensive, so we might want to
	// avoid this calling here, or at least limit the regions
	// we process to the ones that are close to the new culture.
	m.ExpandCultures()
	return c
}

type CultureType int

// Culture types.
const (
	CultureTypeWildland CultureType = iota
	CultureTypeGeneric
	CultureTypeRiver
	CultureTypeLake
	CultureTypeNaval
	CultureTypeNomadic
	CultureTypeHunting
	CultureTypeHighland
)

// String returns the string representation of a given culture type.
func (c CultureType) String() string {
	switch c {
	case CultureTypeWildland:
		return "Wildland"
	case CultureTypeGeneric:
		return "Generic"
	case CultureTypeRiver:
		return "River"
	case CultureTypeLake:
		return "Lake"
	case CultureTypeNaval:
		return "Naval"
	case CultureTypeNomadic:
		return "Nomadic"
	case CultureTypeHunting:
		return "Hunting"
	case CultureTypeHighland:
		return "Highland"
	default:
		return "Unknown"
	}
}

// Expansionism returns the expansionism of a given culture type.
func (t CultureType) Expansionism() float64 {
	// TODO: This is a random attractiveness value of the capital.
	// https://azgaar.wordpress.com/2017/11/21/settlements/
	// I introduced two custom parameters — disbalance and power.
	// Each capital has unique attractiveness power, which is randomly
	// assigned to it based on a disbalance value. Disbalance is the same
	// for all capitals, it only controls the randomness of power
	// definition. Calculating a distance to the closest capital we
	// multiply this value by capital’s power. If capital located not on
	// the same island, we double the distance as it should not be easy
	// for city to get an overseas possessions. As all capitals have
	// different “powers”, the regions vary in area. For some reasons
	// user may want regions having almost the same area, so the disbalance
	// value could be changed.
	powerInputValue := 1.0
	base := 1.0 // Generic
	switch t {
	case CultureTypeLake:
		base = 0.8
	case CultureTypeNaval:
		base = 1.5
	case CultureTypeRiver:
		base = 0.9
	case CultureTypeNomadic:
		base = 1.5
	case CultureTypeHunting:
		base = 0.7
	case CultureTypeHighland:
		base = 1.2
	}
	return roundToDecimals(((rand.Float64()*powerInputValue)/2+1)*base, 1)
}

// Martialism returns the martialism of a given culture type.
func (t CultureType) Martialism() float64 {
	powerInputValue := 1.0
	base := 1.0 // Generic
	switch t {
	case CultureTypeLake:
		base = 0.8
	case CultureTypeNaval:
		base = 1.5
	case CultureTypeRiver:
		base = 0.9
	case CultureTypeNomadic:
		base = 1.4
	case CultureTypeHunting:
		base = 1.4
	case CultureTypeHighland:
		base = 1.1
	}
	return roundToDecimals(((rand.Float64()*powerInputValue)/2+1)*base, 1)
}

// CellTypeCost returns the cost of crossing / navigating a given cell type for a given culture.
func (t CultureType) CellTypeCost(cellType int) float64 {
	// TODO: Make use of this

	// Land near coast / coastline / coastal land strip / "beach"?.
	if cellType == 1 {
		if t == CultureTypeNaval || t == CultureTypeLake {
			// Naval cultures or lake cultures have an easier time navigating
			// coastal areas or shores of lakes.
			return 1.0
		}
		if t == CultureTypeNomadic {
			// Nomadic cultures have a harder time navigating coastal areas or
			// shores of lakes.
			return 1.6
		}
		// All other cultures have a small penalty for coastal areas.
		return 1.2
	}

	// Land slightly further inland.
	if cellType == 2 {
		if t == CultureTypeNaval || t == CultureTypeNomadic {
			// Small penalty for land with distance 2 to ocean for navals and nomads.
			return 1.3
		}
		// All other cultures do not have appreciable penalty.
		return 1.0
	}

	// Not water near coast (deep ocean/coastal land).
	if cellType != -1 {
		if t == CultureTypeNaval || t == CultureTypeLake {
			// Penalty for mainland for naval and lake cultures
			return 2.0
		}
	}
	return 1.0
}

// BiomeCost returns the cost for traversion / expanding into a given biome.
func (t CultureType) BiomeCost(biome int) float64 {
	if t == CultureTypeHunting {
		// Non-native biome penalty for hunters.
		return 5.0
	}
	if t == CultureTypeNomadic && (biome == genbiome.AzgaarBiomeTropicalSeasonalForest ||
		biome == genbiome.AzgaarBiomeTemperateDeciduousForest ||
		biome == genbiome.AzgaarBiomeTropicalRainforest ||
		biome == genbiome.AzgaarBiomeTemperateRainforest ||
		biome == genbiome.AzgaarBiomeTaiga) {
		// Forest biome penalty for nomads.
		return 10.0
	}
	// General non-native biome penalty.
	return 2.0
}

// getRegHaven returns the closest neighbor region that is a water cell, which
// can be used as a haven, and returns the number of water neighbors, indicating
// the harbor size.
//
// If no haven is found, -1 is returned.
func (m *Civ) getRegHaven(i int) (int, int) {
	// get all neighbors that are below or at sea level.
	var water []int
	for _, nb := range m.GetRegNeighbors(i) {
		if m.Elevation[nb] <= 0.0 {
			water = append(water, nb)
		}
	}

	// No water neighbors, return -1.
	if len(water) == 0 {
		return -1, 0
	}

	// Get distances of i to each water neighbor.
	// get the closest water neighbor.
	iLatLon := m.LatLon[i]
	closest := -1
	var minDist float64
	for _, nb := range water {
		nbLatLon := m.LatLon[nb]
		dist := haversine(iLatLon[0], iLatLon[1], nbLatLon[0], nbLatLon[1])
		if closest == -1 || dist < minDist {
			minDist = dist
			closest = nb
		}
	}
	// store the closest water neighbor as the haven.
	// store the number of water neighbors as the harbor.
	return closest, len(water)
}

// getRegCellTypes maps the region to its cell type.
//
// NOTE: Currently this depends on the region graph, which will break
// things once we increas or decrease the number of regions on the map as
// the distance between regions will change with the region density.
//
// Value meanings:
//
// -2: deep ocean or large lake
// -1: region is a water cell next to a land cell (lake shore/coastal water)
// +1: region is a land cell next to a water cell (lake shore/coastal land)
// +2: region is a land cell next to a coastal land cell
// >2: region is inland
func (m *Civ) getRegCellTypes() []int {
	var oceanRegs, landRegs []int
	for r, elev := range m.Elevation {
		if elev <= 0.0 {
			oceanRegs = append(oceanRegs, r)
		} else {
			landRegs = append(landRegs, r)
		}
	}
	regDistanceOcean := m.assignDistanceField(oceanRegs, make(map[int]bool))
	regDistanceLand := m.assignDistanceField(landRegs, make(map[int]bool))

	cellType := make([]int, m.mesh.numRegions)
	for i := range cellType {
		// Is it water?
		if m.Elevation[i] <= 0.0 {
			// Figure out if it has a land neighbor.
			// If so, it is -1 (water near coast)
			if regDistanceLand[i] <= 1 {
				cellType[i] = -1
			} else {
				// If not, it is -2 (water far from coast)
				cellType[i] = -2
			}
		} else {
			// Figure out if it has a water neighbor.
			// If so, it is 1 (land near coast)
			if regDistanceOcean[i] <= 1 {
				cellType[i] = 1
			} else {
				// If not, it is >=2 (land far from coast)
				cellType[i] = int(regDistanceOcean[i])
			}
		}
	}
	return cellType
}

// Landmark feature types.
const (
	FeatureTypeOcean     = "ocean"
	FeatureTypeSea       = "sea"
	FeatureTypeLake      = "lake"
	FeatureTypeGulf      = "gulf"
	FeatureTypeIsle      = "isle"
	FeatureTypeContinent = "continent"
)

// getRegionFeatureTypeFunc returns a function that returns the feature type of
// a given region.
func (m *Civ) getRegionFeatureTypeFunc() func(int) string {
	return func(i int) string {
		if i < 0 {
			return ""
		}
		if waterbodyID := m.Waterbodies[i]; waterbodyID >= 0 {
			switch wbSize := m.WaterbodySize[waterbodyID]; {
			case wbSize > m.mesh.numRegions/25:
				return FeatureTypeOcean
			case wbSize > m.mesh.numRegions/100:
				return FeatureTypeSea
			case wbSize > m.mesh.numRegions/500:
				return FeatureTypeGulf
			default:
				return FeatureTypeLake
			}
		}
		if landmassID := m.Landmasses[i]; landmassID >= 0 {
			if m.LandmassSize[landmassID] < m.mesh.numRegions/100 {
				return FeatureTypeIsle
			}
			return FeatureTypeContinent
		}
		return ""
	}
}

// getAzgaarRegionBiome returns the biome for a given region as per Azgaar's map generator.
func (m *Civ) getAzgaarRegionBiome(r int, elev, maxElev float64) int {
	return genbiome.GetAzgaarBiome(int(20.0*m.Moisture[r]), int(m.getRegTemperature(r, maxElev)), int(elev*100))
}

// getRegionCutureTypeFunc returns a function that returns the culture type suitable for a given region.
func (m *Civ) getRegionCultureTypeFunc() func(int) CultureType {
	cellType := m.getRegCellTypes()
	getType := m.getRegionFeatureTypeFunc()
	biomeFunc := m.getRegWhittakerModBiomeFunc()
	_, maxElev := minMax(m.Elevation)

	// Return culture type based on culture center region.
	return func(r int) CultureType {
		eleVal := m.Elevation[r] / maxElev
		gotBiome := m.getAzgaarRegionBiome(r, eleVal, maxElev)
		log.Println(gotBiome)
		log.Println(biomeFunc(r))

		// Desert and grassland means a nomadic culture.
		// BUT: Grassland is extremely well suited for farming... Which is not nomadic.
		if eleVal < 0.7 && (gotBiome == genbiome.AzgaarBiomeHotDesert ||
			gotBiome == genbiome.AzgaarBiomeColdDesert ||
			gotBiome == genbiome.AzgaarBiomeGrassland) {
			return CultureTypeNomadic // high penalty in forest biomes and near coastline
		}

		// Montane cultures in high elevations and hills
		// that aren't deserts or grassland.
		if eleVal > 0.3 {
			return CultureTypeHighland // no penalty for hills and moutains, high for other elevations
		}

		// Get the region (if any) that represents the haven for this region.
		// A haven is the closest neighbor that is a water body.
		// NOTE: harborSize indicates the number of neighbors that are water.
		rHaven, harborSize := m.getRegHaven(r)
		havenType := getType(rHaven) // Get the type of the haven region.
		regionType := getType(r)
		log.Println(havenType, regionType)

		// Ensure only larger lakes will result in the 'lake' culture type.
		if havenType == FeatureTypeLake && m.WaterbodySize[rHaven] > 5 {
			return CultureTypeLake // low water cross penalty and high for growth not along coastline
		}

		// If we have a harbor (more than 1 water neighbor), or are on an island,
		// we are potentially a naval culture.
		if (harborSize > 0 && P(0.1) && havenType != FeatureTypeLake) ||
			(harborSize == 1 && P(0.6)) ||
			(regionType == FeatureTypeIsle && P(0.4)) {
			return CultureTypeNaval // low water cross penalty and high for non-along-coastline growth
		}

		// If we are on a big river (flux > 2*rainfall), we are a river culture.
		if m.isRegBigRiver(r) {
			return CultureTypeRiver // no River cross penalty, penalty for non-River growth
		}

		// If we are inland (cellType > 2) and in one of the listed biomes,
		// we are a hunting culture.
		if cellType[r] > 2 && (gotBiome == genbiome.AzgaarBiomeSavanna ||
			gotBiome == genbiome.AzgaarBiomeTropicalRainforest ||
			gotBiome == genbiome.AzgaarBiomeTemperateRainforest ||
			gotBiome == genbiome.AzgaarBiomeWetland ||
			gotBiome == genbiome.AzgaarBiomeTaiga ||
			gotBiome == genbiome.AzgaarBiomeTundra || // Tundra is also nomadic?
			gotBiome == genbiome.AzgaarBiomeGlacier) {
			return CultureTypeHunting // high penalty in non-native biomes
		}

		// TODO: Wildlands?
		// TODO: What culture would have originated in seasonal forests?
		log.Println(gotBiome, gotBiome, gotBiome, gotBiome)
		return CultureTypeGeneric
	}
}

/*

func (m *Map) getBiomeCost(cultureCenter int, biome int, cType CultureType) int {
	_, maxElev := minMax(m.r_elevation)
	eleVal := m.r_elevation[cultureCenter] / maxElev
	gotBiome := m.getRBiomeTEMP(cultureCenter, eleVal, maxElev)
	if gotBiome == biome {
		return 10 // tiny penalty for native biome
	}
	if cType == CultureTypeHunting {
		return genbiome.AzgaarBiomeMovementCost[biome] * 5 // non-native biome penalty for hunters
	}
	if cType == CultureTypeNomadic && biome > 4 && biome < 10 {
		return genbiome.AzgaarBiomeMovementCost[biome] * 10 // forest biome penalty for nomads
	}
	return genbiome.AzgaarBiomeMovementCost[biome] * 2 // general non-native biome penalty
}

func (m *Map) getHeightCost(i int, h float64, cType CultureType) int {
	f = pack.features[cells.f[i]]
	a = cells.area[i]
	if cType == CultureTypeLake && f.Type == "lake" {
		return 10 // no lake crossing penalty for Lake cultures
	}
	if cType == CultureTypeNaval && h < 20 {
		return a * 2 // low sea/lake crossing penalty for Naval cultures
	}
	if cType == CultureTypeNomadic && h < 20 {
		return a * 50 // giant sea/lake crossing penalty for Nomads
	}
	if h < 20 {
		return a * 6 // general sea/lake crossing penalty
	}
	if cType == CultureTypeHighland && h < 44 {
		return 3000 // giant penalty for highlanders on lowlands
	}
	if cType == CultureTypeHighland && h < 62 {
		return 200 // giant penalty for highlanders on lowhills
	}
	if cType == CultureTypeHighland {
		return 0 // no penalty for highlanders on highlands
	}
	if h >= 67 {
		return 200 // general mountains crossing penalty
	}
	if h >= 44 {
		return 30 // general hills crossing penalty
	}
	return 0
}

	if r_waterbodies[i] >= 0 {
		if r_waterbody_size[r_waterbodies[i]] > m.mesh.numRegions/25 {
			return "ocean"
		}
		if r_waterbody_size[r_waterbodies[i]] > m.mesh.numRegions/100 {
			return "sea"
		}
		return "gulf"
	}
	  function defineOceanGroup(number) {
	    if (number > grid.cells.i.length / 25) return "ocean";
	    if (number > grid.cells.i.length / 100) return "sea";
	    return "gulf";
	  }

	  function defineIslandGroup(cell, number) {
	    if (cell && features[cells.f[cell - 1]].type === "lake") return "lake_island";
	    if (number > grid.cells.i.length / 10) return "continent";
	    if (number > grid.cells.i.length / 1000) return "island";
	    return "isle";
	  }*/
/*


  // expand cultures across the map (Dijkstra-like algorithm)
  const expand = function () {
    TIME && console.time("expandCultures");
    cells = pack.cells;

    const queue = new PriorityQueue({comparator: (a, b) => a.p - b.p});
    pack.cultures.forEach(function (c) {
      if (!c.i || c.removed) return;
      queue.queue({e: c.center, p: 0, c: c.i});
    });

    const neutral = (cells.i.length / 5000) * 3000 * neutralInput.value; // limit cost for culture growth
    const cost = [];
    while (queue.length) {
      const next = queue.dequeue(),
        n = next.e,
        p = next.p,
        c = next.c;
      const type = pack.cultures[c].type;
      cells.c[n].forEach(function (e) {
        const biome = cells.biome[e];
        const biomeCost = getBiomeCost(c, biome, type);
        const biomeChangeCost = biome === cells.biome[n] ? 0 : 20; // penalty on biome change
        const heightCost = getHeightCost(e, cells.h[e], type);
        const riverCost = getRiverCost(cells.r[e], e, type);
        const typeCost = getTypeCost(cells.t[e], type);
        const totalCost =
          p + (biomeCost + biomeChangeCost + heightCost + riverCost + typeCost) / pack.cultures[c].expansionism;

        if (totalCost > neutral) return;

        if (!cost[e] || totalCost < cost[e]) {
          if (cells.s[e] > 0) cells.culture[e] = c; // assign culture to populated cell
          cost[e] = totalCost;
          queue.queue({e, p: totalCost, c});
        }
      });
    }

    TIME && console.timeEnd("expandCultures");
  };



  function getHeightCost(i, h, type) {
    const f = pack.features[cells.f[i]],
      a = cells.area[i];
    if (type === "Lake" && f.type === "lake") return 10; // no lake crossing penalty for Lake cultures
    if (type === "Naval" && h < 20) return a * 2; // low sea/lake crossing penalty for Naval cultures
    if (type === "Nomadic" && h < 20) return a * 50; // giant sea/lake crossing penalty for Nomads
    if (h < 20) return a * 6; // general sea/lake crossing penalty
    if (type === "Highland" && h < 44) return 3000; // giant penalty for highlanders on lowlands
    if (type === "Highland" && h < 62) return 200; // giant penalty for highlanders on lowhills
    if (type === "Highland") return 0; // no penalty for highlanders on highlands
    if (h >= 67) return 200; // general mountains crossing penalty
    if (h >= 44) return 30; // general hills crossing penalty
    return 0;
  }

  function getRiverCost(r, i, type) {
    if (type === "River") return r ? 0 : 100; // penalty for river cultures
    if (!r) return 0; // no penalty for others if there is no river
    return minmax(cells.fl[i] / 10, 20, 100); // river penalty from 20 to 100 based on flux
  }

  function getTypeCost(t, type) {
    if (t === 1) return type === "Naval" || type === "Lake" ? 0 : type === "Nomadic" ? 60 : 20; // penalty for coastline
    if (t === 2) return type === "Naval" || type === "Nomadic" ? 30 : 0; // low penalty for land level 2 for Navals and nomads
    if (t !== -1) return type === "Naval" || type === "Lake" ? 100 : 0; // penalty for mainland for navals
    return 0;
  }*/
