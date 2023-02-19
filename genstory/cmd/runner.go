package main

import (
	"log"
	"time"

	"github.com/Flokey82/go_gens/genstory"
)

func main() {
	log.Println(genstory.NewWorld(time.Now().UnixNano()))

	tokens := []genstory.TokenReplacement{{
		Token:       genstory.TokenNoun,
		Replacement: "Battle Tactics",
	}}
	for i := 0; i < 10; i++ {
		log.Println(genstory.GenerateTitle(tokens, genstory.BookVariantTitles))
	}
	tokens = []genstory.TokenReplacement{{
		Token:       genstory.TokenName,
		Replacement: "Glorbnorb",
	}}
	for i := 0; i < 10; i++ {
		log.Println(genstory.GenerateTitle(tokens, genstory.BookInstructionTitles))
	}

	tokens = []genstory.TokenReplacement{{
		Token:       genstory.TokenName,
		Replacement: "Bleepblorp",
	}, {
		Token:       genstory.TokenNoun,
		Replacement: "Trickery",
	}}

	for i := 0; i < 10; i++ {
		log.Println(genstory.GenerateTitle(tokens, genstory.BookVariantTitles))
	}
}
