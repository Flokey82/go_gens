package genstory

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

// TokenReplacement is a replacement for a token in a text.
type TokenReplacement struct {
	Token       string // The token to replace.
	Replacement string // The replacement text.
}

// TextConfig is a configuration for generating text.
type TextConfig struct {
	TokenPools       map[string][]string // A map of token names to a list of possible values.
	TokenIsMandatory map[string]bool     // A map of token names to a boolean indicating whether the token is mandatory.
	Tokens           []string            // A list of tokens that are required to be replaced.
	Templates        []string            // A list of possible templates.
	Title            bool                // Capitalize the first letter of each word in the text.
}

// Generate generates a text from the provided tokens and the configuration.
func (c *TextConfig) Generate(provided []TokenReplacement) (string, error) {
	return generateText(provided, c.Templates, c.Tokens, c.TokenIsMandatory, c.TokenPools, c.Title)
}

// GenerateWithTemplate generates a text from the provided tokens and the provided template.
func (c *TextConfig) GenerateWithTemplate(provided []TokenReplacement, template string) (string, error) {
	return generateText(provided, []string{template}, c.Tokens, c.TokenIsMandatory, c.TokenPools, c.Title)
}

// GenerateTitle generates a text from the provided tokens and a list of
// possible templates.
//   - The provided tokens are used to replace the tokens in the possible templates.
//   - If a token is not provided and optional, it is replaced with a random value.
//   - If a token is not provided and not optional, all templates that require that
//     token are excluded.
//
// TODO: Also return the selected template, and the individual replacements,
// so that the caller can use them for the description or the generation of content.
func GenerateTitle(provided []TokenReplacement, titles []string) (string, error) {
	return generateText(provided, titles, DefaultTitleTokens, DefaultTitleTokenIsMandatory, DefaultTitleTokenPool, true)
}

func generateText(provided []TokenReplacement, templates, tokens []string, isMandatory map[string]bool, tokenRandom map[string][]string, capitalize bool) (string, error) {
	// Count how many token replacements we have for each token.
	tokenReplacements := map[string]int{}
	for _, replacement := range provided {
		tokenReplacements[replacement.Token]++
	}

	// Loop over all templates and find the ones where we have all required tokens.
	possibleTemplates := []string{}
	for _, i := range rand.Perm(len(templates)) {
		template := templates[i]
		// Check if we have all required tokens the required number of times.
		var missingToken bool
		for _, token := range tokens {
			if tokenReplacements[token] < strings.Count(template, token) {
				if isMandatory[token] {
					missingToken = true
					break
				}
			}
		}

		// Something is missing, skip this template.
		if missingToken {
			continue
		}

		// Also make sure all tokens we have provided are available in the template,
		// since we want to pick a complete template, referencing all provided tokens.
		for _, replacement := range provided {
			if strings.Count(template, replacement.Token) < tokenReplacements[replacement.Token] {
				missingToken = true
				break
			}
		}

		// Something is missing, skip this template.
		if missingToken {
			continue
		}

		// We have all required tokens, add the template to the list of possible templates.
		possibleTemplates = append(possibleTemplates, template)
	}

	// If we have no possible templates, return an error.
	if len(possibleTemplates) == 0 {
		return "", errors.New("no possible templates satisfying the provided tokens")
	}

	// Pick a random text.
	text := randArrayString(possibleTemplates)

	// Replace all tokens with the provided replacements.
	for _, replacement := range provided {
		text = strings.Replace(text, replacement.Token, replacement.Replacement, 1)
	}

	// Replace all optional tokens with random replacements.
	remainingTokens, err := ExtractTokens(text)
	if err != nil {
		return "", fmt.Errorf("failed to extract tokens from template: %v", err)
	}

	// Relplace each token one by one until we can't find any more.
	for _, token := range remainingTokens {
		if !isMandatory[token] && strings.Contains(text, token) {
			// Pick a random replacement.
			// TODO: What to do if we don't have any replacements for a token?
			replacement := randArrayString(tokenRandom[token])

			// Replace the token.
			text = strings.Replace(text, token, replacement, 1)
		}
	}
	if capitalize {
		// Capitalize the first letter of each word in the text.
		text = strings.Title(text)
	} else {
		// Capitalize the first letter of the text.
		text = strings.ToUpper(text[:1]) + text[1:]
	}
	return text, nil
}

func randArrayString(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	return arr[rand.Intn(len(arr))]
}
