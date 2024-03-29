package gamecs

import (
	"log"
	"math/rand"

	"github.com/Flokey82/aistate"
)

// ProfessionType is the general type of a profession.
// (e.g.: Baker, farmer, butcher, ...)
type ProfessionType struct {
	Name     string
	CanCraft []*ItemType
}

// NewProfessionType returns a new profession of a given type.
func NewProfessionType(name string, canCraft ...*ItemType) *ProfessionType {
	return &ProfessionType{
		Name:     name,
		CanCraft: canCraft,
	}
}

func (p *ProfessionType) New(w *World, a *Agent, workshop *Location) *Profession {
	return &Profession{
		w:              w,
		ProfessionType: p,
		ai:             a.CompAi,
		workshop:       workshop,
	}
}

// Profession represents a career of an individual and performs
// tasks related to the production of items. Implements aistate.State.
type Profession struct {
	*ProfessionType
	w              *World
	ai             *CompAi
	CurrentProject *Project  // Current project? Would a queue be better?
	workshop       *Location // Workshop inventory
	Missing        []*ItemType
}

const StateTypeProfession aistate.StateType = 6

// TODO: Implement state that solely handles work activity.
func (s *Profession) Type() aistate.StateType {
	return StateTypeProfession
}

// Tick advances the tasks associated with the profession by the
// given time interval.
func (s *Profession) Tick(delta uint64) {
	if len(s.Missing) > 0 {
		return // TODO: Implement acquisition of missing materials.
	}
	if s.CurrentProject == nil {
		s.selectNewProject()
	} else {
		s.CurrentProject.Tick(delta)
		if s.CurrentProject.Complete {
			s.workshop.Add(s.CurrentProject.Produce.New(s.w, s.workshop.Pos))
			log.Println("completed item!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			s.CurrentProject = nil
		}
	}
}

// selectNewProject selects a new Item to produce.
func (s *Profession) selectNewProject() {
	// Do we have a current project?
	if s.CurrentProject == nil {
		// No: Select a new project
		s.CurrentProject = newProject(s.CanCraft[rand.Intn(len(s.CanCraft))])
	}

	// Any missing resources?
	s.Missing = nil

	// TODO: Nil pointer exception protection.
	for _, req := range s.CurrentProject.Produce.Requires {
		if !s.workshop.Has(req) {
			//   Yes: Go and get missing items to the task list
			s.Missing = append(s.Missing, req)
		}
	}
}

// OnEnter is called when the state machine switches
// to this state.
func (s *Profession) OnEnter() {
	log.Println("Start work")
	if len(s.CanCraft) == 0 {
		return
	}
	// Plan for the day.
	s.selectNewProject()
	// Work on current project.
}

// OnExit is called when the state machine switches
// to another state.
func (s *Profession) OnExit() {
	log.Println("Stop work")
}

// Project represents a production task.
type Project struct {
	Produce  *ItemType
	Progress uint64 // Amount of time invested
	Duration uint64
	Complete bool
}

// newProject returns a project to produce an item of the
// given ItemType.
func newProject(it *ItemType) *Project {
	return &Project{
		Produce:  it,
		Duration: 50,
	}
}

// Tick advances the project by the given duration.
func (p *Project) Tick(delta uint64) {
	if p.Complete {
		return
	}
	p.Progress += delta
	if p.Progress >= p.Duration {
		p.Complete = true
	}
}

/*
func (s *Profession) CanComplete(it *ItemType) bool {
	return s.CInventory.HasAll(it.Requires...)
}

func (s *Profession) NeedsFor(it *ItemType) []*ItemType {
	// magic
	return items
}

func (s *Profession) Complete(it *ItemType) bool {
	if !s.CanComplete(it) {
		return false
	}
	// TODO: Prevent partial drops.
	if !s.CInventory.RemoveAll(it.Requires...) {
		return false
	}
	// TODO: Drop project from queue?
	return s.CInventory.Add(it.New(/ current location /))
}
*/
