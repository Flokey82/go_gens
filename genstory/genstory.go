// Package genstory provides means to generate text from templates and / or rules.
package genstory

import (
	"errors"
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
	UseAllProvided   bool                           // Use all provided tokens, even if they are not used in the template.
	Modifiers        map[string]func(string) string // A map of token names to a function that modifies the replacement text.
}

// Generate generates a text from the provided tokens and the configuration.
func (c *TextConfig) Generate(provided []TokenReplacement) (*Generated, error) {
	return DefaultTextGenerator.GenerateFromConfig(provided, c, nil)
}

// GenerateAndGiveMeTheTemplate generates a text from the provided tokens and the configuration.
// It also returns the template that was used to generate the text.
func (c *TextConfig) GenerateAndGiveMeTheTemplate(provided []TokenReplacement) (*Generated, error) {
	return DefaultTextGenerator.GenerateFromConfig(provided, c, nil)
}

// GenerateWithTemplate generates a text from the provided tokens and the provided template.
func (c *TextConfig) GenerateWithTemplate(provided []TokenReplacement, template string) (*Generated, error) {
	return DefaultTextGenerator.GenerateFromConfig(provided, c, []string{template})
}

// TextGenerator is a generator for text using TextConfigs.
// Using this over the TextConfig methods allows you to use a custom random number generator
// and to (re-)set the random number generator's seed.
type TextGenerator struct {
	RandInterface
}

// NewTextGenerator creates a new TextGenerator using the provided random number generator.
func NewTextGenerator(rng RandInterface) *TextGenerator {
	return &TextGenerator{
		RandInterface: rng,
	}
}

// Generate generates a text from the provided tokens and the configuration.
func (g *TextGenerator) Generate(provided []TokenReplacement, config *TextConfig) (*Generated, error) {
	return g.GenerateFromConfig(provided, config, nil)
}

// GenerateAndGiveMeTheTemplate generates a text from the provided tokens and the configuration.
// It also returns the template that was used to generate the text.
func (g *TextGenerator) GenerateAndGiveMeTheTemplate(provided []TokenReplacement, config *TextConfig) (*Generated, error) {
	return g.GenerateFromConfig(provided, config, nil)
}

// GenerateButUseThisTemplate generates a text from the provided tokens and the provided template.
func (g *TextGenerator) GenerateButUseThisTemplate(provided []TokenReplacement, config *TextConfig, template string) (*Generated, error) {
	return g.GenerateFromConfig(provided, config, []string{template})
}

// Generated provides information about a generated text.
type Generated struct {
	Text     string
	Template string
	Tokens   []TokenReplacement
}

func (g *TextGenerator) GenerateFromConfig(provided []TokenReplacement, config *TextConfig, altTemplates []string) (*Generated, error) {
	templates := config.Templates
	if altTemplates != nil {
		templates = altTemplates
	}
	tokens := config.Tokens
	isMandatory := config.TokenIsMandatory
	tokenRandom := config.TokenPools
	modifierFuncs := config.Modifiers
	capitalize := config.Title
	useAllProvided := config.UseAllProvided

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
			return nil, err
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
		if useAllProvided {
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
		}

		// We have all required tokens, add the template to the list of possible templates.
		possibleTemplates = append(possibleTemplates, candidate{
			template:        template,
			extractedTokens: extracted,
		})
	}

	// If we have no possible templates, return an error.
	if len(possibleTemplates) == 0 {
		return nil, errors.New("no possible templates satisfying the provided tokens")
	}

	// Pick a random text.
	chosen := possibleTemplates[rand.Intn(len(possibleTemplates))]
	generated := &Generated{
		Text:     chosen.template,
		Template: chosen.template,
	}

	// Relplace each token one by one until we can't find any more.
	replacementsUsed := make(map[string]int)
	for _, token := range chosen.extractedTokens {
		// Replace all tokens with the provided replacements or a random value.
		var replacement string
		if replacementsUsed[token.Token] < len(tokenReplacements[token.Token]) {
			replacement = tokenReplacements[token.Token][replacementsUsed[token.Token]]
			replacementsUsed[token.Token]++
		} else {
			replacement = randArrayString(g, tokenRandom[token.Token])
		}

		// Remember the token and the replacement we used.
		generated.Tokens = append(generated.Tokens, TokenReplacement{
			Token:       token.Token,
			Replacement: replacement,
		})

		// Apply modifiers.
		replacement = applyModifiers(replacement, token.Modifiers)

		// Replace the token with the replacement.
		generated.Text = strings.Replace(generated.Text, token.FullToken, replacement, 1)
	}
	if capitalize {
		// Capitalize the first letter of each word in the text.
		generated.Text = strings.Title(generated.Text)
	} else {
		// Capitalize the first letter of the text.
		// We have to make sure we don't just use the slice operator, since that
		// might corrupt UTF-8 characters.
		generated.Text = genlanguage.Capitalize(generated.Text)
	}

	return generated, nil
}

func randArrayString(rng RandInterface, arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	return arr[rng.Intn(len(arr))]
}

func ApplyModifiers(s string, modifiers []string) string {
	for _, modifier := range modifiers {
		if fn, ok := DefaultModifiers[modifier]; ok {
			s = fn(s)
		}
	}
	return s
}

// DefaultModifiers is a map of default modifiers that can be used in templates.
var DefaultModifiers = map[string]func(string) string{
	"title":             strings.Title,          // Title case.
	"capitalize":        genlanguage.Capitalize, // Capitalize the first letter.
	"upper":             strings.ToUpper,
	"lower":             strings.ToLower,
	"adjecive":          genlanguage.GetAdjective,
	"nounplural":        genlanguage.GetNounPlural,
	"past":              genlanguage.GetPastTense,
	"presentsingular":   genlanguage.GetPresentSingular,
	"presentparticiple": genlanguage.GetPresentParticiple,
	"quote": func(s string) string {
		return "'" + s + "'"
	},
	"doublequote": func(s string) string {
		return "\"" + s + "\""
	},
	"trimvowels": func(s string) string {
		return genlanguage.TrimVowels(s, 3)
	},
	"a": func(s string) string {
		// Add an article to the string.
		return genlanguage.GetArticle(s) + " " + s
	},
}

type RandInterface interface {
	Intn(n int) int
	Seed(seed int64)
}

type randWrapper struct{}

func (r *randWrapper) Intn(n int) int {
	return rand.Intn(n)
}

func (r *randWrapper) Seed(seed int64) {
	rand.Seed(seed)
}

var defaultRand = &randWrapper{}

var DefaultTextGenerator = NewTextGenerator(defaultRand)
