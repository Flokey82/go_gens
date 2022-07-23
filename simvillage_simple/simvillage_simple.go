package simvillage_simple

import (
	"log"
	"math/rand"
)

type Village struct {
	People []*Person
	maxID  int
	tick   int
}

func New() *Village {
	return new(Village)
}

func (v *Village) getNextID() int {
	id := v.maxID
	v.maxID++
	return id
}

func (v *Village) Tick() {
	log.Println("Tick", v.tick, "Population", len(v.People))
	v.tick++

	// Increase age of villagers.
	for _, p := range v.People {
		p.age++
	}

	// Get eligible singles.
	var singleMen, singleWomen []*Person
	for _, p := range v.People {
		if p.isEligibleSingle() {
			switch p.gender {
			case GenderFemale:
				singleWomen = append(singleWomen, p)
			case GenderMale:
				singleMen = append(singleMen, p)
			}
		}
	}

	// Pair up singles.
	maxLen := len(singleMen)
	if len(singleWomen) < maxLen {
		maxLen = len(singleWomen)
	}

	// TODO: Shuffle the slices first.
	// TODO: Ensure that siblings don't marry...?
	for i := 0; i < maxLen; i++ {
		singleMen[i].spouse = singleWomen[i]
		singleWomen[i].spouse = singleMen[i]
		log.Println(singleWomen[i].id, "and", singleMen[i].id, "are in love")
	}

	// Random pregnancies.
	var children []*Person
	for _, p := range v.People {
		if p.canBePregnant() {
			if rand.Intn(20) < 3 {
				p.pregnant++
			}
		} else if p.pregnant > 0 {
			p.pregnant++

			// Give birth if we are far along enough.
			if p.pregnant > 10 {
				p.pregnant = 0
				children = append(children, &Person{
					id:     v.getNextID(),
					gender: randGender(),
				})
				log.Println(p.id, "had a baby")
			}
		}
	}
	v.People = append(v.People, children...)

	// Random deaths.
	// TODO: Increase chances of death with rising population through disease.
	var livingPeople []*Person
	for _, p := range v.People {
		// absolute value
		// - lowest chance at 20 years old (2 in 80)
		// - high child mortality (2 in 60)
		// - increasing mortality at > 40 years old (2 in < 60)
		chance := p.age - 20
		if chance != 0 {
			chance = chance * chance / chance
		}
		if rand.Intn(80-chance) < 2 {
			// Kill vilager.
			if spouse := p.spouse; spouse != nil {
				spouse.spouse = nil // Remove dead spouse from spouse.
			}
			log.Println(p.id, "died at age", p.age)
		} else {
			// Filter out dead people.
			livingPeople = append(livingPeople, p)
		}
	}
	v.People = livingPeople

	// Random arrival.
	if rand.Intn(10) < 2 {
		p := &Person{
			id:     v.getNextID(),
			gender: randGender(),
			age:    rand.Intn(10) + 16,
		}
		v.People = append(v.People, p)
		log.Println(p.id, "arrived")
	}
}

const (
	GenderFemale = iota
	GenderMale
)

func randGender() int {
	return rand.Intn(2)
}

type Person struct {
	id       int
	age      int
	spouse   *Person
	gender   int
	pregnant int
}

func (p *Person) isEligibleSingle() bool {
	// Old enough and single.
	return p.age > 16 && p.spouse == nil
}

func (p *Person) canBePregnant() bool {
	// Female, has a spouse (implies old enough), and is currently not pregnant.
	return p.gender == GenderFemale && p.spouse != nil && p.pregnant == 0
}
