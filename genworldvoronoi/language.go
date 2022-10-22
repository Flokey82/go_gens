package genworldvoronoi

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/Flokey82/naming-language-gen/naming"
)

func GenLanguage(seed int64) *Language {
	rand.Seed(seed)
	lang := naming.RandomLanguage(true, true)
	lang.Describe()

	l := &Language{
		lang: lang,
		wordConfig: &naming.WordParams{
			MinSyllables: 1,
			MaxSyllables: 3,
			Structure:    naming.DefaultSyllableStructures,
		},
		nameConfig: &naming.NameParams{
			MinLength: naming.RandomRange(3, 5),
			MaxLength: naming.RandomRange(6, 20),
			WordParams: &naming.WordParams{
				MinSyllables: 2,
				MaxSyllables: naming.RandomRange(2, 7),
				Structure:    naming.DefaultSyllableStructures,
			},
			Joiners: "  -",
			Group:   "words",
		},
	}

	// generate some words based on our new awesome language
	group := "words"
	for i := 0; i < 10; i++ {
		fmt.Println(l.GetWord(group))
	}

	wordList := strings.Join(lang.Words.General[group], ", ")
	fmt.Printf("[%v]: %v\n", group, wordList)

	// generate some names
	for i := 0; i < 10; i++ {
		l.MakeName()
	}

	nameList := strings.Join(lang.Words.Names, ", ")
	fmt.Printf("[%v]: %v\n", "names", nameList)
	return l
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
