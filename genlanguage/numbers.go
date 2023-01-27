package genlanguage

var (
	numBelow20    = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}
	numTens       = []string{"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}
	numMagnitudes = []string{
		"",
		"thousand",
		"million",
		"billion",
		"trillion",
		"quadrillion",
		"quintillion",
		"sextillion",
		"septillion",
		"octillion",
		"nonillion",
		"decillion",
		"undecillion",
		"duodecillion",
		"tredecillion",
	}
)

// NumberToWords returns a string representation of the number in English.
func NumberToWords(num int) string {
	if num < 0 {
		return "negative " + NumberToWords(-num)
	}

	if num < 20 {
		return numBelow20[num]
	}

	if num < 100 {
		if num%10 == 0 {
			return numTens[num/10]
		}
		return numTens[num/10] + "-" + NumberToWords(num%10)
	}

	// Hundreds.
	if num < 1000 {
		if num%100 == 0 {
			return NumberToWords(num/100) + " hundred"
		}
		return NumberToWords(num/100) + " hundred and " + NumberToWords(num%100)
	}

	// Thousands and above... loop through the magnitudes and assemble the result.
	var result string
	for i := 0; num > 0; i++ {
		if num%1000 != 0 {
			nextMag := NumberToWords(num % 1000)
			if i > 0 {
				nextMag += " " + numMagnitudes[i]
			} else if num%1000 < 100 {
				nextMag = "and " + nextMag
			}
			if result == "" {
				result = nextMag
			} else {
				result = nextMag + " " + result
			}
		}
		num /= 1000
	}
	return result
}
