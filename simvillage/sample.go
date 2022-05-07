package simvillage

import (
	"math/rand"
)

type RandomEffects struct {
}

func NewRandomEffects() *RandomEffects {
	return &RandomEffects{}
}

func (r *RandomEffects) get_mod() float64 {
	p := rand.Intn(100)
	if p >= 0 && p <= 80 {
		return 1.0
	} else if p >= 81 && p <= 95 {
		return 1.5
	} else if p >= 96 && p <= 100 {
		return 2.0
	} else {
		return 1.0
	}
}

// Returns output for a day of farming
func (r *RandomEffects) get_farming() float64 {
	p := rand.Intn(100)
	if p >= 0 && p <= 80 {
		return 1.0
	} else if p >= 81 && p <= 95 {
		return 1.5
	} else if p >= 96 && p <= 100 {
		return 2.0
	} else {
		return 1.0
	}
}
