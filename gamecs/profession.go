package gamecs

import (
	"log"
	"math/rand"

	"github.com/Flokey82/aistate"
)

type ProfessionType struct {
	Name     string
	CanCraft []*ItemType
}

func NewProfessionType(name string) *ProfessionType {
	return &ProfessionType{
		Name: name,
	}
}

func (p *ProfessionType) New(w *World, a *Agent, workshop *Location) *Profession {
	return &Profession{
		w:              w,
		ProfessionType: p,
		ai:             a.CAi,
		workshop:       workshop,
	}
}

type Profession struct {
	*ProfessionType
	w              *World
	ai             *CAi
	CurrentProject *Project  // Current project? Would a queue be better?
	workshop       *Location // Workshop inventory
	Missing        []*ItemType
}

const StateTypeProfession aistate.StateType = 6

// TODO: Implement state that solely handles work activity.
func (s *Profession) Type() aistate.StateType {
	return StateTypeProfession
}

func (s *Profession) Tick(delta uint64) {
	if len(s.Missing) > 0 {
		return
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

func (s *Profession) selectNewProject() {
	// Do we have a current project?
	if s.CurrentProject == nil {
		//   No: Select a new project
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

func (s *Profession) OnEnter() {
	log.Println("Start work")
	if len(s.CanCraft) == 0 {
		return
	}
	// Plan for the day.
	s.selectNewProject()
	// Work on current project.
}

func (s *Profession) OnExit() {
	log.Println("Stop work")
}

type Project struct {
	Produce  *ItemType
	Progress uint64 // Amount of time invested
	Duration uint64
	Complete bool
}

func newProject(it *ItemType) *Project {
	return &Project{
		Produce:  it,
		Duration: 50,
	}
}

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
