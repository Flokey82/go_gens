package genstory

import "fmt"

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
	// Scan through the string, looking for tokens.
	var tokens []ExtractedToken
	var tokenStart int
	var tokenEnd int // Either the end of the token or the start of the modifiers.
	var inToken bool

	// Look for modifiers.
	var modifiers []string
	var inModifier bool
	var modifierStart int

	appendModifier := func(pos int) {
		if !inModifier {
			return
		}
		modifiers = append(modifiers, s[modifierStart:pos])
		inModifier = false
	}

	for i, c := range s {
		if c == '[' {
			if inToken {
				return nil, fmt.Errorf("unexpected token start at %d", i)
			}
			inToken = true
			tokenStart = i
		} else if c == ']' {
			if !inToken {
				return nil, fmt.Errorf("unexpected token end at %d", i)
			}
			if inModifier {
				appendModifier(i)
			} else {
				tokenEnd = i
			}
			inToken = false
			tokens = append(tokens, ExtractedToken{
				Token:     s[tokenStart:tokenEnd] + "]",
				FullToken: s[tokenStart : i+1],
				Modifiers: modifiers,
				Start:     tokenStart,
				End:       i + 1,
			})
			modifiers = nil
		} else if c == ':' && inToken {
			if inModifier {
				appendModifier(i)
			} else {
				tokenEnd = i
			}
			inModifier = true
			modifierStart = i + 1
		}

	}
	if inToken {
		return nil, fmt.Errorf("unexpected end of string")
	}
	return tokens, nil
}
