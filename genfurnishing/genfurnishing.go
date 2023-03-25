package genfurnishing

type FurnishingQuantity struct {
	Name        string
	Min, Max    int
	Probability float64
}

var (
	QuantityNone = &FurnishingQuantity{
		Name:        "none",
		Min:         0,
		Max:         0,
		Probability: 0.25,
	}
	QuantityMinimum = &FurnishingQuantity{
		Name:        "minimum",
		Min:         1,
		Max:         1,
		Probability: 0.75,
	}
	QuantitySparse = &FurnishingQuantity{
		Name:        "sparse",
		Min:         2,
		Max:         2,
		Probability: 0.92,
	}
	QuantityAverage = &FurnishingQuantity{
		Name:        "average",
		Min:         3,
		Max:         4,
		Probability: 0.98,
	}
	QuantityFurnished = &FurnishingQuantity{
		Name:        "furnished",
		Min:         5,
		Max:         6,
		Probability: 1.0,
	}
)
