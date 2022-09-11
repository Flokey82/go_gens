// Package genflora will at some point provide procedurally generated plants for specific biomes
// and climates.
// See: http://link.springer.com/content/pdf/bbm%3A978-94-009-8680-0%2F1.pdf
// And: https://storage.googleapis.com/pirk.io/papers/Makowski.etal-2019-Synthetic-Silviculture.pdf
package genflora

import (
	"github.com/Flokey82/go_gens/utils"
)

// Some interesting values.
var (
	Temperature             = utils.FloatRange{-10, 33}   // Avg. yearly temperature (in °C).
	Precipitation           = utils.FloatRange{10, 4300}  // Avg. yearly precipitation (in mm).
	SeedingFrequency        = utils.IntRange{1, 10}       // Num. of seeds/seeding period. (in 1/year)
	SeedingRadius           = utils.FloatRange{0.01, 100} // Radius of placing plant seeds (in m).
	ShadeTolerance          = utils.FloatRange{0, 1}      // Plant adaptation to shade.
	TemperatureAdaptation   = utils.FloatRange{-10, 33}   // Optimal habitat temperature (in °C).
	PrecipitationAdaptation = utils.FloatRange{10, 4300}  // Optimal habitat precipitation (in mm).
)

type PlantClass int

// TODO: Extend this list
const (
	Tree        PlantClass = iota // Normal tree
	Treelike                      // Plants that look like a tree
	Krummholz                     // Stunted tree types
	Shrub                         // Bushes
	Grass                         // Normal grass as well as grains
	Forb                          // Herbaceous flowering plant
	Herb                          // Undifferentiated small herbs
	Vine                          // Climbing plants
	Fern                          // Ferns
	Epiphyte                      // Plants that grow on other plants
	Thallophyte                   // Lower plants (algea, funghi, moss)
	Succulent                     // Plants that store water in their leaves
)

// plantClassConfig is a configuration for a plant class.
type plantClassConfig struct {
	MaxHeight utils.FloatRange
	MaxAge    utils.FloatRange
	Stem      *PropertyRange
	Leaf      *PropertyRange
	Flower    *PropertyRange
	Fruit     *PropertyRange
	Seed      *PropertyRange
}

// Random returns a random plant lineage for the given class.
func (pcc plantClassConfig) Random() *PlantLineage {
	plant := &PlantLineage{
		MaxHeight: pcc.MaxHeight.Random(),
		MaxAge:    pcc.MaxAge.Random(),
	}

	// Generate plant parts if they are defined.
	if pcc.Stem != nil {
		plant.Stem = pcc.Stem.Random()
	}
	if pcc.Leaf != nil {
		plant.Leaf = pcc.Leaf.Random()
	}
	if pcc.Flower != nil {
		plant.Flower = pcc.Flower.Random()
	}
	if pcc.Fruit != nil {
		plant.Fruit = pcc.Fruit.Random()
	}
	if pcc.Seed != nil {
		plant.Seed = pcc.Seed.Random()
	}
	return plant
}

func GenRandomLineage() *PlantLineage {
	return PlantClassConfigs[Tree].Random()
}

var PlantClassConfigs = map[PlantClass]plantClassConfig{
	Tree: {
		MaxHeight: utils.FloatRange{1, 130},
		MaxAge:    utils.FloatRange{1, 100},
		Stem: &PropertyRange{
			Diameter: utils.FloatRange{1 / 80, 1 / 60},
			Length:   utils.FloatRange{0.75, 1},
			Density:  utils.IntRange{1, 1},
		},
		Leaf: &PropertyRange{
			Diameter: utils.FloatRange{0.01, 0.05},
			Length:   utils.FloatRange{0.01, 0.05},
			Density:  utils.IntRange{1, 10},
		},
		Flower: &PropertyRange{
			Diameter: utils.FloatRange{0.1, 0.5},
			Length:   utils.FloatRange{0.5, 1},
			Density:  utils.IntRange{1, 1},
		},
		Fruit: &PropertyRange{
			Diameter: utils.FloatRange{0.1, 0.5},
			Length:   utils.FloatRange{0.5, 1},
			Density:  utils.IntRange{1, 1},
		},
		Seed: &PropertyRange{
			Diameter: utils.FloatRange{0.1, 0.5},
			Length:   utils.FloatRange{0.5, 1},
			Density:  utils.IntRange{1, 1},
		},
	},
	Treelike: {
		MaxHeight: utils.FloatRange{0.5, 10.0},
		MaxAge:    utils.FloatRange{1, 100},
		Stem: &PropertyRange{
			Diameter: utils.FloatRange{1 / 80, 1 / 60},
			Length:   utils.FloatRange{0.75, 1},
			Density:  utils.IntRange{3, 10},
		},
		Leaf: &PropertyRange{
			Diameter: utils.FloatRange{0.01, 0.05},
			Length:   utils.FloatRange{0.01, 0.05},
			Density:  utils.IntRange{1, 10},
		},
		Flower: &PropertyRange{
			Diameter: utils.FloatRange{0.1, 0.5},
			Length:   utils.FloatRange{0.5, 1},
			Density:  utils.IntRange{1, 1},
		},
		Fruit: &PropertyRange{
			Diameter: utils.FloatRange{0.1, 0.5},
			Length:   utils.FloatRange{0.5, 1},
			Density:  utils.IntRange{10, 20},
		},
		Seed: &PropertyRange{
			Diameter: utils.FloatRange{0.1, 0.5},
			Length:   utils.FloatRange{0.5, 1},
			Density:  utils.IntRange{1, 1},
		},
	},
}
