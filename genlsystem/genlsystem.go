package genlsystem

// Lindenmayer iterative function.
// See: http://en.wikipedia.org/wiki/L-system
func Lindenmayer(start []string, rules map[string][]string, n int) []string {
	if n <= 0 {
		return start // If we ran out of iterations, return the start.
	}

	// Apply rules.
	var result []string
	for _, c := range start {
		if r, ok := rules[c]; ok {
			result = append(result, r...)
			continue
		}
		result = append(result, c)
	}
	return Lindenmayer(result, rules, n-1)
}

// TODO: https://jsantell.com/l-systems/#3d-interpretation
