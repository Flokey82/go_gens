package simvillage

import (
	"fmt"
	"math"
	"math/rand"
)

type JobManager struct {
	job_ratios     map[string]float64
	people_manager *PeopleManager
	city_stats     *CityManager
	rng            *RandomEffects
	farmers        []*Person
	woodcutters    []*Person
	miners         []*Person
	hunters        []*Person
	mothers        []*Person
	unemployed     []*Person
	logs           []string
}

func NewJobManager(people_manager *PeopleManager, city_stats *CityManager) *JobManager {
	m := &JobManager{
		people_manager: people_manager,
		city_stats:     city_stats,
		job_ratios: map[string]float64{
			"Farmer":     0.30,
			"Woodcutter": 0.25,
			"Miner":      0.25,
			"Hunter":     0.2,
		},
	}

	m.rng = NewRandomEffects()
	m.farmers = nil
	m.woodcutters = nil
	m.miners = nil
	m.hunters = nil
	m.mothers = nil
	m.unemployed = nil
	m.logs = nil

	// Init workers jobs
	m.init_workers()
	return m
}

func (m *JobManager) Tick() []string {
	m.age_based_jobs()
	m.assign_workers()
	m.tick_jobs()
	cp_logs := m.logs
	m.logs = nil
	return cp_logs
}

func (m *JobManager) age_based_jobs() {
	for _, p := range m.people_manager.people {
		if (0 < p.age && p.age < 5) && (p.job != "Infant") {
			p.job = "Infant"
			p.can_work = false
		} else if (6 < p.age && p.age < 10) && (p.job != "Child") {
			p.job = "Child"
			p.can_work = false
		} else if (65 < p.age) && (p.job != "Old Person") {
			// remove from lists
			p.job = "Old Person"
			p.can_work = false
		}
	}
}

func (m *JobManager) update_unemployed() {
	// Get a list of unnasigned workers
	var unassigned []*Person
	for _, person := range m.people_manager.people {
		if (person.job == "") && (person.can_work) {
			unassigned = append(unassigned, person)
		}
	}
	m.unemployed = unassigned
}

// Call when first init village

