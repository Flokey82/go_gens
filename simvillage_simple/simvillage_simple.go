package simvillage_simple

import (
	"fmt"
	"log"
	"math/rand"
	"sort"

	"github.com/s0rg/fantasyname"
)

type Village struct {
	People []*Person
	maxID  int
	tick   int

	// Food int
	// Wood int
	firstGen [2]fmt.Stringer
	lastGen  fmt.Stringer
}

const firstNamePrefix = "!(bil|bal|ban|hil|ham|hal|hol|hob|wil|me|or|ol|od|gor|for|fos|tol|ar|fin|ere|leo|vi|bi|bren|thor)"

func New() *Village {
	v := new(Village)

	// Female first names.
	genFirstF, err := fantasyname.Compile(firstNamePrefix+"(|ga|orbise|apola|adure|mosi|ri|i|na|olea|ne)", fantasyname.Collapse(true), fantasyname.RandFn(rand.Intn))
	if err != nil {
		log.Fatal(err)
	}
	v.firstGen[0] = genFirstF

	// Male first names.
	genFirstM, err := fantasyname.Compile(firstNamePrefix+"(|go|orbis|apol|adur|mos|ole|n)", fantasyname.Collapse(true), fantasyname.RandFn(rand.Intn))
	if err != nil {
		log.Fatal(err)
	}
	v.firstGen[1] = genFirstM

	genLast, err := fantasyname.Compile("!BsVc", fantasyname.Collapse(true), fantasyname.RandFn(rand.Intn))
	if err != nil {
		log.Fatal(err)
	}
	v.lastGen = genLast
	return v
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
	v.popAge()

	// Calc required resources.
	// - if not enough food: Famine
	// - if not enough wood: Freeze in winter, increased disease chance

	// Assign jobs to jobless.
	// - determine which jobs are needed

	// Do job.
	// - produce resources

	// Pair up singles.
	v.popMatchMaker()

	// Pregnancies, settlers etc.
	v.popGrowth()

	// Random deaths.
	v.popDeath()
}

func (v *Village) popAge() {
	for _, p := range v.People {
		p.age++
	}
}

func (v *Village) popMatchMaker() {
	// Get eligible singles.
	var single []*Person
	for _, p := range v.People {
		if p.isEligibleSingle() {
			single = append(single, p)
		}
	}

	// Pair up singles.
	sort.Slice(single, func(a, b int) bool {
		return single[a].age > single[b].age
	})
	for i, p := range single {
		if !p.isEligibleSingle() {
			continue // Not single anymore.
		}
		for j, pc := range single {
			if !pc.isEligibleSingle() {
				continue // Not single anymore.
			}

			// TODO: Allow same sex couples (which can adopt children/orphans).
			if i == j || p.gender == pc.gender || isRelated(p, pc) {
				continue
			}

			// At most 33% age difference.
			if absInt(p.age-pc.age) > minInt(p.age, pc.age)/3 {
				continue
			}
			p.spouse = pc
			pc.spouse = p
			// Update family name.
			// TODO: This is not optimal... There should be a better way to do this.
			if p.gender == GenderFemale {
				p.lastName = pc.lastName
			} else {
				pc.lastName = p.lastName
			}
			log.Println(p.String(), "and", pc.String(), "are in love")
			break
		}
	}
}

func (v *Village) popGrowth() {
	// Pregnancies.
	var children []*Person
	for _, p := range v.People {
		if p.canBePregnant() {
			if rand.Intn(20) < 3 {
				c := v.newPerson()
				c.lastName = p.lastName
				c.mother = p
				c.father = p.spouse
				p.pregnantWith = c
			}
		} else if p.pregnantWith != nil {
			p.pregnant++

			// Give birth if we are far along enough.
			if p.pregnant > 10 {
				c := p.pregnantWith
				p.pregnantWith = nil
				p.pregnant = 0
				c.mother.children = append(c.mother.children, c)
				c.father.children = append(c.father.children, c)
				children = append(children, c)
				log.Println(c.mother.String(), "and", c.father.String(), "had a baby")
			}
		}
	}
	v.People = append(v.People, children...)

	// Random arrivals.
	if rand.Intn(10) < 2 {
		p := v.newPerson()
		p.lastName = v.lastGen.String()
		p.age = rand.Intn(20) + 16
		v.People = append(v.People, p)
		log.Println(p.String(), "arrived")
	}
}

func (v *Village) popDeath() {
	// TODO:
	// - Increase chances of death with rising population through disease.
	// - Increase chances of death if there is a famine or other states.
	var livingPeople []*Person
	for _, p := range v.People {
		// absolute value
		// - lowest chance at 30 years old (2 in 80)
		// - high child mortality (2 in 50)
		// - increasing mortality at > 60 years old (2 in < 50)
		chance := absInt(p.age - 30)
		if rand.Intn(80-chance) < 2 {
			// Kill villager.
			p.dead = true
			if spouse := p.spouse; spouse != nil {
				spouse.spouse = nil // Remove dead spouse from spouse.
			}
			log.Println(p.String(), "died and has", len(p.children), "children")
		} else {
			// Filter out dead people.
			livingPeople = append(livingPeople, p)
		}
	}
	v.People = livingPeople
}

func (v *Village) newPerson() *Person {
	p := &Person{
		id:     v.getNextID(),
		gender: randGender(),
	}
	p.firstName = v.firstGen[p.gender].String()
	return p
}

const (
	GenderFemale = iota
	GenderMale
)

func randGender() int {
	return rand.Intn(2)
}

type Person struct {
	id           int
	firstName    string
	lastName     string
	age          int
	dead         bool
	mother       *Person
	father       *Person
	spouse       *Person // TODO: keep track of former spouses?
	children     []*Person
	gender       int
	pregnant     int
	pregnantWith *Person
}

func (p *Person) Name() string {
	return p.firstName + " " + p.lastName
}

func (p *Person) String() string {
	gen := "F"
	if p.gender == GenderMale {
		gen = "M"
	}
	return p.Name() + fmt.Sprintf(" (%d%s)", p.age, gen)
}

func (p *Person) isEligibleSingle() bool {
	// Old enough and single.
	return p.age > 16 && p.spouse == nil
}

func (p *Person) canBePregnant() bool {
	// Female, has a spouse (implies old enough), and is currently not pregnant.
	// TODO: Set randomized upper age limit.
	return p.gender == GenderFemale && p.spouse != nil && p.pregnantWith == nil
}

func isRelated(a, b *Person) bool {
	if a == b.father || a == b.mother || b == a.father || b == a.mother {
		return true
	}
	if (a.father == nil && a.mother == nil) || (b.father == nil && b.mother == nil) {
		return false
	}
	return a.mother == b.mother || a.father == b.father
}

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
