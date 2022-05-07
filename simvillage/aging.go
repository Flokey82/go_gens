package simvillage

import (
	"math/rand"
)

// Age manages birthdays, and age-related job eligability.
type Age struct {
	age      int
	bday     int
	log      []string
	age_text string
}

func NewAge(age int) *Age {
	a := &Age{
		age:  age,
		bday: rand.Intn(359) + 1, // r.randint(1, 360)
	}
	a.reassign_age_text()
	return a
}

func NewAgeWithBDay(age, bday int) *Age {
	a := &Age{
		age:  age,
		bday: bday,
	}
	a.reassign_age_text()
	return a
}

func (a *Age) Tick() []string {
	return nil
}

func (a *Age) reassign_age_text() {
	if a.age < 4 {
		a.age_text = "Infant"
	} else if 4 < a.age && a.age < 10 {
		a.age_text = "Child"
	} else if 15 < a.age && a.age < 24 {
		a.age_text = "Young Adult"
	} else if 25 < a.age && a.age < 50 {
		a.age_text = "Adult"
	} else if 50 < a.age && a.age < 70 {
		a.age_text = "Old Person"
	} else {
		a.age_text = "Elder"
	}
}
