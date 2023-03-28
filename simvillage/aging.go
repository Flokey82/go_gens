package simvillage

import (
	"math/rand"
)

const (
	AgeInfant     = "Infant"
	AgeChild      = "Child"
	AgeYoungAdult = "Young Adult"
	AgeAdult      = "Adult"
	AgeOldPerson  = "Old Person"
	AgeElder      = "Elder"
)

// Age manages birthdays, and age-related job eligability.
type Age struct {
	age      int
	bday     int
	age_text string
	log      []string
}

func NewAge(age int) *Age {
	a := &Age{
		age:  age,
		bday: rand.Intn(359) + 1,
	}
	a.reassignAgeText()
	return a
}

func NewAgeWithBDay(age, bday int) *Age {
	a := &Age{
		age:  age,
		bday: bday,
	}
	a.reassignAgeText()
	return a
}

func (a *Age) Tick() []string {
	return nil
}

func (a *Age) reassignAgeText() {
	if a.age < 4 {
		a.age_text = AgeInfant
	} else if 4 < a.age && a.age < 10 {
		a.age_text = AgeChild
	} else if 15 < a.age && a.age < 24 {
		a.age_text = AgeYoungAdult
	} else if 25 < a.age && a.age < 50 {
		a.age_text = AgeAdult
	} else if 50 < a.age && a.age < 70 {
		a.age_text = AgeOldPerson
	} else {
		a.age_text = AgeElder
	}
}
