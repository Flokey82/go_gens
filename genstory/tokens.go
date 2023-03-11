package genstory

import (
	"strings"
)

// ExtractedToken is a token extracted from a string.
type ExtractedToken struct {
	Token     string   // The name of the token.
	FullToken string   // The full token, including the square brackets and modifiers.
	Modifiers []string // A list of modifiers for the token.
	Start     int      // The start position of the token in the string.
	End       int      // The end position of the token in the string.
}

// ExtractTokens extracts all tokens from a string.
// A token is a string surrounded by square brackets ('[]')
// and can have modifiers, separated by a colon (':').
// Example: "[token:upper:quote]"
func ExtractTokens(s string) ([]ExtractedToken, error) {
	var tokens []ExtractedToken
	for _, t := range findAllTokens(s) {
		token := ExtractToken(t.token)
		token.Start = t.start
		token.End = t.end
		tokens = append(tokens, token)
	}
	return tokens, nil
}

// ExtractToken extracts a single token from a string.
func ExtractToken(s string) ExtractedToken {
	// Trim off the brackets.
	token := s[1 : len(s)-1]

	// Split off all modifiers.
	modifiers := strings.Split(token, ":")
	token = modifiers[0]
	modifiers = modifiers[1:]

	return ExtractedToken{
		Token:     "[" + token + "]", // Add the brackets back.
		FullToken: s,
		Modifiers: modifiers,
	}
}

type tokenLocation struct {
	token string
	start int
	end   int
}

func findAllTokens(s string) []tokenLocation {
	var tokens []tokenLocation
	var inToken bool
	var start int

	for i, c := range s {
		if c == '[' {
			if inToken {
				continue
			}
			inToken = true
			start = i
		} else if c == ']' {
			if !inToken {
				continue
			}
			inToken = false
			tokens = append(tokens, tokenLocation{
				token: s[start : i+1],
				start: start,
				end:   i + 1,
			})
		}
	}

	return tokens
}
