# genstory: Quick and dirty story generation

This is just a playground for flavor text generation. It's not very good, but it's a start.
**Please note that this is still a work in progress, but please feel free to make suggestions or contribute / improve.**
Currently, there are two ways to set up text generation:

1. Using the templates
2. Using the rules / grammar

## Templates

The template system uses simple text templates with placeholders for the various parts of the story which will be replaced using the token pools defined in the configs for the story generation.

### Example

```go
package main

import (
    "fmt"
    "github.com/Flokey82/go_gens/genstory"
)

func main() {
	// Define some tokens.
	tokenColor := "[color]"
	tokenItem := "[item]"
	tokenLocation := "[location]"
	tokenPerson := "[person]"

	// Set up a config for the generator.
	cfg := &genstory.TextConfig{
		TokenPools: map[string][]string{
			tokenColor:    {"red", "blue", "green", "yellow"},
			tokenItem:     {"skirt", "dress", "shirt", "sock"},
			tokenLocation: {"forest", "mountain", "cave", "lake"},
		},
		TokenIsMandatory: map[string]bool{
			tokenPerson: true,
		},
		Tokens: []string{tokenColor, tokenItem, tokenLocation, tokenPerson},
		Templates: []string{
			"The [color] [item] was found in the [location] by someone called [person:quote].",
			"The [item] was found by [person] in the [location].",
			"[person] lost my [item] in the [location].",
		},
		UseAllProvided: true,
	}
	for i := 0; i < 10; i++ {
		// Generate a story.
		story, err := cfg.Generate([]genstory.TokenReplacement{
			{Token: tokenPerson, Replacement: "John"},
		})
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(story.Text)
	}
}
```

## Rules / Grammar

The rule system is a bit more complex and allows for more flexibility like recursive expansions, assignment of individual tokens to labels (like the names of people, places, etc. for reuse), and more.

### Example

```go
package main

import (
    "fmt"
    "github.com/Flokey82/go_gens/genstory"
)

func main() {
	rules := &genstory.Rules{
		Expansions: map[string][]string{
			"animal": {
				"cat",
				"dog",
				"mouse",
				"elephant",
				"tiger",
				"lion",
			},
			"name": {
				"John",
				"Jane",
				"Bob",
				"Mary",
				"Peter",
				"Paul",
			},
			"event": {
				"snuggled by [animal:a] named [attacker/name:quote]",
				"bitten by [animal:a] named [attacker/name:quote]",
				"chased by [animal:a] named [attacker/name:quote]",
				"licked by [animal:a] named [attacker/name:quote]",
				"purred at by [animal:a] named [attacker/name:quote]",
			},
		},
		Start: "Yesterday, the quiet town of Everwood was shocked to discover that mayor [mayor/name] was [event]. [mayor/name] was visibly shaken by the incident." +
			"\"I'm not sure what to make of this,\" said [mayor/name]. \"I've never been [event] before.\"",
	}

	story := rules.NewStory(time.Now().UnixNano())
	txt, err := story.Expand()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(txt)
}
```
