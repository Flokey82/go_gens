package genstory

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/Flokey82/go_gens/genlanguage"
)

// TokenReplacement is a replacement for a token in a text.
// TODO: Allow infinite re-use of tokens.
type TokenReplacement struct {
	Token       string // The token to replace.
	Replacement string // The replacement text.
}

// TextConfig is a configuration for generating text.
type TextConfig struct {
	TokenPools       map[string][]string            // A map of token names to a list of possible values.
	TokenIsMandatory map[string]bool                // A map of token names to a boolean indicating whether the token is mandatory.
	Tokens           []string                       // A list of tokens that are required to be replaced.
	Templates        []string                       // A list of possible templates.
	Title            bool                           // Capitalize the first letter of each word in the text.
	Modifiers        map[string]func(string) string // A map of token names to a function that modifies the replacement text.
}

// Generate generates a text from the provided tokens and the configuration.
func (c *TextConfig) Generate(provided []TokenReplacement) (string, error) {
	return generateText(provided, c.Templates, c.Tokens, c.TokenIsMandatory, c.TokenPools, c.Modifiers, c.Title)
}

// GenerateWithTemplate generates a text from the provided tokens and the provided template.
func (c *TextConfig) GenerateWithTemplate(provided []TokenReplacement, template string) (string, error) {
	return generateText(provided, []string{template}, c.Tokens, c.TokenIsMandatory, c.TokenPools, c.Modifiers, c.Title)
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
	return generateText(provided, titles, DefaultTitleTokens, DefaultTitleTokenIsMandatory, DefaultTitleTokenPool, nil, true)
}

func generateText(
	provided []TokenReplacement,
	templates, tokens []string,
	isMandatory map[string]bool,
	tokenRandom map[string][]string,
	modifierFuncs map[string]func(string) string,
	capitalize bool) (string, error) {

	// Function for applying modifiers to a string.
	applyModifiers := func(s string, modifiers []string) string {
		for _, modifier := range modifiers {
			// Check if we have a custom modifier function for this token.
			if modifierFuncs != nil {
				if fn, ok := modifierFuncs[modifier]; ok {
					s = fn(s)
				}
			}
			if fn, ok := DefaultModifiers[modifier]; ok {
				s = fn(s)
			}
		}
		return s
	}

	// Count how many token replacements we have for each token.
	tokenReplacements := make(map[string][]string)
	for _, replacement := range provided {
		tokenReplacements[replacement.Token] = append(tokenReplacements[replacement.Token], replacement.Replacement)
	}

	type candidate struct {
		template        string
		extractedTokens []ExtractedToken
	}

	// Loop over all templates and find the ones where we have all required tokens.
	var possibleTemplates []candidate
	for _, i := range rand.Perm(len(templates)) {
		template := templates[i]

		// TODO: Maybe cache the extracted tokens somewhere.
		extracted, err := ExtractTokens(template)
		if err != nil {
			return "", err
		}

		// Count how many times each token appears in the template.
		tokenCounts := make(map[string]int)
		for _, token := range extracted {
			tokenCounts[token.Token]++
		}

		// Check if we have all required tokens the required number of times.
		// TODO: Allow tokens to be re-usable.
		var missingToken bool
		for _, token := range tokens {
			if len(tokenReplacements[token]) < tokenCounts[token] {
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
			// TODO: Take in account that tokens might have modifiers, which defeats strings.Count.
			if tokenCounts[replacement.Token] < len(tokenReplacements[replacement.Token]) {
				missingToken = true
				break
			}
		}

		// Something is missing, skip this template.
		if missingToken {
			continue
		}

		// We have all required tokens, add the template to the list of possible templates.
		possibleTemplates = append(possibleTemplates, candidate{
			template:        template,
			extractedTokens: extracted,
		})
	}

	// If we have no possible templates, return an error.
	if len(possibleTemplates) == 0 {
		return "", errors.New("no possible templates satisfying the provided tokens")
	}

	// Pick a random text.
	chosen := possibleTemplates[rand.Intn(len(possibleTemplates))]
	text := chosen.template

	// Relplace each token one by one until we can't find any more.
	replacementsUsed := make(map[string]int)
	for _, token := range chosen.extractedTokens {
		// Replace all tokens with the provided replacements or a random value.
		var replacement string
		if replacementsUsed[token.Token] < len(tokenReplacements[token.Token]) {
			replacement = tokenReplacements[token.Token][replacementsUsed[token.Token]]
			replacementsUsed[token.Token]++
		} else {
			replacement = randArrayString(tokenRandom[token.Token])
		}

		// Apply modifiers.
		replacement = applyModifiers(replacement, token.Modifiers)

		// Replace the token with the replacement.
		text = strings.Replace(text, token.FullToken, replacement, 1)
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

// DefaultModifiers is a map of default modifiers that can be used in templates.
var DefaultModifiers = map[string]func(string) string{
	"title": func(s string) string {
		return strings.Title(s)
	},
	"upper": func(s string) string {
		return strings.ToUpper(s)
	},
	"lower": func(s string) string {
		return strings.ToLower(s)
	},
	"quote": func(s string) string {
		return fmt.Sprintf("'%s'", s)
	},
	"doublequote": func(s string) string {
		return fmt.Sprintf("%q", s)
	},
	"trimvowels": func(s string) string {
		return genlanguage.TrimVowels(s, 7)
	},
	"adjecive": func(s string) string {
		return genlanguage.GetAdjective(s)
	},
	"nounplural": func(s string) string {
		return genlanguage.GetNounPlural(s)
	},
	"a": func(s string) string {
		return genlanguage.GetArticle(s) + " " + s
	},
	"past": func(s string) string {
		return genlanguage.GetPastTense(s)
	},
	"presentsingular": func(s string) string {
		return genlanguage.GetPresentSingular(s)
	},
	"presentparticiple": func(s string) string {
		return genlanguage.GetPresentParticiple(s)
	},
}
