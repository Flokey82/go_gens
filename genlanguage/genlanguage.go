package genlanguage

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"github.com/Flokey82/naming-language-gen/naming"
	"github.com/s0rg/fantasyname"
)

func GenLanguage(seed int64) *Language {
	// TODO: Based on the culture, we should seed the city names with
	// words that are preferred by the culture.
	// A culture might prefer to name cities to be named after their
	// gods, or their heroes, or their ancestors, or their enemies.
	// Or they might prefer to name cities after their natural
	// resources, or their natural features, or their regional natural
	// disasters... maybe even by the local fauna and flora.

	// We should define a number of different categories of words:
	//
	// - gods
	// - fauna
	// - flora
	// - weather / climate
	// - natural features
	// - natural disasters
	// - titles
	//
	// Then we store all base words that we generate in a dictionary
	// and we can use that to generate names for cities, people, etc.

	rand.Seed(seed)
	lang := naming.RandomLanguage(true, true)

	wordConfig := &naming.WordParams{
		MinSyllables: 1,
		MaxSyllables: 3,
		Structure:    naming.DefaultSyllableStructures,
	}

	nameConfig := &naming.NameParams{
		MinLength: naming.RandomRange(3, 5),
		MaxLength: naming.RandomRange(6, 20),
		WordParams: &naming.WordParams{
			MinSyllables: 2,
			MaxSyllables: naming.RandomRange(2, 7),
			Structure:    naming.DefaultSyllableStructures,
		},
		Joiners: "  -",
		Group:   "words",
	}

	cityConfig := &naming.NameParams{
		MinLength: naming.RandomRange(3, 5),
		MaxLength: naming.RandomRange(6, 15),
		WordParams: &naming.WordParams{
			MinSyllables: 2,
			MaxSyllables: naming.RandomRange(3, 7),
			Structure:    naming.DefaultSyllableStructures,
		},
		Joiners: "  -",
		Group:   "city",
	}

	firstNameConfig := &naming.NameParams{
		MinLength: naming.RandomRange(3, 5),
		MaxLength: naming.RandomRange(6, 15),
		WordParams: &naming.WordParams{
			MinSyllables: 2,
			MaxSyllables: naming.RandomRange(3, 7),
			Structure:    naming.DefaultSyllableStructures,
		},
		Joiners: "  -",
		Group:   "firstname",
	}

	lastNameConfig := &naming.NameParams{
		MinLength: naming.RandomRange(3, 5),
		MaxLength: naming.RandomRange(6, 15),
		WordParams: &naming.WordParams{
			MinSyllables: 2,
			MaxSyllables: naming.RandomRange(3, 7),
			Structure:    naming.DefaultSyllableStructures,
		},
		Joiners: "  -",
		Group:   "lastname",
	}

	return &Language{
		lang:            lang,
		wordConfig:      wordConfig,
		nameConfig:      nameConfig,
		cityConfig:      cityConfig,
		firstNameConfig: firstNameConfig,
		lastNameConfig:  lastNameConfig,
	}
}

type Language struct {
	lang            *naming.Language
	wordConfig      *naming.WordParams
	nameConfig      *naming.NameParams
	cityConfig      *naming.NameParams
	firstNameConfig *naming.NameParams
	lastNameConfig  *naming.NameParams
}

// GetWord returns a random, probably new, word from the language.
func (l *Language) GetWord(group string) string {
	return l.lang.GetWord(l.wordConfig, group)
}

// MakeName returns a random, probably new, name from the language.
func (l *Language) MakeName() string {
	return l.lang.MakeName(l.nameConfig)
}

// MakeCityName returns a random, probably new, city name from the language.
func (l *Language) MakeCityName() string {
	return l.lang.MakeName(l.cityConfig)
}

// GetFirstNamePoolSize returns the number of first names in the language.
func (l *Language) GetFirstNamePoolSize() int {
	return len(l.lang.Words.General[l.firstNameConfig.Group])
}

// MakeFirstName returns a random, probably new, first name from the language.
func (l *Language) MakeFirstName() string {
	return l.lang.MakeName(l.firstNameConfig)
}

// GetFirstName returns a random, pre-existing first name from the language.
func (l *Language) GetFirstName() string {
	return strings.Title(l.lang.GetWord(l.firstNameConfig.WordParams, l.firstNameConfig.Group))
}

// GetLastNamePoolSize returns the number of last names in the language.
func (l *Language) GetLastNamePoolSize() int {
	return len(l.lang.Words.General[l.lastNameConfig.Group])
}

// MakeLastName returns a random, probably new, last name from the language.
func (l *Language) MakeLastName() string {
	return l.lang.MakeName(l.lastNameConfig)
}

// GetLastName returns a random, pre-existing last name from the language.
func (l *Language) GetLastName() string {
	return strings.Title(l.lang.GetWord(l.lastNameConfig.WordParams, l.lastNameConfig.Group))
}

func NewFantasyName() (fmt.Stringer, error) {
	return fantasyname.Compile("")
}

type adjectivizationRule struct {
	name        string
	probability float64
	condition   *regexp.Regexp
	action      func(string) string
}

// stringSlice is the equivalent of JavaScript's string.slice(start, end).
func stringSlice(s string, start, end int) string {
	if end < 0 {
		end = len(s) + end
	}
	return s[start:end]
}

