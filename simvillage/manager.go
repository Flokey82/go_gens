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

	in.cal = NewCalendar()                           // Controls the incrementing of date
	in.log = NewHistoryManager(in.cal)               // Controls logging events
	in.markov = NewMarkovGen()                       // Better flavor text generation
	in.villagers = NewPeopleManager()                // Create some people
	in.stats = NewCityManager(in.villagers)          // Controls city metrics and stockpiles
	in.prof = NewJobManager(in.villagers, in.stats)  // Controls who has what job
	in.social_events = NewSocialEvents(in.villagers) // Make social events happen
	in.reaper = NewDeath(in.villagers, in.markov)    // Controls death and dying events

	in.managers = nil                                   // Holds a list of all known objects that need to tick
	in.managers = append(in.managers, in.log)           // Create a new logging entry for the day
	in.managers = append(in.managers, in.cal)           // Date goes up
	in.managers = append(in.managers, in.stats)         // People eat and food is lost
	in.managers = append(in.managers, in.villagers)     // Villagers tick
	in.managers = append(in.managers, in.prof)          // Professions are managed, jobs reassigned
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
		// Tick current manager and record output
		in.log.add_events(manager.Tick())
	}
	time.Sleep(500000000)
	os.Stdout.Write([]byte{0x1B, 0x5B, 0x33, 0x3B, 0x4A, 0x1B, 0x5B, 0x48, 0x1B, 0x5B, 0x32, 0x4A})
}
