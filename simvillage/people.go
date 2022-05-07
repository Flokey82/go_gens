package simvillage

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

type PeopleManager struct {
	people   []*Person
	last_pop int
	log      []string
}

func NewPeopleManager() *PeopleManager {
	m := &PeopleManager{}
	m.people = nil

	m.init_population(STARTING_POP)

	// Give all existing people relationship objects
	m.init_pop_relationships()

	// Used to only sort when size changes
	m.last_pop = 0
	m.log = nil
	return m
}

func (m *PeopleManager) init_population(to_generate int) {
	if to_generate < 4 {
		to_generate = 4
	}
	for i := 0; i <= to_generate; i++ {
		if rand.Intn(1) < 1 {
			m.people = append(m.people, NewPerson("", rand.Intn(20)+10, 1))
			m.people = append(m.people, NewPerson("", rand.Intn(20)+10, 0))
		} else {
			m.people = append(m.people, NewPerson("", rand.Intn(20)+10, 0))
			m.people = append(m.people, NewPerson("", rand.Intn(20)+10, 1))
		}
	}
}
func (m *PeopleManager) init_pop_relationships() {
	for _, p := range m.people {
		// Every person has a relationship with every
		// other person
		p.init_relationships(m.people)
	}
}

func (m *PeopleManager) marriage_manager() {
	var single_men, single_women []*Person

	for _, p := range m.people {
		if p.gender == "Male" && p.romance {
			single_men = append(single_men, p)
		} else if p.gender == "Female" && p.romance {
			single_women = append(single_women, p)
		}
	}
}

func (m *PeopleManager) add_child(father, mother *Person) {
	mother.job = ""
	child := NewPerson("", 0, GenderRandom)
	child.lname = father.lname
}

func (m *PeopleManager) settlers() {
	if rand.Intn(360-SETTLER_CHANCE)+SETTLER_CHANCE == 360 {
		num_settlers := rand.Intn(10) + 10

		for i := 0; i <= num_settlers; i++ {
			m.people = append(m.people, NewPerson("", rand.Intn(56)+4, -1))
		}

		m.log = append(m.log, fmt.Sprintf("%d  new settlers have arrived!", num_settlers))
	}
}

func (m *PeopleManager) children() {
	if rand.Intn(19) > 18 {
		m.people = append(m.people, NewPerson("", rand.Intn(59)+1, -1))

		// Init persons relationships
		m.people[len(m.people)-1].relationships.init_relationships(m.people)

		// Update everyone but the persons relationship
		// objects to include the new person
		for _, p := range m.people[:len(m.people)-1] {
			p.relationships.add_relationship(m.people[len(m.people)-1], 0, "")
		}

		m.log = append(m.log, fmt.Sprintf("%s has appeared!", m.people[len(m.people)-1].name))
	}
	// Adding a child
	if rand.Intn(24) == 1 {

		// Create a new person object
		m.people = append(m.people, NewPerson("", 1, -1))

		// Init childs relationships
		m.people[len(m.people)-1].relationships.init_relationships(m.people)

		// Update everyone but the childs relationship
		// objects to include the new child
		for _, p := range m.people[:len(m.people)-1] {
			p.relationships.add_relationship(m.people[len(m.people)-1], 0, "")
		}
		// Log the event
		m.log = append(m.log, fmt.Sprintf("%s was born!", m.people[len(m.people)-1].name))
	}
}

func (m *PeopleManager) sort_people(sort_by string) {
	if sort_by == "job" {
		sort.Slice(m.people, func(i, j int) bool {
			return m.people[i].job < m.people[j].job
		})
	}
}

func (m *PeopleManager) Tick() []string {
	// only sort when needed
	if len(m.people) != m.last_pop {
		// Sorts the people
		m.sort_people("job")
		m.last_pop = len(m.people)
	}
	// Check if a child is randomly born
	m.children()

	// Tick if new settlers arrive
	m.settlers()

	// Tick people
	for _, p := range m.people {
		m.log = append(m.log, p.Tick()...)
	}
	// Get stats
	for _, p := range m.people {
		m.log = append(m.log, p.show_stats())
	}
	cp_log := m.log
	m.log = nil
	return cp_log
}

