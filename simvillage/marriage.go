package simvillage

import "fmt"

// Subclass of the people manager
// 1. Decides who is eligable for Marriage
// 2. Allows romantic events to happen in the
//    social manager

type Marriage struct {
	people []*Person
}

func NewMarriage(people []*Person) *Marriage {
	return &Marriage{
		people: people,
	}
}

func (m *Marriage) check_marriage(p *Person) {
	// Check for spouse
	if (p.romance == false) && (18 < p.age && p.age < 50) && (p.spouse == "") {
		// Now eligable to marry
		p.romance = true
		p.log = append(p.log, fmt.Sprintf("%s (%s) is looking for a partner.", p.name, p.gender))
	}

}
