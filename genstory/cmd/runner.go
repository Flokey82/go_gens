package main

import (
	"log"
	"time"

	"github.com/Flokey82/go_gens/genstory"
	"github.com/Flokey82/go_gens/genstory/genbooks"
	"github.com/Flokey82/go_gens/genstory/genweapons"
)

func main() {
	for _, tc := range []struct {
		desc   string
		tokens []genstory.TokenReplacement
		titles []string
	}{{
		desc: "one token, variant titles",
		tokens: []genstory.TokenReplacement{{
			Token:       genbooks.TokenNoun,
			Replacement: "Battle Tactics",
		}},
		titles: genbooks.BookVariantTitles,
	}, {
		desc: "one token, instruction titles",
		tokens: []genstory.TokenReplacement{{
			Token:       genbooks.TokenName,
			Replacement: "Glorbnorb",
		}},
		titles: genbooks.BookInstructionTitles,
	}, {
		desc: "two tokens, variant titles",
		tokens: []genstory.TokenReplacement{{
			Token:       genbooks.TokenName,
			Replacement: "Bleepblorp",
		}, {
			Token:       genbooks.TokenNoun,
			Replacement: "Trickery",
		}},
		titles: genbooks.BookVariantTitles,
	}, {
		desc: "three tokens, variant titles",
		tokens: []genstory.TokenReplacement{{
			Token:       genbooks.TokenName,
			Replacement: "Herdle the Great",
		}, {
			Token:       genbooks.TokenNoun,
			Replacement: "Goats",
		}, {
			Token:       genbooks.TokenAdj,
			Replacement: "Boring",
		}},
		titles: genbooks.BookVariantTitles,
	}} {
		log.Println(tc.desc)
		cfg := genbooks.NewSimpleTitleConfig(tc.titles)
		for i := 0; i < 10; i++ {
			if title, err := cfg.Generate(tc.tokens); err != nil {
				log.Println(i, err)
			} else {
				log.Println(i, title)
			}
		}
	}

	log.Println("weapon flavor text")
	for i := 0; i < 10; i++ {
		if title, err := genweapons.WeaponsTextConfig.Generate([]genstory.TokenReplacement{{
			Token:       genweapons.WeaponTokenName,
			Replacement: "Soul Cleaver",
		}, {
			Token:       genweapons.WeaponTokenType,
			Replacement: "longsword",
		}}); err != nil {
			log.Println(i, err)
		} else {
			log.Println(i, title)
		}
	}

	log.Println("grammar example")
	for i := 0; i < 10; i++ {
		g := genstory.ExampleRules.NewStory(time.Now().UnixNano())
		if title, err := g.Expand(); err != nil {
			log.Println(i, err)
		} else {
			log.Println(i, title)
		}
	}
}
