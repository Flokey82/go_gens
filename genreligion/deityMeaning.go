package genreligion

import (
	"github.com/Flokey82/go_gens/genlanguage"
	"github.com/Flokey82/go_gens/genstory"
)

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
	case TemplateMeaningNumber:
		return g.Ra(GenBase[GenBaseNumber])
	case TemplateMeaningBeing:
		return g.Ra(GenBase[GenBaseBeing])
	case TemplateMeaningAdjective:
		return g.Ra(GenBase[GenBaseAdjective])
	case TemplateMeaningColorAnimal:
		return g.Ra(GenBase[GenBaseColor]) + " " + g.Ra(GenBase[GenBaseAnimal])
	case TemplateMeaningAdjectiveAnimal:
		return g.Ra(GenBase[GenBaseAdjective]) + " " + g.Ra(GenBase[GenBaseAnimal])
	case TemplateMeaningAdjectiveBeing:
		return g.Ra(GenBase[GenBaseAdjective]) + " " + g.Ra(GenBase[GenBaseBeing])
	case TemplateMeaningAdjectiveGenitive:
		return g.Ra(GenBase[GenBaseAdjective]) + " " + g.Ra(GenBase[GenBaseGenitive])
	case TemplateMeaningColorBeing:
		return g.Ra(GenBase[GenBaseColor]) + " " + g.Ra(GenBase[GenBaseBeing])
	case TemplateMeaningColorGenitive:
		return g.Ra(GenBase[GenBaseColor]) + " " + g.Ra(GenBase[GenBaseGenitive])
	case TemplateMeaningBeingOfGenitive:
		return g.Ra(GenBase[GenBaseBeing]) + " of " + g.Ra(GenBase[GenBaseGenitive])
	case TemplateMeaningBeingOfTheGenitive:
		return g.Ra(GenBase[GenBaseBeing]) + " of the " + g.Ra(GenBase[GenBaseTheGenitive])
	case TemplateMeaningAnimalOfGenitive:
		return g.Ra(GenBase[GenBaseAnimal]) + " of " + g.Ra(GenBase[GenBaseGenitive])
	case TemplateMeaningAdjectiveBeingOfGenitive:
		return g.Ra(GenBase[GenBaseAdjective]) + " " + g.Ra(GenBase[GenBaseBeing]) + " of " + g.Ra(GenBase[GenBaseGenitive])
	case TemplateMeaningAdjectiveAnimalOfGenitive:
		return g.Ra(GenBase[GenBaseAdjective]) + " " + g.Ra(GenBase[GenBaseAnimal]) + " of " + g.Ra(GenBase[GenBaseGenitive])
	default:
		return "ERROR"
	}
}

// GenerateDeityMeaningV2 generates a meaning for a deity name.
func (g *Generator) GenerateDeityMeaningV2(provided []genstory.TokenReplacement, approach string) (string, error) {
	return approachConfig.GenerateWithTemplate(provided, approach)
}

// GenMeaningApproaches contains a map of name generation
// approaches and their relative chance to be selected.
var GenMeaningApproaches = map[string]int{
	TemplateMeaningNumber:                    1,
	TemplateMeaningBeing:                     3,
	TemplateMeaningAdjective:                 5,
	TemplateMeaningColorAnimal:               5,
	TemplateMeaningAdjectiveAnimal:           5,
	TemplateMeaningAdjectiveBeing:            5,
	TemplateMeaningAdjectiveGenitive:         1,
	TemplateMeaningColorBeing:                3,
	TemplateMeaningColorGenitive:             3,
	TemplateMeaningBeingOfGenitive:           2,
	TemplateMeaningBeingOfTheGenitive:        1,
	TemplateMeaningAnimalOfGenitive:          1,
	TemplateMeaningAdjectiveBeingOfGenitive:  2,
	TemplateMeaningAdjectiveAnimalOfGenitive: 2,
}

const (
	TokenNumber    = "[NUMBER]"
	TokenBeing     = "[BEING]"
	TokenAdjective = "[ADJECTIVE]"
	TokenColor     = "[COLOR]"
	TokenAnimal    = "[ANIMAL]"
	TokenGenitive  = "[GENITIVE]"
)

const (
	TemplateMeaningNumber                    = "[NUMBER]"
	TemplateMeaningBeing                     = "[BEING]"
	TemplateMeaningAdjective                 = "[ADJECTIVE]"
	TemplateMeaningColorAnimal               = "[COLOR] [ANIMAL]"
	TemplateMeaningAdjectiveAnimal           = "[ADJECTIVE] [ANIMAL]"
	TemplateMeaningAdjectiveBeing            = "[ADJECTIVE] [BEING]"
	TemplateMeaningAdjectiveGenitive         = "[ADJECTIVE] [GENITIVE]"
	TemplateMeaningColorBeing                = "[COLOR] [BEING]"
	TemplateMeaningColorGenitive             = "[COLOR] [GENITIVE]"
	TemplateMeaningBeingOfGenitive           = "[BEING] of [GENITIVE]"
	TemplateMeaningBeingOfTheGenitive        = "[BEING] of the [GENITIVE]"
	TemplateMeaningAnimalOfGenitive          = "[ANIMAL] of [GENITIVE]"
	TemplateMeaningAdjectiveBeingOfGenitive  = "[ADJECTIVE] [BEING] of [GENITIVE]"
	TemplateMeaningAdjectiveAnimalOfGenitive = "[ADJECTIVE] [ANIMAL] of [GENITIVE]"
)

var approachConfig = &genstory.TextConfig{
	TokenPools: map[string][]string{
		TokenNumber:    GenBase[GenBaseNumber],
		TokenBeing:     GenBase[GenBaseBeing],
		TokenAdjective: GenBase[GenBaseAdjective],
		TokenColor:     GenBase[GenBaseColor],
		TokenAnimal:    GenBase[GenBaseAnimal],
		TokenGenitive:  GenBase[GenBaseGenitive],
	},
	TokenIsMandatory: map[string]bool{},
	Tokens: []string{
		TokenNumber,
		TokenBeing,
		TokenAdjective,
		TokenColor,
		TokenAnimal,
		TokenGenitive,
	},
	Templates: []string{
		TemplateMeaningNumber,
		TemplateMeaningBeing,
		TemplateMeaningAdjective,
		TemplateMeaningColorAnimal,
		TemplateMeaningAdjectiveAnimal,
		TemplateMeaningAdjectiveBeing,
		TemplateMeaningAdjectiveGenitive,
		TemplateMeaningColorBeing,
		TemplateMeaningColorGenitive,
		TemplateMeaningBeingOfGenitive,
		TemplateMeaningBeingOfTheGenitive,
		TemplateMeaningAnimalOfGenitive,
		TemplateMeaningAdjectiveBeingOfGenitive,
		TemplateMeaningAdjectiveAnimalOfGenitive,
	},
}
