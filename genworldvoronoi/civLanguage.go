package genworldvoronoi

import (
	"fmt"
	"math/rand"

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

func (l *Language) GetWord(group string) string {
	return l.lang.GetWord(l.wordConfig, group)
}

func (l *Language) MakeName() string {
	return l.lang.MakeName(l.nameConfig)
}

func (l *Language) MakeCityName() string {
	return l.lang.MakeName(l.cityConfig)
}

func (l *Language) MakeFirstName() string {
	return l.lang.MakeName(l.firstNameConfig)
}

func (l *Language) MakeLastName() string {
	return l.lang.MakeName(l.lastNameConfig)
}

func newFantasyName() (fmt.Stringer, error) {
	return fantasyname.Compile("")
}
