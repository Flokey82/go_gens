package simvillage_simple

import (
	"fmt"
	"log"
	"math/rand"
	"sort"

	"github.com/Flokey82/genetics"
	"github.com/Flokey82/genetics/geneticshuman"
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

func (v *Village) popAge() {
	for _, p := range v.People {
		if p.bday == v.day {
			p.age++
			log.Println(p.String(), "has a birthday!")
		}
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
				fixGenes(c)
				c.bday = v.day // Birthday!
				children = append(children, c)
				log.Println(c.mother.String(), "\n", geneticshuman.String(c.mother.g), "\nand", c.father.String(), "\n", geneticshuman.String(c.father.g), "\nhad a baby\n", geneticshuman.String(c.g))
			}
		}
	}
	v.People = append(v.People, children...)

	// Random arrivals.
	if rand.Intn(365) < 1 {
		v.AddRandomPerson()
	}
}

func fixGenes(p *Person) {
	switch p.gender {
	case GenderFemale:
		geneticshuman.SetGender(&p.g, geneticshuman.GenderFemale)
	case GenderMale:
		geneticshuman.SetGender(&p.g, geneticshuman.GenderMale)
	default:
		geneticshuman.SetGender(&p.g, 0)
	}
}

func (v *Village) AddRandomPerson() {
	p := v.newPerson()
	p.lastName = v.lastGen.String()
	p.age = rand.Intn(20) + 16
	p.bday = rand.Intn(365)
	p.g = genetics.NewRandom()
	fixGenes(p)
	v.People = append(v.People, p)

	log.Println(p.String(), "arrived")
}

func (v *Village) popDeath() {
	// TODO:
	// - Increase chances of death with rising population through disease.
	// - Increase chances of death if there is a famine or other states.
	var livingPeople []*Person
	for _, p := range v.People {
		// TODO: Figure out proper chance of death.
		// TODO: Child mortality?
		// From: https://github.com/Kontari/Village/blob/master/src/death.py
		if 35 < p.age && p.age < 50 { // Adult
			p.dead = rand.Intn(241995) == 0
		} else if 50 < p.age && p.age < 70 { // Old Person
			p.dead = rand.Intn(29380579) == 0
		} else if p.age > 70 { // Elderly
			p.dead = rand.Intn(5475) == 0
		}
		if p.dead {
			// Kill villager.
			if spouse := p.spouse; spouse != nil {
				spouse.spouse = nil // Remove dead spouse from spouse.
			}
			// TODO: Remove child from parents?
			log.Println(p.String(), "died and has", len(p.children), "children !!!!!!!!", p.numLivingChildren(), "alive")
			for _, c := range p.children {
				log.Println(c.String())
			}
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

type Gender int

const (
	GenderFemale Gender = iota
	GenderMale
)

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

func randGender() Gender {
	return Gender(rand.Intn(2))
}

type Person struct {
	id           int
	firstName    string
	lastName     string
	age          int
	bday         int
	dead         bool
	mother       *Person
	father       *Person
	spouse       *Person // TODO: keep track of former spouses?
	children     []*Person
	gender       Gender
	pregnant     int
	pregnantWith *Person
	g            genetics.Genes
}

func (p *Person) Name() string {
	return p.firstName + " " + p.lastName
}

func (p *Person) String() string {
	deadStr := ""
	if p.dead {
		deadStr = " dead"
	}
	return p.Name() + fmt.Sprintf(" (%d %s%s)", p.age, p.gender, deadStr)
}

func (p *Person) numLivingChildren() int {
	var n int
	for _, c := range p.children {
		if !c.dead {
			n++
		}
	}
	return n
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
