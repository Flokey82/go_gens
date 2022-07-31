package simnpcs

import (
	"fmt"
	"log"
	"math/rand"
)

type Career struct {
	ID           uint64
	Start        int         // Start cycle
	Active       int         // Active cycles
	End          int         // End cycle
	Profession   *Profession // Type of profession
	Location     *Location   // Main career location
	Storage      *Inventory  // Produced items, resources...
	WorkingHours [7][24]bool // Active working hours

	// C          *Character // Wo is working on stuff
	// WorkingFor *Character // Employed by someone?
	// Any apprentices / employees?

	// TODO: Change this to requested items list / orders.
	WorkingOn *ProductionTask
}

type ProductionTask struct {
	*Item
	Remaining int
}

func (c *Career) IsWorkTime(day, hour int) bool {
	if c == nil {
		return false
	}
	if hour > 23 {
		hour = hour % 23
		day = (day + 1) % 7
	}
	return c.WorkingHours[day][hour]
}

func (c *Career) NeedsItems() []*Item {
	if c == nil {
		return nil
	}
	// This should only list items required for current projects.
	return c.BuysItems()
}

func (c *Career) BuysItems() []*Item {
	var res []*Item
	for _, item := range c.Profession.CanCraft(c.Active) {
		res = append(res, item.Consumes...)
	}
	return res
}

func (c *Career) SellsItems() []*Item {
	if c == nil {
		return nil
	}
	return c.Profession.CanCraft(c.Active)
}

func (c *Career) Update() {
	// TODO: Improve style.
	if c.Profession == nil {
		return
	}

	if reqLoc := c.Profession.Requires; reqLoc != LocTypeNone && (c.Location == nil || c.Location.Type != reqLoc) {
		// TODO: Return issues needed to be resolved in order to transform
		// them into actionable tasks.
		return
	}

	// Use required items.
	//for _, reqItem := range c.RequiresItems {
	//	gotItem := c.getInventoryItem(reqItem)
	//	if gotItem == nil || gotItem.durability == 0 {
	//		// TODO: Return issues needed to be resolved in order to transform
	//		// them into actionable tasks.
	//		return
	//	}
	//	gotItem.durability--
	//}

	c.Active++
	// We're not working on anything, so select something to work on.
	if c.WorkingOn == nil {
		// We assume we have all resources needed for creating the Item.
		canCraft := c.Profession.CanCraft(c.Active)
		// Can we craft anything?
		if len(canCraft) > 0 {
			item := canCraft[rand.Intn(len(canCraft))]
			c.WorkingOn = &ProductionTask{
				Item:      item,
				Remaining: item.RequiresTime,
			}
			log.Println(fmt.Sprintf("%q started working on %q", c.Profession.Name, c.WorkingOn.Name))
		}
	} else {
		c.WorkingOn.Remaining--
		if c.WorkingOn.Remaining <= 0 {
			c.Storage.Add(c.WorkingOn.Item.newInstance(1234567)) // TODO: Generate unique id.
			log.Println(fmt.Sprintf("%q produced %q", c.Profession.Name, c.WorkingOn.Name))
			c.WorkingOn = nil
		}
	}
}

type Profession struct {
	ID       uint64
	Name     string       // Name of the profession.
	Requires LocationType // Requires location type.

	// Time of day
	TypicalStart int
	TypicalEnd   int
	TypicalDays  []DayOfWeek
	// Required routines
	// Required items (speer f. hunter, plow f. farmer, ...)
	RequiresItems []*Item

	// Base Salary?

	// TODO: Expertise should also factor in talent.
	Novice  int // Requires number of (active) cycles to be "Novice"
	Skilled int // Requires number of (active) cycles to be "Skilled"
	Expert  int // Requires number of (active) cycles to be "Expert"
	Skills  []*Skill
}

func NewProfession(id uint64, name string, req LocationType) *Profession {
	prof := &Profession{
		ID:           id,
		Name:         name,
		Requires:     req,
		TypicalStart: 9,
		TypicalEnd:   19,
		TypicalDays: []DayOfWeek{
			DayMonday,
			DayTuesday,
			DayWednesday,
			DayThursday,
			DayFriday,
			DaySaturday,
		},
	}
	return prof
}

func (p *Profession) AddSkill(id uint64, name string, produce []*Item, minExp int) {
	skill := &Skill{
		ID:            id,
		Name:          name,
		CanProduce:    produce,
		Requires:      p,
		MinExperience: minExp,
	}

	// Note that the given item can be produced by this profession.
	for _, item := range produce {
		item.ProducedBy = append(item.ProducedBy, p)
	}

	p.Skills = append(p.Skills, skill)
}

func (p *Profession) CanCraft(activeTics int) []*Item {
	var res []*Item
	for _, skill := range p.Skills {
		if activeTics >= skill.MinExperience {
			res = append(res, skill.CanProduce...)
		}
	}
	return res
}

type Skill struct {
	ID            uint64
	Name          string
	CanProduce    []*Item
	Requires      *Profession // Profession required.
	MinExperience int         // Minimum number of active cycles needed
}

type DayOfWeek int

const (
	DayMonday DayOfWeek = iota
	DayTuesday
	DayWednesday
	DayThursday
	DayFriday
	DaySaturday
	DaySunday
)