type GenderType int

const (
	GenderMale   GenderType = 0
	GenderFemale GenderType = 1
	GenderRandom GenderType = -1
)

type Person struct {
	id_num        int
	name          string
	lname         string
	can_work      bool
	job           string
	age           int
	gender        string
	hunger        int
	relationships *Relations
	pers          *Personality
	mood          *Mood
	romance       bool
	spouse        string
	children      []*Person
	log           []string
	history_log   []string
}

func NewPerson(job string, age int, gender GenderType) *Person {
	p := &Person{}
	// Unique identifier
	p.id_num = int(rand.Int31())

	// Assigned a name
	p.name = get_name()
	p.lname = get_last_name()

	// Assign job
	p.can_work = false
	if 0 < age && age < 5 {
		p.job = "Infant"
	}
	if 4 < age && age < 10 {
		p.job = "Child"
	} else {
		p.job = job
		p.can_work = true
	}
	p.age = age

	// Assign gender
	if gender == GenderMale {
		p.gender = "Male"
	} else if gender == GenderFemale {
		p.gender = "Female"
	} else if rand.Intn(1) < 1 {
		p.gender = "Male"
	} else {
		p.gender = "Female"
	}
	// Personality information
	p.pers = NewPersonality(p.name)
	p.mood = NewMood()
	p.romance = false
	p.spouse = ""

	// Relationship objects
	p.relationships = NewRelations(p.name, p.id_num)

	// Statuses
	/*
	   Hunger
	   Scale from 1-10
	   1 = Starving
	   2-5 = Hungry
	   6-8 = Sated
	   9-10 = Fed
	*/
	p.hunger = 5
	p.history_log = nil
	p.log = nil

	// Log about the person when they are made
	p.log = append(p.log, p.pers.get_backstory())
	return p
}

func (p *Person) Tick() []string {
	p.log = append(p.log, p.relationships.Tick()...)

	// Manage mood
	p.mood.Tick()

	// Manage romance
	p.check_marriage()
	// TODO
	// else: // Roll for spouse
	//  pass

	if 4 < p.age && p.age < 10 {
		p.job = "Child"
	}
	cp_log := p.log
	p.log = nil
	return cp_log
}

func (p *Person) init_relationships(people []*Person) {
	p.relationships.init_relationships(people)
}

// todo put in romance manager?
func (p *Person) check_marriage() {
	// Check for spouse
	if (p.romance == false) && (18 < p.age && p.age < 50) && (p.spouse == "") {
		// Now eligable to marry
		p.romance = true
		p.log = append(p.log, fmt.Sprintf("%s (%s) is looking for a partner.", p.name, p.gender))
	} else if p.romance == true {
		// Looking for a partner
		//pass
	}
}

func (p *Person) set_spouse(person *Person) {
	p.spouse = person.name
}

func (p *Person) show_stats() string {
	symb := "\u001b[34m♂\u001b[0m"
	if p.gender == "Female" {
		symb = "\u001b[35;1m♀\u001b[0m"
	}
	job := p.job
	if p.job == "" {
		job = "None"
	}
	if job == "Farmer" {
		job = "\u001b[34m" + job + "\u001b[0m"
	} else if job == "Woodcutter" {
		job = "\u001b[33m" + job + "\u001b[0m"
	} else if job == "Miner" {
		job = "\u001b[35m" + job + "\u001b[0m"
	} else if job == "Hunter" {
		job = "\u001b[32m" + job + "\u001b[0m"
	} else {
		job = "\u001b[36m" + job + "\u001b[0m"
	}

	// Basic information about a villager
	basic_info := fmt.Sprintf("%s  %s  Age:%d %s", p.name, job, p.age, symb)

	// Relationship info about a villager
	rel_info := p.relationships.get_rels_str()

	// Advanced information on a villager
	show_rel := "--"
	if p.romance == true && p.spouse == "" {
		show_rel = "Single"
	} else if p.romance == true && p.spouse != "" {
		show_rel = "Married"
	}

	// Show mood
	mood := p.mood.mood
	adv_info := fmt.Sprintf("%s | %s", mood, show_rel)

	return basic_info + adv_info + strings.Join(rel_info, ",")
}
