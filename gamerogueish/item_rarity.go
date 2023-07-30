package gamerogueish

import "math/rand"

// Rarity represents the rarity of an item.
type Rarity struct {
	Name           string // Name of this rarity
	Probability    int    // Probability of this rarity (the higher the more rare)
	IndicateRarity bool   // Indicate rarity in item name
}

// Roll returns true if the item should be generated.
func (r *Rarity) Roll() bool {
	return rand.Intn(101) >= r.Probability
}

var (
	RarityAbundant = &Rarity{
		Name:           "abundant",
		Probability:    25,
		IndicateRarity: false,
	}
	RarityCommon = &Rarity{
		Name:           "common",
		Probability:    45,
		IndicateRarity: false,
	}
	RarityAverage = &Rarity{
		Name:           "average",
		Probability:    65,
		IndicateRarity: false,
	}
	RarityUncommon = &Rarity{
		Name:           "uncommon",
		Probability:    80,
		IndicateRarity: true,
	}
	RarityRare = &Rarity{
		Name:           "rare",
		Probability:    93,
		IndicateRarity: true,
	}
	RarityExotic = &Rarity{
		Name:           "exotic",
		Probability:    99,
		IndicateRarity: true,
	}
	RarityLegendary = &Rarity{
		Name:           "legendary",
		Probability:    100,
		IndicateRarity: true,
	}
)
