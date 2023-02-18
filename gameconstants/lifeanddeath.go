package gameconstants

import (
	"math"
	"math/rand"
)

// https://pubmed.ncbi.nlm.nih.gov/28006969/
// Seems like in the middle ages, death rates were elevated in urban areas, especially in the winter.
// for women.
// https://sharondewitte.files.wordpress.com/2018/02/walter-and-dewitte-2017-urban-and-rural-mortality-and-survival-in-medieval-england.pdf
// ... has some interesting graphs like age at death etc.

// See: https://en.wikipedia.org/wiki/Medieval_demography#Demographic_tables_of_Europe%E2%80%99s_population
// That's just european though... I'd love to see other data, especially from
// the middle east and asia (different cultures, different climates).
const (
	AvgGrowthPerYearMin = 0.09
	AvgGrowthPerYearMax = 0.20
)

// For calculating the population growth, we'd have to take into account
// if there was a famine, war, etc. in the past years.
// For now we just assume a constant growth rate.
// 'startVal' is the starting population.
// 't' is the time that has passed
// 'growthRate' is the average growth rate in the time interval used to measure 't'.
// NOTE: If 't' is in years, then the growth rate should be in growth / year as well.
func CalcPopulationAfterNYears(startVal int, t, growthRate float64) float64 {
	return float64(startVal) * math.Pow(math.E, growthRate*t)
}

// ConvertGrowthRate converts a given growth rate to the growth rate of a given factor.
// For example, to convert yearly growth to daily growth, the factor would be 1/365.
func ConvertGrowthRate(factor, growthRate float64) float64 {
	return math.Pow(growthRate, factor) - 1
}

// DiesAtAge returns true if the person of the given age dies of natural causes.
// From: https://github.com/Kontari/Village/blob/master/src/death.py
// TODO: Mortality under 35 and child mortaity is not yet implemented.
// TODO: Figure out proper chance of death?
func DiesAtAge(age int) bool {
	if 35 < age && age < 50 { // Adult
		return rand.Intn(241995) == 0
	} else if 50 < age && age < 70 { // Old Person
		return rand.Intn(29380579) == 0
	} else if age > 70 { // Elderly
		return rand.Intn(5475) == 0
	}
	return false
}

// DiesAtAgeWithinNDays returns true if the person of the given age dies of
// natural causes within the next n days.
// TODO: Deduplicate code with DiesAtAge.
func DiesAtAgeWithinNDays(age, n int) bool {
	if 35 < age && age < 50 { // Adult
		return rand.Intn(241995) < n
	} else if 50 < age && age < 70 { // Old Person
		return rand.Intn(29380579) < n
	} else if age > 70 { // Elderly
		return rand.Intn(5475) < n
	}
	return false
}
