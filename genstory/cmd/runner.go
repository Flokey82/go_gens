package main

import (
	"log"
	"time"

	"github.com/Flokey82/go_gens/genstory"
)

func main() {
	log.Println(genstory.NewWorld(time.Now().UnixNano()))

	for _, tc := range []struct {
		tokens []genstory.TokenReplacement
		titles []string
	}{{
		tokens: []genstory.TokenReplacement{{
			Token:       genstory.TokenNoun,
			Replacement: "Battle Tactics",
		}},
		titles: genstory.BookVariantTitles,
	}, {
		tokens: []genstory.TokenReplacement{{
			Token:       genstory.TokenName,
			Replacement: "Glorbnorb",
		}},
		titles: genstory.BookInstructionTitles,
	}, {
		tokens: []genstory.TokenReplacement{{
			Token:       genstory.TokenName,
			Replacement: "Bleepblorp",
		}, {
			Token:       genstory.TokenNoun,
			Replacement: "Trickery",
		}},
		titles: genstory.BookVariantTitles,
	}} {
		for i := 0; i < 10; i++ {
			if title, err := genstory.GenerateTitle(tc.tokens, tc.titles); err != nil {
				log.Println(i, err)
			} else {
				log.Println(i, title)
			}
		}
	}
}
