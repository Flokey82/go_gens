package main

import (
	"log"
	"time"

	"github.com/Flokey82/go_gens/genstory"
)

func main() {
	for _, tc := range []struct {
		desc   string
		tokens []genstory.TokenReplacement
		titles []string
	}{{
		desc: "one token, variant titles",
		tokens: []genstory.TokenReplacement{{
			Token:       genstory.TokenNoun,
			Replacement: "Battle Tactics",
		}},
		titles: genstory.BookVariantTitles,
	}, {
		desc: "one token, instruction titles",
		tokens: []genstory.TokenReplacement{{
			Token:       genstory.TokenName,
			Replacement: "Glorbnorb",
		}},
		titles: genstory.BookInstructionTitles,
	}, {
		desc: "two tokens, variant titles",
		tokens: []genstory.TokenReplacement{{
			Token:       genstory.TokenName,
			Replacement: "Bleepblorp",
		}, {
			Token:       genstory.TokenNoun,
			Replacement: "Trickery",
		}},
		titles: genstory.BookVariantTitles,
	}, {
		desc: "three tokens, variant titles",
		tokens: []genstory.TokenReplacement{{
			Token:       genstory.TokenName,
			Replacement: "Herdle the Great",
		}, {
			Token:       genstory.TokenNoun,
			Replacement: "Goats",
		}, {
			Token:       genstory.TokenAdj,
			Replacement: "Boring",
		}},
		titles: genstory.BookVariantTitles,
	}} {
		log.Println(tc.desc)
		cfg := genstory.NewSimpleTitleConfig(tc.titles)
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
		if title, err := genstory.WeaponsTextConfig.Generate([]genstory.TokenReplacement{{
			Token:       genstory.WeaponTokenName,
			Replacement: "Soul Cleaver",
		}, {
			Token:       genstory.WeaponTokenType,
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
