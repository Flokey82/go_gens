package simvillage_simple

import (
	"fmt"
	"math/rand"

	"github.com/Flokey82/genetics"
	"github.com/Flokey82/genetics/geneticshuman"
)

const ageOfAdulthood = 16

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
	id           int            // Unique ID.
	firstName    string         // First name.
	lastName     string         // Last name.
	age          int            // Age in years.
	bday         int            // Day of birth.
	dead         bool           // True if dead.
	mother       *Person        // Mother.
	father       *Person        // Father.
	spouse       *Person        // TODO: keep track of spouses that might have perished?
	children     []*Person      // Children.
	gender       Gender         // Gender.
	pregnant     int            // Number of days pregnant.
	pregnantWith *Person        // Baby that will be born.
	genes        genetics.Genes // Genes.
}

// newPerson creates a new person.
func (v *Village) newPerson() *Person {
	gender := randGender()
	return &Person{
		id:        v.getNextID(),
		gender:    gender,
		firstName: v.firstGen[gender].String(),
	}
}

// Name returns the full name of the person.
func (p *Person) Name() string {
	return p.firstName + " " + p.lastName
}

// String returns the string representation of the person.
func (p *Person) String() string {
	var deadStr string
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
	return p.age > ageOfAdulthood && p.spouse == nil // Old enough and single.
}

// canBePregnant returns true if the person is old enough and not pregnant.
func (p *Person) canBePregnant() bool {
	// Female, has a spouse (implies old enough), and is currently not pregnant.
	// TODO: Set randomized upper age limit.
	return p.gender == GenderFemale && p.spouse != nil && p.pregnantWith == nil
}

// fixGenes makes sure that the gender is set properly in the genes.
// NOTE: This needs to be done due to the genetics package being a bit weird.
func (p *Person) fixGenes() {
	switch p.gender {
	case GenderFemale:
		geneticshuman.SetGender(&p.genes, geneticshuman.GenderFemale)
	case GenderMale:
		geneticshuman.SetGender(&p.genes, geneticshuman.GenderMale)
	default:
		geneticshuman.SetGender(&p.genes, 0)
	}
}

// isRelated returns true if a and b are related (up to first degree).
func isRelated(a, b *Person) bool {
	// Check if there is a parent/child relationship.
	if a == b.father || a == b.mother || b == a.father || b == a.mother {
		return true
	}

	// If either (or both) of the parents are nil, we assume that they are not related.
	if (a.father == nil && a.mother == nil) || (b.father == nil && b.mother == nil) {
		return false
	}

	// Check if there is a (half-) sibling relationship.
	return a.mother == b.mother || a.father == b.father
}
