package genreligion

import (
	"github.com/Flokey82/go_gens/genlanguage"
	"github.com/Flokey82/go_gens/genstory"
)

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

const (
	MethodRandomType     = "[RANDOM] [TYPE]"
	MethodRandomIsm      = "[RANDOM:trimvowels]ism"
	MethodSurpremeIsm    = "[SUPREME:trimvowels]ism"
	MethodFaithOfSupreme = "[FAITH] of [SUPREME]"
	MethodPlaceIsm       = "[PLACE:trimvowels]ism"
	MethodCultureIsm     = "[CULTURE:trimvowels]ism"
	MethodPlaceIanType   = "[PLACE]ian [TYPE]"
	MethodCultureType    = "[CULTURE] [TYPE]"
)

const (
	TokenPlace    = "[PLACE]"
	TokenCulture  = "[CULTURE]"
	TokenRandom   = "[RANDOM]"
	TokenSurpreme = "[SUPREME]"
	TokenType     = "[TYPE]"
	TokenFaith    = "[FAITH]"
)

var NameConfig = genstory.TextConfig{
	TokenIsMandatory: map[string]bool{
		TokenSurpreme: true,
		TokenCulture:  true,
		TokenRandom:   true,
		TokenPlace:    true,
		TokenType:     true,
	},
	TokenPools: map[string][]string{
		TokenFaith: FaitOfSupremePrefixes,
	},
	Tokens: []string{
		TokenPlace,
		TokenCulture,
		TokenRandom,
		TokenSurpreme,
		TokenType,
		TokenFaith,
	},
	Templates: []string{
		MethodRandomType,
		MethodRandomIsm,
		MethodSurpremeIsm,
		MethodFaithOfSupreme,
		MethodPlaceIsm,
		MethodCultureIsm,
		MethodPlaceIanType,
		MethodCultureType,
	},
	Modifiers: map[string]func(string) string{
		"trimvowels": func(s string) string {
			return genlanguage.TrimVowels(s, 3)
		},
	},
	UseAllProvided: false,
}
