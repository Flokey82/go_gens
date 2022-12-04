package genworldvoronoi

// Bio handles the generation of life on the map (plants, animals, etc.).
type Bio struct {
	*Geo
	Species    []*Species // All species on the map.
	NumSpecies int        // Number of species to generate.
}

func newBio(geo *Geo) *Bio {
	return &Bio{
		Geo:        geo,
		NumSpecies: 100,
	}
}

func (b *Bio) generateBiology() {
	b.genNRandomSpecies(b.NumSpecies) // Generate the species.
}

func (b *Bio) Tick() {
	// TODO: Tick the species.
}
