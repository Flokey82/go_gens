package genworldvoronoi

import (
	"log"
	"math"
	"math/rand"

	"github.com/Flokey82/go_gens/genbiome"
)

func (m *Map) getRCulture(r int) *Culture {
	return m.cultures_r[m.r_cultures[r]]
}

// This code is based on:
// https://github.com/Azgaar/Fantasy-Map-Generator/blob/master/modules/cultures-generator.js
func (m *Map) rPlaceNCultures(n int) {
	m.resetRand()
	m.cultures_r = m.placeNCultures(n)

	var seeds []int
	originToCulture := make(map[int]*Culture)
	for _, c := range m.cultures_r {
		seeds = append(seeds, c.Origin)
		originToCulture[c.Origin] = c
	}
	_, maxElev := minMax(m.r_elevation)
	twf := m.getTerritoryWeightFunc()
	biomeWeight := m.getTerritoryBiomeWeightFunc()
	m.r_cultures = m.rPlaceNTerritoriesCustom(seeds, func(o, u, v int) float64 {
		c := originToCulture[o]
		eleVal := m.r_elevation[v] / maxElev
		gotBiome := m.getRBiomeTEMP(v, eleVal, maxElev)
		biomePenalty := biomeWeight(o, u, v) * c.Type.BiomeAffinity(gotBiome) * float64(genbiome.AzgaarBiomeMovementCost[gotBiome]) / 100
		return biomePenalty + twf(o, u, v)/c.Expansionism
	})

	// Update stats?
}

func (m *Map) placeNCultures(n int) []*Culture {
	// The fitness function, returning a score from
	// 0.0 to 1.0 for a given region.
	var sf func(int) float64

	// The distance seed point function, returning
	// seed points/regions that we want to be far
	// away from.
	var dsf func() []int

	var cultures []*Culture
	rctf := m.getRCultureTypeFunc()
	fc := m.getFitnessClimate()
	sf = func(r int) float64 {
		if m.r_elevation[r] <= 0 {
			return 0
		}
		return math.Sqrt((fc(r) + 3.0) / 4.0)
	}

	// For now we just maximize the distance to cities of the same type.
	dsf = func() []int {
		var cultureSeeds []int
		for _, c := range cultures {
			//if c.Type == cType {
			cultureSeeds = append(cultureSeeds, c.Origin)
			//}
		}
		return cultureSeeds
	}
	// Place n cities of the given type.
	for i := 0; i < n; i++ {
		c := m.placeCulture(rctf, sf, dsf)
		log.Printf("placing culture %d: %s", i, c.Name)
		cultures = append(cultures, c)
	}
	return cultures
}

// Culture represents a culture.
type Culture struct {
	ID           int // unique ID
	Origin       int // Same as ID... Remove
	Name         string
	Language     *Language
	Type         CultureType
	Expansionism float64
	// Extremism float64
	// Parent    *Culture
	// Children  []*Culture
	// Extinct   bool
}

func (m *Map) placeCulture(rctf func(int) CultureType, sf func(int) float64, distSeedFunc func() []int) *Culture {
	// Score all regions, pick highest score.
	var newculture int
	lastMax := math.Inf(-1)
	for i, val := range m.rCityScore(sf, distSeedFunc) {
		if val > lastMax {
			newculture = i
			lastMax = val
		}
	}
	cultureType := rctf(newculture)
	lang := GenLanguage(m.seed + int64(newculture))
	return &Culture{
		ID:           newculture,
		Origin:       newculture,
		Name:         lang.MakeName(),
		Type:         cultureType,
		Expansionism: cultureType.Expansionism(),
		Language:     lang,
	}
}

type CultureType int

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

func (t CultureType) BiomeAffinity(biome int) float64 {
	if t == CultureTypeHunting {
		return 5 // non-native biome penalty for hunters
	}
	if t == CultureTypeNomadic && biome > 4 && biome < 10 {
		return 10 // forest biome penalty for nomads
	}
	return 2 // general non-native biome penalty
}

// round value to d decimals
func roundToDecimals(v, d float64) float64 {
	m := math.Pow(10, d)
	return math.Round(v*m) / m
}

