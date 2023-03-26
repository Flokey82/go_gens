package simvillage

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/Flokey82/go_gens/gameconstants"
)

// Death controls chance of death, and death events
type Death struct {
	pop         *PeopleManager
	mark        *MarkovGen
	dead        []*Person
	ways_to_die []string
	old_age     []string
	suicides    []string
	log         []string
}

func NewDeath(pmanager *PeopleManager, mark *MarkovGen) *Death {
	return &Death{
		pop:  pmanager,
		mark: mark,
		ways_to_die: []string{ // By chance ways to die
			"%s dies in their sleep",
			"%s is lost in the night",
			"%s drowned",
		},
		old_age: []string{ // Aging related deaths
			"%s had a heart attack",
		},
		suicides: []string{ // Self-inflicted deaths
			"%s hangs themself",
			"%s jumps from a tree",
		}, // TODO: Work related deaths
	}
}

func (d *Death) Tick() []string {
	for _, p := range d.pop.people {
		d.tick_death(p)
	}
	cp_log := d.log
	d.log = nil
	return cp_log
}

// tick_death check if the villager will die today by aging
func (d *Death) tick_death(v *Person) {
	// Check if the villager dies of natural causes.
	if gameconstants.DiesAtAge(v.age) {
		d.kill_villager(v, "")
	}

	// Check if the villager is depressed
	if v.mood.is_depressed && rand.Intn(10) == 0 {
		d.kill_villager(v, "%s loses the will to exist -- "+d.mark.get_death())
		return
	}

	// Check if villager dies on the job
	// TODO

	// Check if villager starved
	if v.hunger == 0 {
		d.kill_villager(v, "%s starved to death.")
		return
	}
}

func (d *Death) kill_villager(villager *Person, reason string) {
	// Time to kick the bucket
	if reason == "" {
		reason = d.mark.get_death()
	}
	if strings.Contains(reason, "starved") {
		reason = "%s's hunger lead them to " + d.mark.get_death()
	}
	d.log = append(d.log, fmt.Sprintf(reason, "\u001b[31;1m"+villager.name+"\u001b[0m"))

	// Clean up lists
	for i, v := range d.pop.people {
		if v == villager {
			d.pop.people = append(d.pop.people[:i], d.pop.people[i+1:]...)
			break
		}
	}

	// d.pop.people.remove(villager)
	d.dead = append(d.dead, villager)

	// Clean up relationship objects
	for _, people := range d.pop.people {

		// Get rel value and remove the person
		rel_strength := people.relationships.del_relationship(villager)

		// Stronger relationships mean more sadness
		people.mood.death_event(rel_strength, people.name, villager.name, "")
	}
}
