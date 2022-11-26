package simvillage_simple

import (
	"fmt"
	"log"
	"math/rand"
	"sort"

	"github.com/Flokey82/genetics"
	"github.com/Flokey82/genetics/geneticshuman"
	"github.com/Flokey82/go_gens/gameconstants"
	"github.com/Flokey82/go_gens/utils"
	"github.com/s0rg/fantasyname"
)

type Village struct {
	People []*Person
	maxID  int
	tick   int
	day    int
	year   int

	// Food int
	// Wood int
	firstGen [2]fmt.Stringer
	lastGen  fmt.Stringer
}

const firstNamePrefix = "!(bil|bal|ban|hil|ham|hal|hol|hob|wil|me|or|ol|od|gor|for|fos|tol|ar|fin|ere|leo|vi|bi|bren|thor)"

// New returns a new village.
func New() *Village {
	v := new(Village)

	// Initialize name generation.
	genFirstF, err := fantasyname.Compile(firstNamePrefix+"(|ga|orbise|apola|adure|mosi|ri|i|na|olea|ne)", fantasyname.Collapse(true), fantasyname.RandFn(rand.Intn))
	if err != nil {
		log.Fatal(err)
	}
	v.firstGen[0] = genFirstF // Female first names.

	genFirstM, err := fantasyname.Compile(firstNamePrefix+"(|go|orbis|apol|adur|mos|ole|n)", fantasyname.Collapse(true), fantasyname.RandFn(rand.Intn))
	if err != nil {
		log.Fatal(err)
	}
	v.firstGen[1] = genFirstM // Male first names.

	genLast, err := fantasyname.Compile("!BsVc", fantasyname.Collapse(true), fantasyname.RandFn(rand.Intn))
	if err != nil {
		log.Fatal(err)
	}
	v.lastGen = genLast // Last names.
	return v
}

// getNextID returns the next unique ID.
func (v *Village) getNextID() int {
	id := v.maxID
	v.maxID++
	return id
}

