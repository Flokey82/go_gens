package simvillage

import (
	"math/rand"
	"strings"
)

var vowel = []string{"a", "e", "i", "o", "u"}
var consonant = []string{"b", "c", "d", "f", "g", "h", "j",
	"k", "l", "m", "n", "p", "r", "s", "t"}

// Input: Structure of a word
// Output: Word generated from said structure

func name_gen() string {
	var new_name string
	for x := 0; x <= 10; x++ { // This algo uses 1/3 vowel and 2/3 consts
		if (rand.Float64() * 3) > 1 {
			new_name += consonant[int(float64(len(consonant))*rand.Float64())]
		} else {
			new_name += vowel[int(float64(len(vowel))*rand.Float64())]
		}
	}
	return new_name
}

func get_name() string {
	if rand.Intn(1) > 0 {
		if rand.Intn(4) > 3 {
			return strings.ToTitle(name_cvcv()) + name_cvcv()
		} else {
			return strings.ToTitle(name_cvcv())
		}
	} else {
		if rand.Intn(4) > 2 {
			return strings.ToTitle(name_cvcvc())
		} else {
			return strings.ToTitle(name_cvcvc())
		}
	}
}

func get_first_name() string {
	return strings.ToTitle(name_cvcv())
}

func get_last_name() string {
	return strings.ToTitle(name_cvcv()) + name_cvcv()
}

func name_cvcv() string {
	//name_cvcv = ""
	return (add_const() + add_vowel() + add_const() + add_vowel())
}

func name_cvcvc() string {
	//name_cvcv = ""
	return (add_const() + add_vowel() + add_const() + add_vowel() + add_const())
}

// Generates a word sructure in the form of a string
func word_struc() string {
	return add_const() + add_vowel() + add_const() + add_const()
}

// adds a vowel

func add_vowel() string {
	return vowel[(int(float64(len(vowel)) * rand.Float64()))]
}

// adds a constonant
func add_const() string {
	return consonant[(int(float64(len(consonant)) * rand.Float64()))]
}
