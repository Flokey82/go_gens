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
	my_name string
	my_id   int
	rel     []*Relationship
	log     []string
}

func NewRelations(my_name string, my_id int) *Relations {
	r := &Relations{}
	r.my_name = my_name
	r.my_id = my_id

	// Holds person obj, rel value, text to describe it
	r.rel = nil
	r.log = nil
	return r
}

func (r *Relations) Tick() []string {
	cp_log := r.log
	r.log = nil
	return cp_log
}

func (r *Relations) get_rels_str() []string {
	var rels []string
	for _, r := range r.rel {
		rels = append(rels, fmt.Sprintf("%.2f ", r.Value))
	}
	return rels
}

func (r *Relations) init_relationships(people []*Person) {
	for _, p := range people {
		r.add_relationship(p, 2.0, "")
	}
}

func (r *Relations) add_relationship(person *Person, rel_value float64, rel_text string) {
	if rel_value == 0 {
		rel_value = 2.0
	}
	r.rel = append(r.rel, &Relationship{person, rel_value, rel_text})
}

func (r *Relations) del_relationship(tgt_person *Person) float64 {
	// When a villager dies remove the relationship.
	// Return strength of relationship and relationship
	// text to create a mood event.
	dead_rel_value := 0.0
	for _, person := range r.rel {
		if person.A == tgt_person {
			dead_rel_value = person.Value
			break
		}
	}
	return dead_rel_value
}

func (r *Relations) mod_relationship(value float64, person *Person) {
	for _, people := range r.rel {
		if people.A == person {
			// Get old relationship text
			old_rel_text := r.get_rel_text(people.Value)

			// Update new relationship data
			fin_rel_value := people.Value * value
			people.Value = fin_rel_value

			// Get new relationship text
			rel_text := r.get_rel_text(fin_rel_value)

			// If the relationship catagorychanged
			if old_rel_text != rel_text {
				if rel_text == RelLiked {
					rel_text = "\u001b[32m" + rel_text + "\u001b[0m"
				} else if rel_text == RelDisliked {
					rel_text = "\u001b[31m" + rel_text + "\u001b[0m"
				} else if rel_text == RelFriendly {
					rel_text = " \u001b[32;1m" + rel_text + "\u001b[0m"
				}
				new_rel_text := fmt.Sprintf("%s is now %s with %s. (%.2f)", r.my_name, rel_text, people.A.name, fin_rel_value)
				r.log = append(r.log, new_rel_text)
			}
		}
	}
}

func (r *Relations) get_relationship(person *Person) float64 {
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

func (r *Relations) get_rel_text(rel_value float64) string {
	if rel_value < 1.00 {
		return RelDisliked
	} else if 1.00 < rel_value && rel_value < 2.00 {
		return RelNeutral
	} else if 2.00 < rel_value && rel_value < 3.00 {
		return RelLiked
	} else if 3.00 < rel_value {
		return RelFriendly
	}
	return "?"
}
