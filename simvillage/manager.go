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
	cal          *Calendar
	log          *HistoryManager
	villagers    *PeopleManager
	prof         *JobManager
	stats        *CityManager
	reaper       *Death
	socialEvents *SocialEvents
	markov       *MarkovGen
	managers     []Mngi
}

func NewInstance() *Instance {
	in := &Instance{}
	in.cal = NewCalendar()                          // Controls the incrementing of date
	in.log = NewHistoryManager(in.cal)              // Controls logging events
	in.markov = NewMarkovGen()                      // Better flavor text generation
	in.villagers = NewPeopleManager()               // Create some people
	in.stats = NewCityManager(in.villagers)         // Controls city metrics and stockpiles
	in.prof = NewJobManager(in.villagers, in.stats) // Controls who has what job
	in.socialEvents = NewSocialEvents(in.villagers) // Make social events happen
	in.reaper = NewDeath(in.villagers, in.markov)   // Controls death and dying events

	in.managers = []Mngi{ // Holds a list of all known objects that need to tick
		in.log,          // Create a new logging entry for the day
		in.cal,          // Date goes up
		in.stats,        // People eat and food is lost
		in.villagers,    // Villagers tick
		in.prof,         // Professions are managed, jobs reassigned
		in.socialEvents, // Random social events occur
		in.reaper,       // Sometimes people kick the bucket
	}
	return in
}

func (in *Instance) tickMonth() {
	for i := 1; i <= 10; i++ {
		in.nextTick()
	}
	in.log.listTodaysEvents()
}

func (in *Instance) TickDay() {
	in.nextTick()
	in.log.listTodaysEvents()
}

func (in *Instance) nextTick() {
	for _, manager := range in.managers {
		// Tick current manager and record output
		in.log.addEvents(manager.Tick())
	}
	time.Sleep(500000000)
	os.Stdout.Write([]byte{0x1B, 0x5B, 0x33, 0x3B, 0x4A, 0x1B, 0x5B, 0x48, 0x1B, 0x5B, 0x32, 0x4A})
}