// Tick advances the simulation by one day.
func (v *Village) Tick() {
	v.tick++
	v.year += (v.day + 1) / 365
	v.day = (v.day + 1) % 365
	log.Println("Tick", v.tick, "day", v.day, "year", v.year, "Population", len(v.People))

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

// popAge ages the population.
func (v *Village) popAge() {
	for _, p := range v.People {
		if p.bday == v.day {
			p.age++
			log.Println(p.String(), "has a birthday!")
		}
	}
}

// popMatchMaker pairs up singles.
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
			if utils.Abs(p.age-pc.age) > utils.Min(p.age, pc.age)/3 {
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

// popGrowth triggers pregnancies, births, and random settlers arriving at the village.
func (v *Village) popGrowth() {
	// Pregnancies.
	var children []*Person
	for _, p := range v.People {
		if p.canBePregnant() {
			// Approximately once every 4 years if no children.
			// TODO: Figure out proper chance of birth.
			chance := 4 * 365
			if p.age > 40 {
				// Over 40, it becomes more and more unlikely.
				// TODO: Genetic variance?
				chance *= (p.age - 40)
			}

			// The more children, the less likely it becomes
			// that more children are on the way.
			// NOTE: Not because of biological reasons, but
			// who wants more children after having some.
			chance *= p.numLivingChildren() + 1

			if rand.Intn(chance) < 1 {
				c := v.newPerson()
				c.lastName = p.lastName
				c.mother = p
				c.father = p.spouse
				p.pregnantWith = c
			}
		} else if p.pregnantWith != nil {
			p.pregnant++

			// Give birth if we are far along enough.
			if p.pregnant > 9*30 { // TODO: Add some variance
				c := p.pregnantWith
				p.pregnantWith = nil
				p.pregnant = 0
				c.mother.children = append(c.mother.children, c)
				c.father.children = append(c.father.children, c)
				c.g = genetics.Mix(c.mother.g, c.father.g, 2)
				c.fixGenes()
				c.bday = v.day // Birthday!
				children = append(children, c)
				log.Println(c.mother.String(), "\n", geneticshuman.String(c.mother.g), "\nand", c.father.String(), "\n", geneticshuman.String(c.father.g), "\nhad a baby\n", geneticshuman.String(c.g))
			}
		}
	}
	v.People = append(v.People, children...)

	// Random arrival of settlers.
	if rand.Intn(365) < 1 {
		v.AddRandomPerson()
	}
}

// AddRandomPerson adds a random settler to the village.
func (v *Village) AddRandomPerson() {
	p := v.newPerson()
	p.lastName = v.lastGen.String()
	p.age = rand.Intn(20) + 16
	p.bday = rand.Intn(365)
	p.g = genetics.NewRandom()
	p.fixGenes()
	v.People = append(v.People, p)

	log.Println(p.String(), "arrived")
}

// popDeath causes the death random people based on their age.
func (v *Village) popDeath() {
	// TODO:
	// - Increase chances of death with rising population through disease.
	// - Increase chances of death if there is a famine or other states.
	var livingPeople []*Person
	for _, p := range v.People {
		// Check if villager dies of natural causes.
		p.dead = gameconstants.DiesAtAge(p.age)

		// Filter out dead people.
		if !p.dead {
			livingPeople = append(livingPeople, p)
			continue
		}

		// Kill villager.
		if spouse := p.spouse; spouse != nil {
			spouse.spouse = nil // Remove dead spouse from spouse.
		}
		// TODO: Remove child from parents?
		log.Println(p.String(), "died and has", len(p.children), "children !!!!!!!!", p.numLivingChildren(), "alive")
		for _, c := range p.children {
			log.Println(c.String())
		}
	}
	v.People = livingPeople
}

// Gender represents a gender.
type Gender int

const (
	GenderFemale Gender = iota
	GenderMale
)

// String returns the string representation of the gender.
func (g Gender) String() string {
	switch g {
	case GenderFemale:
		return "F"
	case GenderMale:
		return "M"
	default:
		return "X"
	}
}

// randGender returns a random gender.
func randGender() Gender {
	return Gender(rand.Intn(2))
}

// Person represents a person in the village.
type Person struct {
	id           int
	firstName    string
	lastName     string
	age          int
	bday         int
	dead         bool
	mother       *Person
	father       *Person
	spouse       *Person // TODO: keep track of spouses that might have perished?
	children     []*Person
	gender       Gender
	pregnant     int
	pregnantWith *Person
	g            genetics.Genes
}

// newPerson creates a new person.
func (v *Village) newPerson() *Person {
	p := &Person{
		id:     v.getNextID(),
		gender: randGender(),
	}
	p.firstName = v.firstGen[p.gender].String()
	return p
}

// Name returns the name of the person.
func (p *Person) Name() string {
	return p.firstName + " " + p.lastName
}

// String returns the string representation of the person.
func (p *Person) String() string {
	deadStr := ""
	if p.dead {
		deadStr = " dead"
	}
	return p.Name() + fmt.Sprintf(" (%d %s%s)", p.age, p.gender, deadStr)
}

// numLivingChildren returns the number of children that are still alive.
func (p *Person) numLivingChildren() int {
	var n int
	for _, c := range p.children {
		if !c.dead {
			n++
		}
	}
	return n
}

// isElegibleSingle returns true if the person is old enough and single.
func (p *Person) isEligibleSingle() bool {
	// Old enough and single.
	return p.age > 16 && p.spouse == nil
}

// canBePregnant returns true if the person is old enough and not pregnant.
func (p *Person) canBePregnant() bool {
	// Female, has a spouse (implies old enough), and is currently not pregnant.
	// TODO: Set randomized upper age limit.
	return p.gender == GenderFemale && p.spouse != nil && p.pregnantWith == nil
}

// fixGenes makes sure that the gender is set properly.
// NOTE: This needs to be done due to the genetics package being a bit weird.
func (p *Person) fixGenes() {
	switch p.gender {
	case GenderFemale:
		geneticshuman.SetGender(&p.g, geneticshuman.GenderFemale)
	case GenderMale:
		geneticshuman.SetGender(&p.g, geneticshuman.GenderMale)
	default:
		geneticshuman.SetGender(&p.g, 0)
	}
}

// isRelated returns true if a and b are related (first degree).
func isRelated(a, b *Person) bool {
	if a == b.father || a == b.mother || b == a.father || b == a.mother {
		return true
	}
	if (a.father == nil && a.mother == nil) || (b.father == nil && b.mother == nil) {
		return false
	}
	return a.mother == b.mother || a.father == b.father
}
