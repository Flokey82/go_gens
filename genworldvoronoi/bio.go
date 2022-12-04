package genworldvoronoi

// Bio handles the generation of life on the map (plants, animals, etc.).
type Bio struct {
	*Geo
	Species    []*Species // All species on the map.
	GrowthDays []int      // Number of days within the growth period for each region.
	NumSpecies int        // Number of species to generate.
}

func newBio(geo *Geo) *Bio {
	return &Bio{
		Geo:        geo,
		GrowthDays: make([]int, geo.mesh.numRegions),
		NumSpecies: 100,
	}
}

func (b *Bio) generateBiology() {
	// Calculate the duration of the potential growth period for each region.
	b.calcGrowthPeriod()

	// TODO: Calculate a score for each region that reflects how well
	// suited it is for agriculture during the growth period. This
	// would be based on insolation, temperature, rainfall, steepness,
	// etc.

	// TODO: Calculate a score for each region that reflects how well
	// herbivores would be able to survive there.
	// As long as there are shrubs, etc. then herbivores should be able
	// to survive even in winter.

	// TODO: Calculate a score for each region that reflects how well
	// carnivores would be able to survive there (presence of prey animals,
	// total survivability).

	// Generate the species.
	b.genNRandomSpecies(b.NumSpecies)
}

// calcGrowthPeriod calculates the duration of the potential growth
// period for each region (dormancy can be inferred), which will
// give us the potential for agricultural output (nr of harvests etc).
func (b *Bio) calcGrowthPeriod() {
	for r := range b.GrowthDays {
		var growthDays int
		for i := 0; i < 356; i++ {
			// Calculate daily average temperature.
			min, max := b.getMinMaxTemperatureOfDay(b.LatLon[r][0], i)
			avg := (min + max) / 2

			// TODO: Right now we only count days where the average temperature
			// is above 0. This is not correct, as we should be counting days
			// where the average temperature is above a certain minimum.
			// We should also take in account when there is precipitation.
			if avg > 0 && b.Rainfall[r] > 0 {
				growthDays++
			}
		}
		b.GrowthDays[r] = growthDays
	}
}

func (b *Bio) Tick() {
	// TODO: Tick the species.
}
