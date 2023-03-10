package main

import (
	"log"

	"github.com/Flokey82/go_gens/genthesaurus"
)

func main() {
	// Load thesaurus from file
	ts, err := genthesaurus.NewFromJSON("data/ea-thesaurus-export.json")
	if err != nil {
		log.Fatal(err)
	}

	// Log the entries.
	ts.Log()
}
