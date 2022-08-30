package genflora

import (
	"fmt"
	"math/rand"
)

type PlantLineage struct {
	MaxHeight float64
	MaxAge    float64   // Life expectency
	Stem      *Property // Stem properties
	Leaf      *Property // Leaf properties

	Flower *Property // Flower properties
	Fruit  *Property // Fruit properties
	Seed   *Property // Seed properties
	// Environmental conditions
	//
	// Min / max
	//  - temperature
	//  - moisture
	//  - soil depth
	//  - slope
	//  - altitude
	//  - precipitation

	// Lifespan
	// Growth rate
	// Area required

	// TODO: Synthesize L-systems from plant types.
	// Appearance
	//
	// Leaf (*)
	//   - frequency
	//   - density
	//   - diameter / length / size
	//   - shape
	//
	// Stem (*)
	//   - frequency
	//   - density
	//   - diameter / length / size
	//   - shape
	//
	// Flower (*)
	//   - frequency
	//   - density
	//   - diameter / length / size
	//   - shape
	//
	// Fruit / Seed (*)
	//   - frequency
	//   - density
	//   - diameter / length / size
	//   - shape
	//
	// Plant
	//   - diameter / length / size
	//   - shape
	//
	// Effect
	//   - poision (none, paralyze, diarrhea, etc)
	//   - magic (none, heal, cure, etc)
	//
	// (*) = optional
	// - frequency = the number of ocurrences along the
}

func (p PlantLineage) Mutate() *PlantLineage {
	// TODO: Add age, height etc to the mutateable properties.

	// Pick a random part of the plant to mutate.
	var plantParts []*Property
	if p.Stem != nil {
		plantParts = append(plantParts, p.Stem)
	}
	if p.Leaf != nil {
		plantParts = append(plantParts, p.Leaf)
	}
	if p.Flower != nil {
		plantParts = append(plantParts, p.Flower)
	}
	if p.Fruit != nil {
		plantParts = append(plantParts, p.Fruit)
	}
	if p.Seed != nil {
		plantParts = append(plantParts, p.Seed)
	}

	// Mutate the part.
	part := plantParts[rand.Intn(len(plantParts))]
	part.Length += rand.Float64() * 0.1
	part.Diameter += rand.Float64() * 0.1
	part.Shape = rand.Intn(10) // TODO: Make sure we don't exceed the shape range
	return &p
}

func (p PlantLineage) String() string {
	return fmt.Sprintf("height: %f, age: %f, water: %f", p.MaxHeight, p.MaxAge, p.getWaterRequirement())
}

func (p PlantLineage) Crossover(other PlantLineage) *PlantLineage {
	// TODO: Pick random parts of the plant to crossover, creating a hybrid.
	return &p
}

// getWaterRequirement returns the base amount of water required by the plant (per height).
// NOTE: The following code is just experimental. I'll have to think this through in more detail.
// TODO: Take in account frequency and density of each plant part.
func (p PlantLineage) getWaterRequirement() float64 {
	var waterRequirement float64
	// If leaves are present, use them to calculate the water requirement.
	// - Smaller leaves require less water.
	// - However thick leaves (like succulents) can store more water.
	if p.Leaf != nil {
		// Roughly calculate leaf area.
		leafArea := p.Leaf.Area() * float64(p.Leaf.Density)

		// Succulent leaves require less water.
		if p.Leaf.Shape == ShapeSucculent {
			waterRequirement = leafArea * 0.01
		} else {
			waterRequirement = leafArea * 0.1
		}
	}

	// If stems are present, use them to calculate the water requirement.
	// Thicker stems can store more water.
	if p.Stem != nil {
		// Roughly calculate stem volume.
		stemVolume := p.Stem.Volume() * float64(p.Stem.Density)

		// Succulent stems (like cacti) require less water.
		if p.Stem.Shape == ShapeSucculent {
			waterRequirement += stemVolume * 0.01
		} else {
			waterRequirement += stemVolume * 0.1
		}
	}

	// If we need to produce fruit, add the water requirement of the fruit.
	if p.Fruit != nil {
		// Roughly calculate fruit volume.
		fruitVolume := p.Fruit.Volume() * float64(p.Fruit.Density)
		waterRequirement += fruitVolume * 0.01
	}

	// If we need to produce seeds, add the water requirement of the seeds.
	if p.Seed != nil {
		// Roughly calculate seed volume.
		seedVolume := p.Seed.Volume() * float64(p.Seed.Density)
		waterRequirement += seedVolume * 0.01
	}

	// If we need to produce flowers, add the water requirement of the flowers.
	if p.Flower != nil {
		// Roughly calculate flower area.
		flowerArea := p.Flower.Area() * float64(p.Flower.Density)
		waterRequirement += flowerArea * 0.01
	}

	// Larger plants require more water.
	// NOTE: We return instead the base-value per height, so that we can calculate the actual
	// water requirement based on the current height of the plant.
	// This would also require to take in account the age of the plant
	// as plants only produce fruit, seeds, and flowers after a certain age.
	waterRequirement *= p.MaxHeight
	return waterRequirement
}

// getTemperatureRequirement returns the temperature requirement of the plant.
//
// TODO: This should be a mix of averaging and modifying the requirement
// Fruits require a certain temperature range while succulent plant parts
// help the plant to survive under more arid conditions.
func (p PlantLineage) getTemperatureRequirement() float64 {
	// Start out with a base temperature requirement that we will modify
	// based on the plant parts.
	temperatureRequirement := 20.0

	return temperatureRequirement
}
