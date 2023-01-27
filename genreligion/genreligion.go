package genreligion

import (
	"math/rand"
	"strings"

	"github.com/Flokey82/go_gens/genlanguage"
)

// Generator is a generator for religions.
type Generator struct {
	rng *rand.Rand
}

// NewGenerator creates a new religion generator.
func NewGenerator(seed int64) *Generator {
	return &Generator{
		rng: rand.New(rand.NewSource(seed)),
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

// GenNamedTypeOfForm generates a name for a named religion type based on a given
// religion form ("Polytheism", "Dualism", etc).
// E.g. "Pradanium deities".
func (g *Generator) GenNamedTypeOfForm(name, form string) string {
	return name + " " + g.Rw(Types[form])
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

// FaitOfSupremePrefixes contains a list of prefixes identifying the group of
// followers of a supreme deity or leader.
var FaitOfSupremePrefixes = []string{
	"Faith",
	"Way",
	"Path",
	"Word",
	"Truth",
	"Law",
	"Order",
	"Light",
	"Darkness",
	"Gift",
	"Grace",
	"Witnesses",
	"Servants",
	"Messengers",
	"Believers",
	"Disciples",
	"Followers",
	"Children",
	"Brothers",
	"Sisters",
	"Brothers and Sisters",
	"Sons",
	"Daughters",
	"Sons and Daughters",
	"Brides",
	"Grooms",
	"Brides and Grooms",
}

const (
	// Religion groups.
	GroupFolk      = "Folk"
	GroupOrganized = "Organized"
	GroupCult      = "Cult"
	GroupHeresy    = "Heresy"

	// Religion forms.
	FormShamanism       = "Shamanism"
	FormAnimism         = "Animism"
	FormAncestorWorship = "Ancestor worship"
	FormPolytheism      = "Polytheism"
	FormDualism         = "Dualism"
	FormMonotheism      = "Monotheism"
	FormNontheism       = "Non-theism"
	FormCult            = "Cult"
	FormDarkCult        = "Dark Cult"
	FormHeresy          = "Heresy"
	// FormNature = "Nature"
)

var Forms = map[string]map[string]int{
	GroupFolk: {
		FormShamanism:       2,
		FormAnimism:         2,
		FormAncestorWorship: 1,
		FormPolytheism:      2,
	},
	GroupOrganized: {
		FormPolytheism: 5,
		FormDualism:    1,
		FormMonotheism: 4,
		FormNontheism:  1,
	},
	FormCult: {
		FormCult:     1,
		FormDarkCult: 1,
	},
	FormHeresy: {
		FormHeresy: 1,
	},
}

var Types = map[string]map[string]int{
	FormShamanism: {
		"Beliefs":   3,
		"Shamanism": 2,
		"Spirits":   1,
	},
	FormAnimism: {
		"Spirits": 1,
		"Beliefs": 1,
	},
	FormAncestorWorship: {
		"Beliefs":     1,
		"Forefathers": 2,
		"Ancestors":   2,
	},
	FormPolytheism: {
		"Deities":  3,
		"Faith":    1,
		"Gods":     1,
		"Pantheon": 1,
	},
	FormDualism: {
		"Religion": 3,
		"Faith":    1,
		"Cult":     1,
	},
	FormMonotheism: {
		"Religion": 1,
		"Church":   1,
	},
	FormNontheism: {
		"Beliefs": 3,
		"Spirits": 1,
	},
	FormCult: {
		"Cult":    4,
		"Sect":    4,
		"Arcanum": 1,
		"Coterie": 1,
		"Order":   1,
		"Worship": 1,
	},
	FormDarkCult: {
		"Cult":      2,
		"Sect":      2,
		"Blasphemy": 1,
		"Circle":    1,
		"Coven":     1,
		"Idols":     1,
		"Occultism": 1,
	},
	FormHeresy: {
		"Heresy":      3,
		"Sect":        2,
		"Apostates":   1,
		"Brotherhood": 1,
		"Circle":      1,
		"Dissent":     1,
		"Dissenters":  1,
		"Iconoclasm":  1,
		"Schism":      1,
		"Society":     1,
	},
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
	MethodRandomType     = "Random + type"
	MethodRandomIsm      = "Random + ism"
	MethodSurpremeIsm    = "Supreme + ism"
	MethodFaithOfSupreme = "Faith of + Supreme"
	MethodPlaceIsm       = "Place + ism"
	MethodCultureIsm     = "Culture + ism"
	MethodPlaceIanType   = "Place + ian + type"
	MethodCultureType    = "Culture + type"
)

// genReligionMethods contains a map of religion name generation
// methods and their relative chance to be selected.
var GenReligionMethods = map[string]int{
	MethodRandomType:     3,
	MethodRandomIsm:      1,
	MethodSurpremeIsm:    5,
	MethodFaithOfSupreme: 5,
	MethodPlaceIsm:       1,
	MethodCultureIsm:     2,
	MethodPlaceIanType:   6,
	MethodCultureType:    4,
}
