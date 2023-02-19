package main

import (
	"log"
	"time"

	"github.com/Flokey82/go_gens/genstory"
)

func main() {
	log.Println(genstory.NewWorld(time.Now().UnixNano()))

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
}
