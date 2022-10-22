package genworldvoronoi

import (
	"math/rand"

	"github.com/Flokey82/naming-language-gen/naming"
)

func GenLanguage(seed int64) *Language {
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

	return &Language{
		lang:       lang,
		wordConfig: wordConfig,
		nameConfig: nameConfig,
	}
}

type Language struct {
	lang       *naming.Language
	wordConfig *naming.WordParams
	nameConfig *naming.NameParams
}

func (l *Language) GetWord(group string) string {
	return l.lang.GetWord(l.wordConfig, group)
}

func (l *Language) MakeName() string {
	return l.lang.MakeName(l.nameConfig)
}
