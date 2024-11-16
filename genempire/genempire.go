// Package genempire provides a generator for empires (right now just names).
// It is based on the generator for religions in genreligion.
package genempire

import (
	"math/rand"

	"github.com/Flokey82/go_gens/genlanguage"
	"github.com/Flokey82/go_gens/genstory"
)

// Generator is a generator for empires.
type Generator struct {
	rng    *rand.Rand
	txtGen *genstory.TextGenerator
	lang   *genlanguage.Language
}

// NewGenerator creates a new religion generator.
func NewGenerator(seed int64, lang *genlanguage.Language) *Generator {
	rng := rand.New(rand.NewSource(seed))
	if lang == nil {
		lang = genlanguage.GenLanguage(seed)
	}
	return &Generator{
		rng:    rng,
		txtGen: genstory.NewTextGenerator(rng),
		lang:   lang,
	}
}

// SetSeed sets the seed for the generator.
func (g *Generator) SetSeed(seed int64) {
	g.rng.Seed(seed)
}

func (g *Generator) Ra(array []string) string {
	return array[g.rng.Intn(len(array))]
}

func (g *Generator) Rw(mp map[string]int) string {
	// TODO: Cache weighted arrays?
	return g.Ra(weightedToArray(mp))
}

// RandGenMethod returns a random empire generation Method.
func (g *Generator) RandGenMethod() string {
	return g.Rw(GenEmpireMethods)
}

// GeneEmpireName generates a name for an empire.
func (g *Generator) GenEmpireName(tokens []genstory.TokenReplacement) (*genstory.Generated, error) {
	return g.txtGen.GenerateAndGiveMeTheTemplate(tokens, &NameConfig)
}

// weightedToArray converts a map of weighted values to an array.
func weightedToArray(weighted map[string]int) []string {
	var res []string
	for key, weight := range weighted {
		for j := 0; j < weight; j++ {
			res = append(res, key)
		}
	}
	return res
}

// GenEmpireMethods contains a map of empire name generation
// methods and their relative chance to be selected.
var GenEmpireMethods = map[string]int{
	MethodFormOfPlace:                      3,
	MethodTheAdjFormOfPlace:                2,
	MethodPlaceToAdjForm:                   2,
	MethodPlaceToAdjSubFormToPlural:        2,
	MethodFormOfPlaceToAdjSubFormToPlural:  1,
	MethodFormOfThePlaceAdjSubFormToPlural: 1,
	MethodFormOfTheAdjPlaceSubFormToPlural: 1,
	MethodTheAdjFormOfFoundingFigure:       3,
	MethodFormOfTheAdjFoundingFigure:       2,
	MethodFoundingFigureForm:               2,
	MethodFoundingFigureAdjForm:            2,
}

const (
	// Place based (like capital city, region, etc):
	MethodFormOfPlace                      = "[FORM] of [PLACE]"
	MethodTheAdjFormOfPlace                = "The [ADJECTIVE] [FORM] of [PLACE]"
	MethodPlaceToAdjForm                   = "[PLACE:adjective] [FORM]"
	MethodPlaceToAdjSubFormToPlural        = "[PLACE:adjective] [SUB_FORM:nounplural]"
	MethodFormOfPlaceToAdjSubFormToPlural  = "[FORM] of [PLACE:adjective] [SUB_FORM:nounplural]"
	MethodFormOfThePlaceAdjSubFormToPlural = "[FORM] of the [PLACE:adjective] [ADJECTIVE] [SUB_FORM:nounplural]"
	MethodFormOfTheAdjPlaceSubFormToPlural = "[FORM] of the [ADJECTIVE] [PLACE:adjective] [SUB_FORM:nounplural]"
	// Founding figure based (use last/family/dynasty name):
	MethodTheAdjFormOfFoundingFigure = "The [ADJECTIVE] [FORM] of [FOUNDING_FIGURE]"
	MethodFormOfTheAdjFoundingFigure = "[FORM] of the [ADJECTIVE] [FOUNDING_FIGURE]"
	MethodFoundingFigureForm         = "[FOUNDING_FIGURE]'s [FORM]"
	MethodFoundingFigureAdjForm      = "[FOUNDING_FIGURE]'s [ADJECTIVE] [FORM]"
	// Geographical location based:
	// TODO: Implement
	// Ethnic or cultural identity based:
	// TODO: Implement
	// MethodFormOfFaith       = "[FORM] of [FAITH]"
	// MethodTheAdjFormOfFaith = "The [ADJECTIVE] [FORM] of [FAITH]"
	// MethodAdjFormOfDeity = "[ADJECTIVE] [FORM] of [DEITY]"
	// MethodFormOfDeity    = "[FORM] of [DEITY]"
	// MethodFormOfCulture     = "[FORM] of [CULTURE]"
)

const (
	TokenPlace          = "[PLACE]"
	TokenAdjective      = "[ADJECTIVE]"
	TokenForm           = "[FORM]"
	TokenFoundingFigure = "[FOUNDING_FIGURE]"
	TokenSubForm        = "[SUB_FORM]"
	// TokenFaith          = "[FAITH]"
	// TokenDeity = "[DEITY]"
	// TokenCulture        = "[CULTURE]"
	// TokenLanguage       = "[LANGUAGE]"
)

var NameConfig = genstory.TextConfig{
	TokenIsMandatory: map[string]bool{
		TokenPlace:          true,
		TokenFoundingFigure: true,
		// TokenDeity:          true,
	},
	TokenPools: map[string][]string{
		TokenAdjective: Adjectives,
		TokenForm:      EmpireForms,
		TokenSubForm:   EmpireSubForms,
	},
	Tokens: []string{
		TokenPlace,
		TokenAdjective,
		TokenForm,
		TokenFoundingFigure,
		TokenSubForm,
		// TokenFaith,
		// TokenDeity,
		// TokenCulture,
		// TokenLanguage,
	},
	Templates: []string{
		MethodFormOfPlace,
		MethodTheAdjFormOfPlace,
		MethodPlaceToAdjForm,
		MethodPlaceToAdjSubFormToPlural,
		MethodFormOfPlaceToAdjSubFormToPlural,
		MethodFormOfThePlaceAdjSubFormToPlural,
		MethodFormOfTheAdjPlaceSubFormToPlural,
		MethodTheAdjFormOfFoundingFigure,
		MethodFormOfTheAdjFoundingFigure,
		MethodFoundingFigureForm,
		MethodFoundingFigureAdjForm,
		// MethodAdjFormOfDeity,
		// MethodFormOfDeity,
	},
	Modifiers:      map[string]func(string) string{},
	UseAllProvided: false,
}

// Adjectives contains a list of possible Adjectives for empires.
var Adjectives = []string{
	"great",
	"grand",
	"mighty",
	"holy",
	"just",
	"divine",
	"righteous",
	"noble",
}

// EmpireForms contains a list of possible forms (names) for empires.
var EmpireForms = []string{
	"Empire",
	"Dominion",
	"Union",
	"Commonwealth",
	"Domain",
	"Lands",
	"Fields",
	"Dynasty",
}

// EmpireSubForms contains a list of possible sub-forms for empire member governments.
var EmpireSubForms = []string{
	"Republic",
	"State",
	"Kingdom",
	"Territory",
	"Colony",
	"Province",
	"Region",
}