// These rules were borrowed from Azgaar's Fantasy Map Generator.
// See: https://github.com/Azgaar/Fantasy-Map-Generator/blob/master/utils/languageUtils.js
var adjectivizationRules = []adjectivizationRule{{
	name:        "guo",
	probability: 1,
	condition:   regexp.MustCompile(" Guo$"),
	action:      func(noun string) string { return stringSlice(noun, 0, -4) },
}, {
	name:        "orszag",
	probability: 1,
	condition:   regexp.MustCompile("orszag$"),
	action: func(noun string) string {
		if len(noun) < 9 {
			return noun + "ian"
		}
		return stringSlice(noun, 0, -6)
	},
}, {
	name:        "stan",
	probability: 1,
	condition:   regexp.MustCompile("stan$"),
	action: func(noun string) string {
		if len(noun) < 9 {
			return noun + "i"
		}
		return TrimVowels(stringSlice(noun, 0, -4), 3)
	},
}, {
	name:        "land",
	probability: 1,
	condition:   regexp.MustCompile("land$"),
	action: func(noun string) string {
		if len(noun) > 9 {
			return stringSlice(noun, 0, -4)
		}
		root := TrimVowels(stringSlice(noun, 0, -4), 0)
		if len(root) < 3 {
			return noun + "ic"
		}
		if len(root) < 4 {
			return root + "lish"
		}
		return root + "ish"
	},
}, {
	name:        "que",
	probability: 1,
	condition:   regexp.MustCompile("que$"),
	action: func(noun string) string {
		re := regexp.MustCompile("/que$/")
		return re.ReplaceAllString(noun, "can")
	},
}, {
	name:        "a",
	probability: 1,
	condition:   regexp.MustCompile("a$"),
	action: func(noun string) string {
		return noun + "n"
	},
}, {
	name:        "o",
	probability: 1,
	condition:   regexp.MustCompile("o$"),
	action: func(noun string) string {
		re := regexp.MustCompile("/o$/")
		return re.ReplaceAllString(noun, "an")
	},
}, {
	name:        "u",
	probability: 1,
	condition:   regexp.MustCompile("u$"),
	action: func(noun string) string {
		return noun + "an"
	},
}, {
	name:        "i",
	probability: 1,
	condition:   regexp.MustCompile("i$"),
	action: func(noun string) string {
		return noun + "an"
	},
}, {
	name:        "e",
	probability: 1,
	condition:   regexp.MustCompile("e$"),
	action: func(noun string) string {
		return noun + "an"
	},
}, {
	name:        "ay",
	probability: 1,
	condition:   regexp.MustCompile("ay$"),
	action: func(noun string) string {
		return noun + "an"
	},
}, {
	name:        "os",
	probability: 1,
	condition:   regexp.MustCompile("os$"),
	action: func(noun string) string {
		root := TrimVowels(stringSlice(noun, 0, -2), 0)
		if len(root) < 4 {
			return stringSlice(noun, 0, -1)
		}
		return root + "ian"
	},
}, {
	name:        "es",
	probability: 1,
	condition:   regexp.MustCompile("es$"),
	action: func(noun string) string {
		root := TrimVowels(stringSlice(noun, 0, -2), 0)
		if len(root) > 7 {
			return stringSlice(noun, 0, -1)
		}
		return root + "ian"
	},
}, {
	name:        "l",
	probability: 0.8,
	condition:   regexp.MustCompile("l$"),
	action: func(noun string) string {
		return noun + "ese"
	},
}, {
	name:        "n",
	probability: 0.8,
	condition:   regexp.MustCompile("n$"),
	action: func(noun string) string {
		return noun + "ese"
	},
}, {
	name:        "ad",
	probability: 0.8,
	condition:   regexp.MustCompile("ad$"),
	action: func(noun string) string {
		return noun + "ian"
	},
}, {
	name:        "an",
	probability: 0.8,
	condition:   regexp.MustCompile("an$"),
	action: func(noun string) string {
		return noun + "ian"
	},
}, {
	name:        "ish",
	probability: 0.25,
	condition:   regexp.MustCompile("^[a-zA-Z]{6}$"),
	action: func(noun string) string {
		return TrimVowels(stringSlice(noun, 0, -1), 3) + "ish"
	},
}, {
	name:        "an",
	probability: 0.5,
	condition:   regexp.MustCompile("^[a-zA-Z]{0-7}$"),
	action: func(noun string) string {
		return TrimVowels(noun, 3) + "an"
	},
}}

// GetAdjective get adjective form from noun
func GetAdjective(noun string) string {
	for _, rule := range adjectivizationRules {
		if P(rule.probability) && rule.condition.MatchString(noun) {
			return rule.action(noun)
		}
	}
	return noun // no rule applied, return noun as is
}

// chars that serve as vowels
const vowelRange = `aeiouyɑ'əøɛœæɶɒɨɪɔɐʊɤɯаоиеёэыуюяàèìòùỳẁȁȅȉȍȕáéíóúýẃőűâêîôûŷŵäëïöüÿẅãẽĩõũỹąęįǫųāēīōūȳăĕĭŏŭǎěǐǒǔȧėȯẏẇạẹịọụỵẉḛḭṵṳ`

// IsVowel returns true if the given rune is a vowel.
func IsVowel(c rune) bool {
	return strings.IndexRune(vowelRange, c) != -1
}

// TrimVowels remove vowels from the end of the string.
func TrimVowels(str string, minLength int) string {
	for len(str) > minLength && IsVowel(rune(str[len(str)-1])) {
		str = str[:len(str)-1]
	}
	return str
}

// GetNounPlural returns the plural form of a noun.
// This takes in account "witch" and "fish" which are
// irregular.
func GetNounPlural(noun string) string {
	if strings.HasSuffix(noun, "s") || strings.HasSuffix(noun, "x") || strings.HasSuffix(noun, "z") || strings.HasSuffix(noun, "ch") || strings.HasSuffix(noun, "sh") {
		return noun + "es"
	}
	return noun + "s"
}

// probability shorthand
func P(probability float64) bool {
	if probability >= 1.0 {
		return true
	}
	if probability <= 0 {
		return false
	}
	return rand.Float64() < probability
}
