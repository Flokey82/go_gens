package simvillage

import (
	"math/rand"
	"strings"
)

var vowel = []string{"a", "e", "i", "o", "u"}
var consonant = []string{"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "r", "s", "t"}

// Input: Structure of a word
// Output: Word generated from said structure

func nameGen() string {
	var newName string
	for x := 0; x <= 10; x++ { // This algo uses 1/3 vowel and 2/3 consts
		if (rand.Float64() * 3) > 1 {
			newName += pickRandString(consonant)
		} else {
			newName += pickRandString(vowel)
		}
	}
	return newName
}

func getName() string {
	if rand.Intn(2) > 0 {
		if rand.Intn(4) > 3 {
			return strings.ToTitle(nameCVCV()) + nameCVCV()
		} else {
			return strings.ToTitle(nameCVCV())
		}
	} else {
		if rand.Intn(4) > 2 {
			return strings.ToTitle(nameCVCVC())
		} else {
			return strings.ToTitle(nameCVCVC())
		}
	}
}

func getFirstName() string {
	return strings.ToTitle(nameCVCV())
}

func getLastName() string {
	return strings.ToTitle(nameCVCV()) + nameCVCV()
}

func nameCVCV() string {
	return (addConsonant() + addVowel() + addConsonant() + addVowel())
}

func nameCVCVC() string {
	return (addConsonant() + addVowel() + addConsonant() + addVowel() + addConsonant())
}

// wordStruc generates a word sructure in the form of a string
func wordStruc() string {
	return addConsonant() + addVowel() + addConsonant() + addConsonant()
}

// adds a vowel
func addVowel() string {
	return pickRandString(vowel)
}

// addConsonant adds a consonant.
func addConsonant() string {
	return pickRandString(consonant)
}
