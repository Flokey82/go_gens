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

func NewSocialEvents(people_objs *PeopleManager) *SocialEvents {
	return &SocialEvents{
		people: people_objs,
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
	prct_tick := SOCIAL_CHANCE
	loops := int(math.Floor(prct_tick * float64(len(s.people.people))))
	for i := 0; i < loops; i++ {
		s.random_event()
	}
	cp_log := s.log
	s.log = nil
	return cp_log
}

func randPerson(pp []*Person) *Person {
	return pp[rand.Intn(len(pp))]
}

func (s *SocialEvents) random_event() {
	// Add timeouts here
	// Select a random villager to have an event happen
	selected_person := randPerson(s.people.people)
	for selected_person.age < 2 {
		selected_person = randPerson(s.people.people) // TODO: Hangs here!
	}

	another_person := randPerson(s.people.people)
	for (another_person == selected_person) && (another_person.age < 2) {
		another_person = randPerson(s.people.people)
	}

	// Now we have two people to trigger an event with
	sel_p_rel := selected_person.relationships.get_relationship(another_person)
	ano_p_rel := another_person.relationships.get_relationship(selected_person)
	if sel_p_rel == 0 || ano_p_rel == 0 {
		return
	}

	// See if their relationship is bad, neutral, or good
	sum_relationship := sel_p_rel + ano_p_rel
	if sum_relationship < 1 {
		// Disliked
		s.negative_event(selected_person, another_person)
	} else if 1 < sum_relationship && sum_relationship < 3 {
		// Neutral
		// See if the neutral event will be positive or negative
		if rand.Float64() < FRIENDLY_CHANCE {
			s.negative_event(selected_person, another_person)
		} else {
			s.positive_event(selected_person, another_person)
		}
	} else if 3 < sum_relationship {
		// Positive
		s.positive_event(selected_person, another_person)
	}
}

func (s *SocialEvents) negative_event(p_one, p_two *Person) {
	event_text := fmt.Sprintf(pickRandString(s.disliked_events), p_one.name, pickRandString(s.negative_verbs), p_two.name, pickRandString(s.places))
	p_one.relationships.mod_relationship(0.9, p_two)
	p_two.relationships.mod_relationship(0.9, p_one)
	s.log = append(s.log, event_text)
}

func (s *SocialEvents) positive_event(p_one, p_two *Person) {
	event_text := fmt.Sprintf(pickRandString(s.friendly_events), p_one.name, pickRandString(s.verbs), p_two.name, pickRandString(s.places))
	p_one.relationships.mod_relationship(1.1, p_two)
	p_two.relationships.mod_relationship(1.1, p_one)
	s.log = append(s.log, event_text)
}
