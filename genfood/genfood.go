package genfood

import (
	"errors"
	"math/rand"

	"github.com/Flokey82/go_gens/genstory"
)

func randomString(strs []string) (string, error) {
	if len(strs) == 0 {
		return "", errors.New("no strings provided")
	}
	return strs[rand.Intn(len(strs))], nil
}

func slicesStringIn(str string, strs []string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

var vegetables = []string{
	"artichoke",
	"asparagus",
	"aubergine",
	"beet",
}

var fruits = []string{
	"apple",
	"apricot",
	"avocado",
	"banana",
	"blackberry",
	"blueberry",
}

var meats = []string{
	"beef",
	"chicken",
	"duck",
	"fish",
	"lamb",
	"pork",
	"turkey",
}

var grains = []string{
	"barley",
	"bread",
}

var dairy = []string{
	"butter",
	"cheese",
}

var spices = []string{
	"pepper",
	"paprika",
}

var herbs = []string{
	"basil",
	"bay",
	"chives",
	"cilantro",
	"coriander",
	"dill",
}

var condiments = []string{
	"butter",
	"cheese",
	"cream",
	"egg",
	"fish",
	"flour",
	"garlic",
	"ginger",
	"milk",
	"oil",
}

var flavorProfiles = []string{
	"aromatic",
	"bitter",
	"pungent",
	"salty",
	"savory",
	"sour",
	"spicy",
	"sweet",
}

var textureProfiles = []string{
	"crunchy",
	"creamy",
	"crispy",
	"juicy",
	"mushy",
	"oily",
	"soft",
	"sticky",
	"stringy",
	"tender",
	"thick",
	"thin",
	"tough",
	"watery",
	"spongy",
}

const (
	TokenVegetable         = "[vegetable]"
	TokenFruit             = "[fruit]"
	TokenMeat              = "[meat]"
	TokenGrain             = "[grain]"
	TokenDairy             = "[dairy]"
	TokenSpice             = "[spice]"
	TokenHerb              = "[herb]"
	TokenCondiment         = "[condiment]"
	TokenFlavorProfile     = "[flavor_profile]"
	TokenTextureProfile    = "[texture_profile]"
	TokenPreparationMethod = "[preparation_method]"
)

var tokens = []string{
	TokenVegetable,
	TokenFruit,
	TokenMeat,
	TokenGrain,
	TokenDairy,
	TokenSpice,
	TokenHerb,
	TokenCondiment,
	TokenFlavorProfile,
	TokenTextureProfile,
	TokenPreparationMethod,
}

var dishStrategies = []string{
	"[preparation_method] [vegetable] with [flavor_profile] [meat]",
	"[preparation_method] [vegetable] with [flavor_profile] [meat] and [flavor_profile] [vegetable]",
	"[preparation_method] [vegetable] with [flavor_profile] [meat] and [flavor_profile] [fruit]",
	"[preparation_method] [fruit] with [texture_profile] [meat] and [herb]-[dairy]",
}

var FoodTextConfig = &genstory.TextConfig{
	TokenPools: map[string][]string{
		TokenVegetable:         vegetables,
		TokenFruit:             fruits,
		TokenMeat:              meats,
		TokenGrain:             grains,
		TokenDairy:             dairy,
		TokenSpice:             spices,
		TokenHerb:              herbs,
		TokenCondiment:         condiments,
		TokenFlavorProfile:     flavorProfiles,
		TokenTextureProfile:    textureProfiles,
		TokenPreparationMethod: preparationMethods,
	},
	TokenIsMandatory: map[string]bool{},
	Tokens:           tokens,
	Templates:        dishStrategies,
	UseAllProvided:   true,
}
