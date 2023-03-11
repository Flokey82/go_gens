package genreligion

import (
	"math/rand"
	"strings"

	"github.com/Flokey82/go_gens/genlanguage"
	"github.com/Flokey82/go_gens/genstory"
)

// Generator is a generator for religions.
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

// RandGenMethod returns a random religion generation method.
func (g *Generator) RandGenMethod() string {
	return g.Rw(GenReligionMethods)
}

// RandFormFromGroup returns a random religion form based on a given religion
// group ("Folk", "Organized", etc).
func (g *Generator) RandFormFromGroup(group string) string {
	return g.Rw(Forms[group])
}

// RandTypeFromForm generates a random religion type based on a given religion
// form ("Polytheism", "Dualism", etc).
func (g *Generator) RandTypeFromForm(form string) string {
	return g.Rw(Types[form])
}

// GenNamedIsm generates a name for a religion based on a given religion form
// ("Polytheism", "Dualism", etc).
// E.g. "Pradaniumism".
func (g *Generator) GenNamedIsm(name string) string {
	return genlanguage.TrimVowels(name, 3) + "ism"
}

// GenFaithName generates a name for a faith.
func (g *Generator) GenFaithName(tokens []genstory.TokenReplacement) (name, method string, err error) {
	return g.txtGen.GenerateAndGiveMeTheTemplate(tokens, &NameConfig)
}

// GenNameFaitOfSurpreme generates a name for a faith of a supreme deity or leader.
// E.g. "Way of Grognark".
func (g *Generator) GenNameFaitOfSupreme(supreme string) string {
	// Select a random name from the list.
	// but ensure that the name is not a subset of the deity name
	// and vice versa. This is to avoid names like "The Way of The Way".
	var prefix string
	for i := 0; i < 100; i++ {
		prefix = g.Ra(FaitOfSupremePrefixes)
		if !strings.Contains(strings.ToLower(supreme), strings.ToLower(prefix)) &&
			!strings.Contains(strings.ToLower(prefix), strings.ToLower(supreme)) {
			break
		}
	}
	return prefix + " of " + supreme
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

const (
	// Expansion modes.
	ReligionExpGlobal  = "global"
	ReligionExpState   = "state"
	ReligionExpCulture = "culture"
)
