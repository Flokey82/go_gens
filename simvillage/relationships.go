package simvillage

import (
	"fmt"
)

type Relationship struct {
	A     *Person
	Value float64
	Text  string
}

type Relations struct {
	myName string
	myID   int
	rel    []*Relationship
	log    []string
}

func NewRelations(myName string, myID int) *Relations {
	return &Relations{
		myName: myName,
		myID:   myID,
	} // Holds person obj, rel value, text to describe it
}

func (r *Relations) Tick() []string {
	cpLog := r.log
	r.log = nil
	return cpLog
}

func (r *Relations) getRelsStr() []string {
	var rels []string
	for _, r := range r.rel {
		rels = append(rels, fmt.Sprintf("%.2f ", r.Value))
	}
	return rels
}

func (r *Relations) initRelationships(people []*Person) {
	for _, p := range people {
		r.addRelationship(p, 2.0, "")
	}
}

func (r *Relations) addRelationship(person *Person, relValue float64, relText string) {
	if relValue == 0.0 {
		relValue = 2.0
	}
	r.rel = append(r.rel, &Relationship{person, relValue, relText})
}

func (r *Relations) delRelationship(tgtPerson *Person) float64 {
	// When a villager dies remove the relationship.
	// Return strength of relationship and relationship
	// text to create a mood event.
	var deadRelValue float64
	for _, person := range r.rel {
		if person.A == tgtPerson {
			deadRelValue = person.Value
			break
		}
	}
	return deadRelValue
}

func (r *Relations) modRelationship(value float64, person *Person) {
	for _, people := range r.rel {
		if people.A == person {
			// Get old relationship text
			oldRelText := r.getRelText(people.Value)

			// Update new relationship data
			finRelValue := people.Value * value
			people.Value = finRelValue

			// Get new relationship text
			relText := r.getRelText(finRelValue)

			// If the relationship catagorychanged
			if oldRelText != relText {
				if relText == RelLiked {
					relText = "\u001b[32m" + relText + "\u001b[0m"
				} else if relText == RelDisliked {
					relText = "\u001b[31m" + relText + "\u001b[0m"
				} else if relText == RelFriendly {
					relText = " \u001b[32;1m" + relText + "\u001b[0m"
				}
				newRelText := fmt.Sprintf("%s is now %s with %s. (%.2f)", r.myName, relText, people.A.name, finRelValue)
				r.log = append(r.log, newRelText)
			}
		}
	}
}

func (r *Relations) getRelationship(person *Person) float64 {
	for _, people := range r.rel {
		if people.A == person {
			return people.Value
		}
	}

	// Person must not exist so add
	r.rel = append(r.rel, &Relationship{person, 2.0, ""})
	return 2.0
}

const (
	RelDisliked = "Disliked"
	RelNeutral  = "Neutral"
	RelLiked    = "Liked"
	RelFriendly = "Friendly"
)

func (r *Relations) getRelText(relValue float64) string {
	if relValue < 1.00 {
		return RelDisliked
	}
	if 1.00 < relValue && relValue < 2.00 {
		return RelNeutral
	}
	if 2.00 < relValue && relValue < 3.00 {
		return RelLiked
	}
	if 3.00 < relValue {
		return RelFriendly
	}
	return "?"
}