func (m *JobManager) init_workers() {
	m.update_unemployed()

	def_jobs := []string{"Farmer", "Woodcutter", "Miner", "Hunter"}

	for i := range m.unemployed {
		chosen := def_jobs[i%len(def_jobs)]
		m.unemployed[i].job = chosen

		m.logs = append(m.logs, fmt.Sprintf("%s was chosen to be a %s.", m.unemployed[i].name, chosen))

		if chosen == "Farmer" {
			m.farmers = append(m.farmers, m.unemployed[i])
		} else if chosen == "Woodcutter" {
			m.woodcutters = append(m.woodcutters, m.unemployed[i])
		} else if chosen == "Miner" {
			m.miners = append(m.miners, m.unemployed[i])
		} else if chosen == "Hunter" {
			m.hunters = append(m.hunters, m.unemployed[i])
		}
	}
	m.update_unemployed()
}
func (m *JobManager) assign_workers() {
	m.update_unemployed()

	// Find jobs that need to be filled
	pop := len(m.people_manager.people)

	var needed_jobs []string
	// Check for farmers
	if (float64(len(m.farmers)) / float64(pop)) < m.job_ratios["Farmer"] {
		needed_jobs = append(needed_jobs, "Farmer")
	}
	if (float64(len(m.woodcutters)) / float64(pop)) < m.job_ratios["Woodcutter"] {
		needed_jobs = append(needed_jobs, "Woodcutter")
	}
	if (float64(len(m.miners)) / float64(pop)) < m.job_ratios["Miner"] {
		needed_jobs = append(needed_jobs, "Miner")
	}
	if (float64(len(m.hunters)) / float64(pop)) < m.job_ratios["Hunter"] {
		needed_jobs = append(needed_jobs, "Hunter")
	}

	if needed_jobs == nil {
		// Give default job
		// TODO: Better fix for no needed jobs being selected
		needed_jobs = append(needed_jobs, "Farmer")
	}
	// Assign workers to jobs that aren't as filled
	for _, worker := range m.unemployed {
		chosen := needed_jobs[rand.Intn(len(needed_jobs))]

		worker.job = chosen

		m.logs = append(m.logs, fmt.Sprintf("%s was chosen to be a %s.", worker.name, chosen))

		if chosen == "Farmer" {
			m.farmers = append(m.farmers, worker)
		} else if chosen == "Woodcutter" {
			m.woodcutters = append(m.woodcutters, worker)
		} else if chosen == "Miner" {
			m.miners = append(m.miners, worker)
		} else if chosen == "Hunter" {
			m.hunters = append(m.hunters, worker)
		}
	}
	m.update_unemployed()
}
func (m *JobManager) tick_jobs() {

	m.tick_farmers()
	m.tick_woodcutters()
	m.tick_miners()
	m.tick_hunters()
	m.tick_unemployed()
	m.tick_mothers()
}
func (m *JobManager) tick_farmers() {
	if len(m.farmers) == 0 {
		return
	}
	for _, n := range m.farmers {
		// base gathered
		base := 10.0

		// Get productivity level
		prod := n.mood.productivity

		// Sample by chance productivity
		chance := m.rng.get_mod()

		// Get final gathering quota
		final := math.Floor((base * prod) * chance)

		m.logs = append(m.logs, fmt.Sprintf("%s harvests %f crops.", n.name, final))

		m.city_stats.food += int(final)
	}
}
func (m *JobManager) tick_woodcutters() {
	if len(m.woodcutters) == 0 {
		return
	}
	for _, n := range m.woodcutters {
		// base gathered
		base := 20.0

		// Get productivity level
		prod := n.mood.productivity

		// Sample by chance productivity
		chance := m.rng.get_mod()

		// Get final gathering quota
		final := math.Floor((base * prod) * chance)
		m.logs = append(m.logs, fmt.Sprintf("%s chops %f wood.", n.name, final))
		m.city_stats.wood += int(final)
	}
}
func (m *JobManager) tick_miners() {

	if len(m.miners) == 0 {
		return
	}
	for _, n := range m.miners {
		// base gathered
		base := 15.0

		// Get productivity level
		prod := n.mood.productivity

		// Sample by chance productivity
		chance := m.rng.get_mod()

		// Get final gathering quota
		final := math.Floor((base * prod) * chance)
		m.logs = append(m.logs, fmt.Sprintf("%s mines %f stone.", n.name, final))
		m.city_stats.stone += int(final)
	}
}
func (m *JobManager) tick_hunters() {
	if len(m.hunters) == 0 {
		return
	}

	for _, n := range m.hunters {
		// base gathered
		base := 20.0

		// Get productivity level
		prod := n.mood.productivity

		// Sample by chance productivity
		chance := m.rng.get_mod()

		// Get final gathering quota
		final := math.Floor((base * prod) * chance)

		// Hunters are high risk high reward
		if rand.Intn(2) > 0 {
			m.logs = append(m.logs, fmt.Sprintf("%s hunts and brings back %f food", n.name, final))
			m.city_stats.food += int(final)
		} else {
			m.logs = append(m.logs, fmt.Sprintf("%s hunts and catches nothing.", n.name))
		}
	}
}
func (m *JobManager) tick_unemployed() {
}
func (m *JobManager) tick_mothers() {
	for _, n := range m.mothers {
		youngest := 99
		for _, c := range n.children {
			if c.age < youngest {
				youngest = c.age
			}
		}
		if youngest > 4 {
			n.job = "Unemployed"
		}
	}
}

type Job struct {
	name string
	mode string
}

func NewJob(name, mode string) *Job {
	j := &Job{}
	// name: name of the job
	// mode: what kind of job it is
	j.name = name
	j.mode = mode
	return j
}
func (j *Job) Tick() []string {
	return nil
}
