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
	People   []*Person       // All people in the village.
	maxID    int             // Current max unique ID.
	tick     int             // Current tick.
	day      int             // Current day.
	year     int             // Current year.
	firstGen [2]fmt.Stringer // First name generators (male/female).
	lastGen  fmt.Stringer    // Last name generators.
	// Food int
	// Wood int
}

// first name prefixes for fantasyname generator.
const firstNamePrefix = "!(bil|bal|ban|hil|ham|hal|hol|hob|wil|me|or|ol|od|gor|for|fos|tol|ar|fin|ere|leo|vi|bi|bren|thor)"

// New returns a new village.
func New() *Village {
	v := new(Village)

	// Initialize name generation.

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

	// Last names.
	genLast, err := fantasyname.Compile("!BsVc", fantasyname.Collapse(true), fantasyname.RandFn(rand.Intn))
	if err != nil {
		log.Fatal(err)
	}
	v.lastGen = genLast
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
	// Advance day and eventually year.
	v.tick++
	v.year += (v.day + 1) / 365
	v.day = (v.day + 1) % 365

	// Log tick.
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

	// Sort by age, so similar age people are more likely to be paired up quicker.
	sort.Slice(single, func(a, b int) bool {
		return single[a].age > single[b].age
	})

	// Pair up singles.
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
				// We are pregnant!
				// Generate a new person and add it to the mother.
				c := v.newPerson()
				c.lastName = p.lastName

				// Set the child's parents.
				c.mother = p
				c.father = p.spouse

				// Set the child's genes.
				c.genes = genetics.Mix(c.mother.genes, c.father.genes, 2)
				c.fixGenes()

				// Assign the child to the mother.
				p.pregnantWith = c
			}
		} else if p.pregnantWith != nil {
			p.pregnant++

			// Give birth if we are far along enough.
			// TODO: Add some variance.
			if p.pregnant > 9*30 {
				// The miracle of life!

				// Get the child and remove it from the mother.
				c := p.pregnantWith
				p.pregnantWith = nil
				p.pregnant = 0

				// Add the new born to the child's parents.
				c.mother.children = append(c.mother.children, c)
				c.father.children = append(c.father.children, c)

				// Set the child's birthday.
				c.bday = v.day

				// Log the joyous event.
				// LOL you will not sleep for a while, enjoy! :P
				log.Println(c.mother.String(), "\n", geneticshuman.String(c.mother.genes), "\nand", c.father.String(), "\n", geneticshuman.String(c.father.genes), "\nhad a baby\n", geneticshuman.String(c.genes))

				// Add the child to the pool that will be added to the village.
				children = append(children, c)
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
	// Generate a new person.
	p := v.newPerson()
	p.lastName = v.lastGen.String()

	// Set the person's age and birthday.
	p.age = ageOfAdulthood + rand.Intn(20)
	p.bday = rand.Intn(365)

	// Set the person's genes.
	p.genes = genetics.NewRandom()
	p.fixGenes()

	// Add the person to the village.
	v.People = append(v.People, p)

	// Log the arrival.
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
