package simvillage

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

type PeopleManager struct {
	people  []*Person
	lastPop int
	log     []string
}

func NewPeopleManager() *PeopleManager {
	m := &PeopleManager{
		people: nil,
	}
	m.initPopulation(STARTING_POP)

	// Give all existing people relationship objects
	m.initPopRelationships()

	// Used to only sort when size changes
	m.lastPop = 0
	return m
}

func (m *PeopleManager) initPopulation(toGenerate int) {
	if toGenerate < 4 {
		toGenerate = 4
	}
	for i := 0; i <= toGenerate; i++ {
		if rand.Intn(2) < 1 {
			m.people = append(m.people, NewPerson("", rand.Intn(20)+14, 1))
			m.people = append(m.people, NewPerson("", rand.Intn(20)+14, 0))
		} else {
			m.people = append(m.people, NewPerson("", rand.Intn(20)+14, 0))
			m.people = append(m.people, NewPerson("", rand.Intn(20)+14, 1))
		}
	}
}

func (m *PeopleManager) initPopRelationships() {
	for _, p := range m.people {
		// Every person has a relationship with every other person
		p.initRelationships(m.people)
	}
}

func (m *PeopleManager) marriageManager() {
	var singleMen, singleWomen []*Person
	for _, p := range m.people {
		if p.gender == GenderStrMale && p.romance {
			singleMen = append(singleMen, p)
		} else if p.gender == GenderStrFemale && p.romance {
			singleWomen = append(singleWomen, p)
		}
	}
}

func (m *PeopleManager) addChild(father, mother *Person) {
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
		m.people[len(m.people)-1].relationships.initRelationships(m.people)

		// Update everyone but the persons relationship
		// objects to include the new person
		for _, p := range m.people[:len(m.people)-1] {
			p.relationships.addRelationship(m.people[len(m.people)-1], 0, "")
		}
		m.log = append(m.log, fmt.Sprintf("%s has appeared!", m.people[len(m.people)-1].name))
	}
	// Adding a child
	if rand.Intn(24) == 1 {
		// Create a new person object
		m.people = append(m.people, NewPerson("", 1, -1))

		// Init childs relationships
		m.people[len(m.people)-1].relationships.initRelationships(m.people)

		// Update everyone but the childs relationship
		// objects to include the new child
		for _, p := range m.people[:len(m.people)-1] {
			p.relationships.addRelationship(m.people[len(m.people)-1], 0, "")
		}
		// Log the event
		m.log = append(m.log, fmt.Sprintf("%s was born!", m.people[len(m.people)-1].name))
	}
}

func (m *PeopleManager) sortPeople(sortBy string) {
	if sortBy == "job" {
		sort.Slice(m.people, func(i, j int) bool {
			return m.people[i].job < m.people[j].job
		})
	}
}

func (m *PeopleManager) Tick() []string {
	// only sort when needed
	if len(m.people) != m.lastPop {
		// Sorts the people
		m.sortPeople("job")
		m.lastPop = len(m.people)
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
		m.log = append(m.log, p.showStats())
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

const (
	GenderStrMale   = "Male"
	GenderStrFemale = "Female"
)

type Person struct {
	idNum         int
	name          string
	lname         string
	canWork       bool
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
	historyLog    []string
}

func NewPerson(job string, age int, gender GenderType) *Person {
	p := &Person{
		idNum: int(rand.Int31()), // Unique identifier
		name:  getName(),         // Assigned a name
		lname: getLastName(),     // Assigned a name
	}

	// Assign job
	p.canWork = false
	if 0 < age && age < 5 {
		p.job = JobInfant.name
	}
	if 4 < age && age < 10 {
		p.job = JobChild.name
	} else {
		p.job = job
		p.canWork = true
	}
	p.age = age

	// Assign gender
	if gender == GenderMale {
		p.gender = GenderStrMale
	} else if gender == GenderFemale {
		p.gender = GenderStrFemale
	} else if rand.Intn(2) < 1 {
		p.gender = GenderStrMale
	} else {
		p.gender = GenderStrFemale
	}

	// Personality information
	p.pers = NewPersonality(p.name)
	p.mood = NewMood()
	p.romance = false
	p.spouse = ""

	// Relationship objects
	p.relationships = NewRelations(p.name, p.idNum)

	// Statuses
	// Hunger
	// Scale from 1-10
	// 1 = Starving
	// 2-5 = Hungry
	// 6-8 = Sated
	// 9-10 = Fed
	p.hunger = 5
	p.historyLog = nil
	p.log = nil

	// Log about the person when they are made
	p.log = append(p.log, p.pers.getBackstory())
	return p
}

func (p *Person) Tick() []string {
	p.log = append(p.log, p.relationships.Tick()...)

	// Manage mood
	p.mood.Tick()

	// Manage romance
	p.checkMarriage()
	// TODO
	// else: // Roll for spouse
	//  pass

	if 4 < p.age && p.age < 10 {
		p.job = JobChild.name
	}
	cp_log := p.log
	p.log = nil
	return cp_log
}

func (p *Person) initRelationships(people []*Person) {
	p.relationships.initRelationships(people)
}

// todo put in romance manager?
func (p *Person) checkMarriage() {
	// Check for spouse
	if !p.romance && 18 < p.age && p.age < 50 && p.spouse == "" {
		p.romance = true // Now eligable to marry
		p.log = append(p.log, fmt.Sprintf("%s (%s) is looking for a partner.", p.name, p.gender))
	} else if p.romance {
		// Looking for a partner
		// TODO: Find a partner that we like.
	}
}

func (p *Person) setSpouse(person *Person) {
	p.spouse = person.name
}

func (p *Person) showStats() string {
	symb := "\u001b[34m♂\u001b[0m"
	if p.gender == GenderStrFemale {
		symb = "\u001b[35;1m♀\u001b[0m"
	}

	job := p.job
	if job == "" {
		job = "None"
	} else if job == JobFarmer.name {
		job = "\u001b[34m" + job + "\u001b[0m"
	} else if job == JobWoodcutter.name {
		job = "\u001b[33m" + job + "\u001b[0m"
	} else if job == JobMiner.name {
		job = "\u001b[35m" + job + "\u001b[0m"
	} else if job == JobHunter.name {
		job = "\u001b[32m" + job + "\u001b[0m"
	} else {
		job = "\u001b[36m" + job + "\u001b[0m"
	}

	// Basic information about a villager
	basicInfo := fmt.Sprintf("%8s  %20s  Age:%3d %s", p.name, job, p.age, symb)

	// Relationship info about a villager
	relInfo := p.relationships.getRelsStr()

	// Advanced information on a villager and mood.
	showRel := "--"
	if p.romance && p.spouse == "" {
		showRel = "Single"
	} else if p.romance && p.spouse != "" {
		showRel = "Married"
	}
	advInfo := fmt.Sprintf("[%12s | %8s]", p.mood.mood, showRel)
	return basicInfo + advInfo + strings.Join(relInfo, ",")
}
