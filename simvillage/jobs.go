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
	employed       map[string][]*Person
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
		employed: make(map[string][]*Person),
	}

	m.rng = NewRandomEffects()
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
		if (0 < p.age && p.age < 5) && (p.job != JobInfant.name) {
			p.job = JobInfant.name
			p.can_work = false
		} else if (6 < p.age && p.age < 10) && (p.job != JobChild.name) {
			p.job = JobChild.name
			p.can_work = false
		} else if (65 < p.age) && (p.job != JobOldPerson.name) {
			// remove from lists
			p.job = JobOldPerson.name
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
		m.employed[chosen] = append(m.employed[chosen], m.unemployed[i])
	}
	m.update_unemployed()
}

func (m *JobManager) assign_workers() {
	m.update_unemployed()

	// Find jobs that need to be filled
	pop := len(m.people_manager.people)

	// Clean up dead people.
	m.employed = make(map[string][]*Person)
	for _, w := range m.people_manager.people {
		m.employed[w.job] = append(m.employed[w.job], w)
	}

	var needed_jobs []string
	for _, j := range def_jobs {
		if (float64(len(m.employed[j.name])) / float64(pop)) < m.job_ratios[j.name] {
			needed_jobs = append(needed_jobs, j.name)
		}
	}

	// Give default job
	if needed_jobs == nil {
		// TODO: Better fix for no needed jobs being selected
		needed_jobs = append(needed_jobs, JobFarmer.name)
	}

	// Assign workers to jobs that aren't as filled
	for _, worker := range m.unemployed {
		chosen := needed_jobs[rand.Intn(len(needed_jobs))]
		worker.job = chosen

		m.logs = append(m.logs, fmt.Sprintf("%s was chosen to be a %s.", worker.name, chosen))
		m.employed[chosen] = append(m.employed[chosen], worker)
	}
	m.update_unemployed()
}

func (m *JobManager) tick_jobs() {
	for _, j := range def_jobs {
		for _, n := range m.employed[j.name] {
			base := j.Tick() // base gathered
			if base > 0.0 {
				// Get productivity level
				prod := n.mood.productivity

				// Sample by chance productivity
				chance := m.rng.get_mod()

				// Get final gathering quota
				final := math.Floor((base * prod) * chance)
				switch j.produces {
				case ResGame, ResCrops:
					m.city_stats.food += int(final)
				case ResStone:
					m.city_stats.stone += int(final)
				case ResWood:
					m.city_stats.wood += int(final)
				}
				m.logs = append(m.logs, fmt.Sprintf(j.successMsg, n.name, final))
			} else {
				m.logs = append(m.logs, fmt.Sprintf(j.failMsg, n.name))
			}
		}
	}
	// TODO: Remove dead.
	m.tick_unemployed()
	m.tick_mothers()
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

type Resource int

const (
	ResNothing Resource = iota
	ResStone
	ResCrops
	ResWood
	ResGame
)

type Job struct {
	name        string // name of the job
	mode        string // what kind of job it is
	successMsg  string
	failMsg     string
	produces    Resource
	produceBase float64
	chance      int
	canFail     bool
}

func NewJob(name, mode string) *Job {
	return &Job{
		name: name,
		mode: mode,
	}
}

func (j *Job) Tick() float64 {
	if j.canFail && rand.Intn(j.chance) == 0 {
		return 0.0
	}
	// base gathered
	return j.produceBase
}

var (
	JobMiner = &Job{
		name:        "Miner",
		successMsg:  "%s mines %.0f stone.",
		produces:    ResStone,
		produceBase: 15.0,
		chance:      0,
	}
	JobFarmer = &Job{
		name:        "Farmer",
		successMsg:  "%s harvests %.0f crops.",
		produces:    ResCrops,
		produceBase: 10.0,
		chance:      0,
	}
	JobWoodcutter = &Job{
		name:        "Woodcutter",
		successMsg:  "%s chops %.0f wood.",
		produces:    ResWood,
		produceBase: 20.0,
		chance:      0,
	}
	JobHunter = &Job{
		name:        "Hunter",
		successMsg:  "%s hunts and brings back %.0f food",
		failMsg:     "%s hunts and catches nothing.",
		produces:    ResGame,
		produceBase: 20.0,
		canFail:     true,
		chance:      2,
	}
	JobInfant = &Job{
		name:        "Infant",
		successMsg:  "%s poops %.0f times",
		failMsg:     "%s didn't poop",
		produces:    ResNothing,
		produceBase: 2,
		canFail:     true,
		chance:      2,
	}
	JobChild = &Job{
		name:        "Child",
		successMsg:  "%s groans %.0f times",
		failMsg:     "%s was helpful in the household",
		produces:    ResNothing,
		canFail:     true,
		produceBase: 1,
		chance:      2,
	}
	JobOldPerson = &Job{
		name:        "Old Person",
		successMsg:  "%s complains %.0f times",
		produces:    ResNothing,
		canFail:     false,
		produceBase: 1,
		chance:      2,
	}
)
var def_jobs = []*Job{JobFarmer, JobWoodcutter, JobMiner, JobHunter}
