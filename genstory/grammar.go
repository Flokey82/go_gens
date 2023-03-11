package genstory

import (
	"fmt"
	"strings"

	"github.com/Flokey82/go_gens/genlanguage"
)

// Rules defines the rules for a story.
type Rules struct {
	Expansions map[string][]string // Pool of expansions for each token
	Start      string              // Start rule
}

// NewStory creates a new story from the given rules.
func (r *Rules) NewStory() *Story {
	return &Story{
		Rules:    r,
		Assigned: make(map[string]string),
	}
}

var ExampleRules = &Rules{
	Expansions: map[string][]string{
		"animal": genlanguage.GenBase[genlanguage.GenBaseAnimal],
		"name": {
			"John",
			"Jane",
			"Bob",
			"Mary",
			"Peter",
			"Paul",
			"George",
			"Ringo",
		},
		"victim": {
			"old lady",
			"blind person",
			"child",
		},
		"loot": {
			"money",
			"jewels",
			"gold",
			"silver",
			"coins",
		},
		"crime": {
			"arson",
			"extortion",
			"racketeering",
			"theft from [victim:a]",
			"theft of [loot]",
			"crimes against humanity",
			"defrauding [victim:a]",
			"running a crypto scam",
		},
		"tragedy": {
			", but [petname] was eaten by [name] the [animal]",
			", but [petname] the [pet] died",
			", but [petname] the [pet] ran away",
			", but [petname] the [pet] was stolen",
			", but [petname] the [pet] was arrested for [crime]",
		},
	},
	Start: "[hero/name:quote] bought [pet/animal:a]. [hero] loved the [pet] and named it [petname/name][tragedy].",
}

// Story is a new story instance that can be expanded from a set of rules.
// After expansion, all assigned token expansions will be stored in Assigned if
// they are needed for later use. (For example, the names picked for the hero,
// name of the pet, and what animal he rode into town).
type Story struct {
	*Rules                     // Rules for the story
	Assigned map[string]string // Assigned expansions for tokens that have assigned expansions
}

// Expand expands the story from the start rule.
func (g *Story) Expand() (string, error) {
	return g.ExpandRule(g.Start)
}

// ExpandRule expands from the start rule.
func (g *Story) ExpandRule(rule string) (string, error) {
	// Find all tokens in the rule.
	tokens, err := ExtractGrammarTokens(rule)
	if err != nil {
		return "", err
	}

	// Expand each token, and replace the token with the expansion.
	for _, token := range tokens {
		// Expand the token, which will give us (possibly) a new rule.
		expansion, err := g.ExpandToken(token)
		if err != nil {
			return "", err
		}

		// Expand the rule.
		expansion, err = g.ExpandRule(expansion)
		if err != nil {
			return "", err
		}

		// Apply modifiers.
		expansion = ApplyModifiers(expansion, token.Modifiers)

		// Replace the token with the expansion.
		rule = strings.Replace(rule, token.FullToken, expansion, 1)
	}

	return rule, nil
}

// ExpandToken expands a single token.
func (g *Story) ExpandToken(token GrammarToken) (string, error) {
	// Check if the token has an assigned expansion.
	if expansion, ok := g.Assigned[token.Token]; ok {
		return expansion, nil
	}

	// Find the pool for the token.
	pool, ok := g.Expansions[token.Pool]
	if !ok {
		return "", fmt.Errorf("no pool found for token %s and reference %s", token.Token, token.Pool)
	}

	// Choose a random expansion from the pool.
	expansion := randArrayString(pool)

	// Store the expansion for the token.
	g.Assigned[token.Token] = expansion

	return expansion, nil
}

// GrammarToken represents a token in a grammar rule.
// TODO:
// - Add flags to prevent re-use of expansions (10 people named John).
// - Add flags for n-repetitions.
// - Allow association of token expansions with tokens. ([green->house])
// - Add conditionals.
type GrammarToken struct {
	Token     string   // The token, including brackets.
	Pool      string   // The pool from which to pick a random value.
	FullToken string   // The full token, including modifiers.
	Modifiers []string // The modifiers for the token.
	Start     int
	End       int
}

// ExtractGrammarTokens extracts all tokens from a string.
// A token is a string surrounded by square brackets ('[]'),
// has a pool reference (e.g. "hero/name") and can have modifiers,
// separated by a colon (':').
// The pool reference (denoted by a slash '/') is used to find the
// expansion pool for the token.
func ExtractGrammarTokens(s string) ([]GrammarToken, error) {
	var tokens []GrammarToken
	for _, t := range findAllTokens(s) {
		token := ExtractGrammarToken(t.token)
		token.Start = t.start
		token.End = t.end
		tokens = append(tokens, token)
	}
	return tokens, nil
}

// ExtractGrammarToken extracts a single token from a string.
func ExtractGrammarToken(s string) GrammarToken {
	// Trim off the brackets.
	token := s[1 : len(s)-1]

	// Split off all modifiers.
	modifiers := strings.Split(token, ":")
	token = modifiers[0]
	modifiers = modifiers[1:]

	// Split off the pool reference.
	var poolReference string
	if strings.Contains(token, "/") {
		parts := strings.Split(token, "/")
		poolReference = parts[1]
		token = parts[0]
	} else {
		// If there is no pool reference, the token is the pool reference.
		poolReference = token
	}

	return GrammarToken{
		Token:     "[" + token + "]", // Add the brackets back.
		FullToken: s,
		Pool:      poolReference,
		Modifiers: modifiers,
	}
}
