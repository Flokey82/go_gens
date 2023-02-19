package genreligion

import "github.com/Flokey82/go_gens/genlanguage"

var GenBase = genlanguage.GenBase

var GenBaseAdjective = genlanguage.GenBaseAdjective
var GenBaseAnimal = genlanguage.GenBaseAnimal
var GenBaseBeing = genlanguage.GenBaseBeing
var GenBaseColor = genlanguage.GenBaseColor
var GenBaseNumber = genlanguage.GenBaseNumber
var GenBaseGenitive = genlanguage.GenBaseGenitive
var GenBaseTheGenitive = genlanguage.GenBaseTheGenitive

// RandDeityMeaningGenMethod returns a random deity meaning generation method.
func (g *Generator) RandDeityGenMethod() string {
	return g.Ra(DeityMeaningApproaches)
}

// GenerateDeityMeaning generates a meaning for a deity name.
// This code is based on:
// https://github.com/Azgaar/Fantasy-Map-Generator/blob/master/modules/religions-generator.js
func (g *Generator) GenerateDeityMeaning(approach string) string {
	switch approach { // select generation approach
	case ApproachNumber:
		return g.Ra(GenBase[GenBaseNumber])
	case ApproachBeing:
		return g.Ra(GenBase[GenBaseBeing])
	case ApproachAdjective:
		return g.Ra(GenBase[GenBaseAdjective])
	case ApproachColorAnimal:
		return g.Ra(GenBase[GenBaseColor]) + " " + g.Ra(GenBase[GenBaseAnimal])
	case ApproachAdjectiveAnimal:
		return g.Ra(GenBase[GenBaseAdjective]) + " " + g.Ra(GenBase[GenBaseAnimal])
	case ApproachAdjectiveBeing:
		return g.Ra(GenBase[GenBaseAdjective]) + " " + g.Ra(GenBase[GenBaseBeing])
	case ApproachAdjectiveGenitive:
		return g.Ra(GenBase[GenBaseAdjective]) + " " + g.Ra(GenBase[GenBaseGenitive])
	case ApproachColorBeing:
		return g.Ra(GenBase[GenBaseColor]) + " " + g.Ra(GenBase[GenBaseBeing])
	case ApproachColorGenitive:
		return g.Ra(GenBase[GenBaseColor]) + " " + g.Ra(GenBase[GenBaseGenitive])
	case ApproachBeingOfGenitive:
		return g.Ra(GenBase[GenBaseBeing]) + " of " + g.Ra(GenBase[GenBaseGenitive])
	case ApproachBeingOfTheGenitive:
		return g.Ra(GenBase[GenBaseBeing]) + " of the " + g.Ra(GenBase[GenBaseTheGenitive])
	case ApproachAnimalOfGenitive:
		return g.Ra(GenBase[GenBaseAnimal]) + " of " + g.Ra(GenBase[GenBaseGenitive])
	case ApproachAdjectiveBeingOfGenitive:
		return g.Ra(GenBase[GenBaseAdjective]) + " " + g.Ra(GenBase[GenBaseBeing]) + " of " + g.Ra(GenBase[GenBaseGenitive])
	case ApproachAdjectiveAnimalOfGenitive:
		return g.Ra(GenBase[GenBaseAdjective]) + " " + g.Ra(GenBase[GenBaseAnimal]) + " of " + g.Ra(GenBase[GenBaseGenitive])
	default:
		return "ERROR"
	}
}

// GenMeaningApproaches contains a map of name generation
// approaches and their relative chance to be selected.
var GenMeaningApproaches = map[string]int{
	ApproachNumber:                    1,
	ApproachBeing:                     3,
	ApproachAdjective:                 5,
	ApproachColorAnimal:               5,
	ApproachAdjectiveAnimal:           5,
	ApproachAdjectiveBeing:            5,
	ApproachAdjectiveGenitive:         1,
	ApproachColorBeing:                3,
	ApproachColorGenitive:             3,
	ApproachBeingOfGenitive:           2,
	ApproachBeingOfTheGenitive:        1,
	ApproachAnimalOfGenitive:          1,
	ApproachAdjectiveBeingOfGenitive:  2,
	ApproachAdjectiveAnimalOfGenitive: 2,
}

const (
	ApproachNumber                    = "Number"
	ApproachBeing                     = "Being"
	ApproachAdjective                 = "Adjective"
	ApproachColorAnimal               = "Color + Animal"
	ApproachAdjectiveAnimal           = "Adjective + Animal"
	ApproachAdjectiveBeing            = "Adjective + Being"
	ApproachAdjectiveGenitive         = "Adjective + Genitive"
	ApproachColorBeing                = "Color + Being"
	ApproachColorGenitive             = "Color + Genitive"
	ApproachBeingOfGenitive           = "Being + of + Genitive"
	ApproachBeingOfTheGenitive        = "Being + of the + Genitive"
	ApproachAnimalOfGenitive          = "Animal + of + Genitive"
	ApproachAdjectiveBeingOfGenitive  = "Adjective + Being + of + Genitive"
	ApproachAdjectiveAnimalOfGenitive = "Adjective + Animal + of + Genitive"
)
