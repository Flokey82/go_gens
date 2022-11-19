// Package genworldvoronoi is a port of redblobgames' amazing planet generator.
// See: https://www.redblobgames.com/x/1843-planet-generation
// And: https://github.com/redblobgames/1843-planet-generation
package genworldvoronoi

// ugh globals, sorry
type Map struct {
	*Geo // Geography / geology
	*Civ // Civilization
}

func NewMap(seed int64, numPlates, numPoints int, jitter float64) (*Map, error) {
	// Initialize the planet.
	geo, err := newGeo(seed, numPlates, numPoints, jitter)
	if err != nil {
		return nil, err
	}

	// Initialize the map.
	m := &Map{
		Geo: geo,
		Civ: NewCiv(geo),
	}
	m.generateMap()
	return m, nil
}

func (m *Map) generateMap() {
	// Build geography / geology / climate.
	m.generateGeology()

	// Build civilization.
	m.generateCivilization()
}