func (m *Map) getRHaven(i int) (int, int) {
	// get all neighbors that are below or at sea level.
	var water []int
	for _, nb := range m.rNeighbors(i) {
		if m.r_elevation[nb] <= 0.0 {
			water = append(water, nb)
		}
	}
	// Get distances of i to each water neighbor.
	// get the closest water neighbor.
	iLatLon := m.r_latLon[i]
	closest := -1
	var minDist float64
	for _, nb := range water {
		nbLatLon := m.r_latLon[nb]
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

func (m *Map) getRCellTypes() []int {
	var ocean_r, land_r []int
	for r, elev := range m.r_elevation {
		if elev <= 0.0 {
			ocean_r = append(ocean_r, r)
		} else {
			land_r = append(land_r, r)
		}
	}
	r_distance_ocean := m.assignDistanceField(ocean_r, make(map[int]bool))
	r_distance_land := m.assignDistanceField(land_r, make(map[int]bool))

	cellType := make([]int, m.mesh.numRegions)
	for i := range cellType {
		// Is it water?
		if m.r_elevation[i] <= 0.0 {
			// Figure out if it has a land neighbor.
			// If so, it is -1 (water near coast)
			if r_distance_land[i] <= 1 {
				cellType[i] = -1
			} else {
				// If not, it is -2 (water far from coast)
				cellType[i] = -2
			}
		} else {
			// Figure out if it has a water neighbor.
			// If so, it is 1 (land near coast)
			if r_distance_ocean[i] <= 1 {
				cellType[i] = 1
			} else {
				// If not, it is >=2 (land far from coast)
				cellType[i] = int(r_distance_ocean[i])
			}
		}
	}
	return cellType
}

func (m *Map) getRFeatureTypeFunc() func(int) string {
	r_waterbodies := m.getWaterBodies()
	r_waterbody_size := m.getWaterBodySizes()
	r_landmasses := m.identifyLandmasses()
	r_landmass_size := m.getLandmassSizes()

	return func(i int) string {
		if i >= 0 {
			if r_waterbodies[i] >= 0 {
				if r_waterbody_size[r_waterbodies[i]] > m.mesh.numRegions/25 {
					return "ocean"
				}
				if r_waterbody_size[r_waterbodies[i]] > m.mesh.numRegions/100 {
					return "sea"
				}
				if r_waterbody_size[r_waterbodies[i]] > m.mesh.numRegions/500 {
					return "gulf"
				}
				return "lake"
			}
			if r_landmasses[i] >= 0 {
				if r_landmass_size[r_landmasses[i]] < m.mesh.numRegions/100 {
					return "isle"
				}
				return "continent"
			}
		}
		return ""
	}
}

func (m *Map) getRBiomeTEMP(r int, elev, maxElev float64) int {
	return genbiome.GetAzgaarBiome(int(20.0*m.r_moisture[r]), int(m.getRTemperature(r, maxElev)), int(elev*100))
}

func (m *Map) getRCultureTypeFunc() func(int) CultureType {
	cellType := m.getRCellTypes()
	getType := m.getRFeatureTypeFunc()
	wmf := m.getRWhittakerModBiomeFunc()

	r_waterbody_size := m.getWaterBodySizes()
	_, maxElev := minMax(m.r_elevation)
	// set culture type based on culture center position
	return func(r int) CultureType {
		eleVal := m.r_elevation[r] / maxElev
		gotBiome := m.getRBiomeTEMP(r, eleVal, maxElev)
		log.Println(gotBiome)
		log.Println(wmf(r))
		if eleVal < 0.7 {
			if gotBiome == genbiome.AzgaarBiomeHotDesert ||
				gotBiome == genbiome.AzgaarBiomeColdDesert ||
				gotBiome == genbiome.AzgaarBiomeGrassland {
				return CultureTypeNomadic // high penalty in forest biomes and near coastline
			}
		}
		if eleVal > 0.3 {
			return CultureTypeHighland // no penalty for hills and moutains, high for other elevations
		}
		haven, harbor := m.getRHaven(r)
		havenType := getType(haven) // opposite feature
		rType := getType(r)
		log.Println(havenType, rType)
		if havenType == "lake" && r_waterbody_size[haven] > 5 {
			// Ensure larger lakes.
			return CultureTypeLake // low water cross penalty and high for growth not along coastline
		}
		if (harbor > 0 && havenType != "lake" && P(0.1)) ||
			(harbor == 1 && P(0.6)) ||
			(rType == "isle" && P(0.4)) {
			return CultureTypeNaval // low water cross penalty and high for non-along-coastline growth
		}
		if m.isBigRiver(r) { // cells.r[r] && cells.fl[r] > 100
			return CultureTypeRiver // no River cross penalty, penalty for non-River growth
		}
		// Probably inland? Distance to ocean?
		if cellType[r] > 2 {
			if gotBiome == genbiome.AzgaarBiomeSavanna ||
				gotBiome == genbiome.AzgaarBiomeTropicalRainforest ||
				gotBiome == genbiome.AzgaarBiomeTemperateRainforest ||
				gotBiome == genbiome.AzgaarBiomeTaiga ||
				gotBiome == genbiome.AzgaarBiomeTundra ||
				gotBiome == genbiome.AzgaarBiomeWetland {
				return CultureTypeHunting // high penalty in non-native biomes
			}
		}
		// TODO: Wildlands
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

func (m *Map) getCellTypeCost(t int, cType CultureType) int {
	if t == 1 {
		if cType == CultureTypeNaval || cType == CultureTypeLake {
			return 0
		}
		if cType == CultureTypeNomadic {
			return 60
		}
		return 20 // penalty for coastline
	}
	if t == 2 {
		if cType == CultureTypeNaval || cType == CultureTypeNomadic {
			return 30
		}
		return 0 // low penalty for land level 2 for Navals and nomads
	}
	if t != -1 {
		if cType == CultureTypeNaval || cType == CultureTypeLake {
			return 100
		}
		return 0 // penalty for mainland for navals
	}
	return 0
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

// probability shorthand
func P(probability float64) bool {
	if probability >= 1.0 {
		return true
	}
	if probability <= 0 {
		return false
	}
	return rand.Float64() < probability
}
