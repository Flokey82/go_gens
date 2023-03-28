package simvillage

import (
	"fmt"
	"math"
	"math/rand"
)

type SocialEvents struct {
	people          *PeopleManager
	places          []string
	verbs           []string
	negative_verbs  []string
	friendly_events []string
	neutral_events  []string
	disliked_events []string
	log             []string
}

func NewSocialEvents(peopleObjs *PeopleManager) *SocialEvents {
	return &SocialEvents{
		people: peopleObjs,
		places: []string{
			"street",
			"store",
			"well",
			"bar",
			"bakery",
			"butchery",
			"town square",
			"forest",
			"gardens",
		},
		verbs: []string{
			"runs into",
			"meets",
			"talks with",
			"hangs with",
			"spots",
		},
		negative_verbs: []string{
			"spits on",
			"gets into a fight with",
			"attacks",
			"insults",
		},
		friendly_events: []string{"%s %s %s at the %s"},
		neutral_events:  []string{"%s %s %s at the %s"},
		disliked_events: []string{"%s %s %s at the %s"},
	}
}

func (s *SocialEvents) Tick() []string {
	prctTick := SOCIAL_CHANCE
	loops := int(math.Floor(prctTick * float64(len(s.people.people))))
	for i := 0; i < loops; i++ {
		s.randomEvent()
	}
	cpLog := s.log
	s.log = nil
	return cpLog
}

func randPerson(pp []*Person) *Person {
	return pp[rand.Intn(len(pp))]
}

func (s *SocialEvents) randomEvent() {
	// Add timeouts here
	// Select a random villager to have an event happen
	selectedPerson := randPerson(s.people.people)
	for selectedPerson.age < 2 {
		selectedPerson = randPerson(s.people.people) // TODO: Hangs here!
	}

	anotherPerson := randPerson(s.people.people)
	for (anotherPerson == selectedPerson) && (anotherPerson.age < 2) {
		anotherPerson = randPerson(s.people.people)
	}

	// Now we have two people to trigger an event with
	selPersonRel := selectedPerson.relationships.getRelationship(anotherPerson)
	anoPersonRel := anotherPerson.relationships.getRelationship(selectedPerson)
	if selPersonRel == 0 || anoPersonRel == 0 {
		return
	}

	// See if their relationship is bad, neutral, or good
	sumRelationship := selPersonRel + anoPersonRel
	if sumRelationship < 1 {
		// Disliked
		s.negativeEvent(selectedPerson, anotherPerson)
	} else if 1 < sumRelationship && sumRelationship < 3 {
		// Neutral
		// See if the neutral event will be positive or negative
		if rand.Float64() < FRIENDLY_CHANCE {
			s.negativeEvent(selectedPerson, anotherPerson)
		} else {
			s.positiveEvent(selectedPerson, anotherPerson)
		}
	} else if 3 < sumRelationship {
		// Positive
		s.positiveEvent(selectedPerson, anotherPerson)
	}
}

func (s *SocialEvents) negativeEvent(p1, p2 *Person) {
	eventText := fmt.Sprintf(pickRandString(s.disliked_events), p1.name, pickRandString(s.negative_verbs), p2.name, pickRandString(s.places))
	p1.relationships.modRelationship(0.9, p2)
	p2.relationships.modRelationship(0.9, p1)
	s.log = append(s.log, eventText)
}

func (s *SocialEvents) positiveEvent(p1, p2 *Person) {
	eventText := fmt.Sprintf(pickRandString(s.friendly_events), p1.name, pickRandString(s.verbs), p2.name, pickRandString(s.places))
	p1.relationships.modRelationship(1.1, p2)
	p2.relationships.modRelationship(1.1, p1)
	s.log = append(s.log, eventText)
}
