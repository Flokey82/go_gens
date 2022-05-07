// Package simvillage is a port of the wonderful village simulator by Kontari.
// See: https://github.com/Kontari/Village/
package simvillage

import (
	"os"
	"time"
)

type Mngi interface {
	Tick() []string
}
type Instance struct {
	cal           *Calendar
	log           *HistoryManager
	villagers     *PeopleManager
	prof          *JobManager
	stats         *CityManager
	reaper        *Death
	social_events *SocialEvents
	markov        *MarkovGen
	managers      []Mngi
}

func NewInstance() *Instance {
	in := &Instance{}
	// Controls the incrementing of date
	in.cal = NewCalendar()
	// Controls logging events
	in.log = NewHistoryManager(in.cal)
	// Better flavor text generation
	in.markov = NewMarkovGen()
	// Create some people
	in.villagers = NewPeopleManager()
	// Controls city metrics and stockpiles
	in.stats = NewCityManager(in.villagers)
	// Controls who has what job
	in.prof = NewJobManager(in.villagers, in.stats)
	// Make social events happen
	in.social_events = NewSocialEvents(in.villagers)
	// Controls death and dying events
	in.reaper = NewDeath(in.villagers, in.markov)
	// Holds a list of all known objects that need to tick
	in.managers = nil
	// Create a new logging entry for the day
	in.managers = append(in.managers, in.log)
	in.managers = append(in.managers, in.cal)       // Date goes up
	in.managers = append(in.managers, in.stats)     // People eat and food is lost
	in.managers = append(in.managers, in.villagers) // Villagers tick
	// Professions are managed, jobs reassigned
	in.managers = append(in.managers, in.prof)
	in.managers = append(in.managers, in.social_events) // Random social events occur
	in.managers = append(in.managers, in.reaper)        // Sometimes people kick the bucket
	return in
}

func (in *Instance) tick_month() {
	for i := 1; i <= 10; i++ {
		in.next_tick()
	}
	in.log.list_todays_events()
}

func (in *Instance) TickDay() {
	in.next_tick()
	in.log.list_todays_events()
}

func (in *Instance) next_tick() {
	for _, manager := range in.managers {
		// print('Calling: {}'.format(manager))
		// Tick current manager and record output
		in.log.add_events(manager.Tick())
		// print(manager)
	}
	time.Sleep(500000000)
	os.Stdout.Write([]byte{0x1B, 0x5B, 0x33, 0x3B, 0x4A, 0x1B, 0x5B, 0x48, 0x1B, 0x5B, 0x32, 0x4A})
}
