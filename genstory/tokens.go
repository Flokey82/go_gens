package genstory

import "fmt"

// ExtractTokens extracts all tokens from a string.
// A token is a string surrounded by square brackets ('[]').
func ExtractTokens(s string) ([]string, error) {
	// Scan through the string, looking for tokens.
	var tokens []string
	var tokenStart int
	var inToken bool
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
			inToken = false
			tokens = append(tokens, s[tokenStart:i+1])
		}
	}
	if inToken {
		return nil, fmt.Errorf("unexpected end of string")
	}
	return tokens, nil
}
