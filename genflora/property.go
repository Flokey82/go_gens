package genflora

import (
	"math"
	"math/rand"

	"github.com/Flokey82/go_gens/utils"
)

// PropertyRange represents a range of properties for a plant part.
//
// NOTE: The values are considered fractions of the current plant height.
// In order to get the actual values, multiply the fractions with the plant height.
// To get individual values, divide by the density. (?)
//
// Well, I haven't really figured out yet how to work out the best values.
//
// TODO: Add frequency and density (how many parts per unit area)
//   - High fruit density could mean a bunch of grapes, while low density could mean a single apple.
//   - Frequency defines how often the plant part occurs per plant height.
//   - Figure out if the ratios make sense and how things are being calculated.
type PropertyRange struct {
	Length    utils.FloatRange // Length of the plant part
	Diameter  utils.FloatRange // Diameter of the plant part
	Frequency utils.FloatRange // Frequency of the plant part
	Density   utils.IntRange   // Density of the plant part
	Shapes    []int            // List of possible shapes
	Effects   []int            // List of possible effects (e.g. poison, healing, etc.)
}

// Random returns a random property from the range.
func (p PropertyRange) Random() *Property {
	prop := &Property{
		Length:    p.Length.Random(),
		Diameter:  p.Diameter.Random(),
		Frequency: p.Frequency.Random(),
		Density:   p.Density.Random(),
	}
	if len(p.Shapes) > 0 {
		prop.Shape = p.Shapes[rand.Intn(len(p.Shapes))]
	}
	if len(p.Effects) > 0 {
		prop.Effect = p.Effects[rand.Intn(len(p.Effects))]
	}
	return prop
}

// Property represents properties of a plant part.
// TODO: Add frequency and density (how many parts per unit area)
//   - High fruit density could mean a bunch of grapes, while low density could mean a single apple.
//   - Frequency defines how often the plant part occurs per plant height.
type Property struct {
	Length    float64 // Length of the plant part
	Diameter  float64 // Diameter of the plant part
	Frequency float64 // Frequency of the plant part
	Density   int     // Density of the plant part
	Shape     int     // Shape of the plant part
	Effect    int     // Effect of the plant part (e.g. poison, healing, etc.)
}

// Area returns the cross-section area of the plant part.
func (p Property) Area() float64 {
	return math.Pi * math.Pow(p.Diameter/2, 2)
}

// Volume returns the rough volume of the plant part.
func (p Property) Volume() float64 {
	return p.Area() * p.Length
}

const (
	ShapeSucculent = iota // Fleshy leaves/stems
	ShapeWoody            // Hard, wooden
)
