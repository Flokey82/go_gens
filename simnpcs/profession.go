package simnpcs

import (
	"fmt"
	"log"
	"math/rand"
)

// Career is an instance of a profession associated with a person.
type Career struct {
	ID           uint64      // Unique ID
	Start        int         // Started career at cycle
	Active       int         // Number of active cycles (experience)
	End          int         // Last cycle of activity
	Profession   *Profession // Type of profession
	Location     *Location   // Main career location
	Storage      *Inventory  // Produced items, resources...
	WorkingHours [7][24]bool // Active working hours

	// C          *Character // Wo is working on stuff
	// WorkingFor *Character // Employed by someone?
	// Any apprentices / employees?

	// TODO: Change this to requested items list / orders.
	WorkingOn *ProductionTask // Task currently being worked on
}

// ProductionTask is a task that is being worked on.
type ProductionTask struct {
	*Item         // Item being produced
	Remaining int // Remaining cycles until completion
}

// IsWorkTime returns true if the given day of the week and hour is active for this profession.
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

// NeedsItems returns all items that are required to produce items for this profession.
func (c *Career) NeedsItems() []*Item {
	if c == nil {
		return nil
	}
	// This should only list items required for current projects.
	return c.BuysItems()
}

// BuysItems returns all items that can be bought by this profession.
func (c *Career) BuysItems() []*Item {
	var res []*Item
	for _, item := range c.Profession.CanCraft(c.Active) {
		res = append(res, item.Consumes...)
	}
	return res
}

// SellsItems returns all items that can be sold by this profession.
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
		// Determine all items that we can craft given the current experience.
		canCraft := c.Profession.CanCraft(c.Active)

		// Can we craft anything?
		if len(canCraft) > 0 {
			// Select a random item to craft.
			// TODO: Select based on availability of required items.
			item := canCraft[rand.Intn(len(canCraft))]

			// Start working on the selected item.
			c.WorkingOn = &ProductionTask{
				Item:      item,
				Remaining: item.RequiresTime,
			}
			log.Println(fmt.Sprintf("%q started working on %q", c.Profession.Name, c.WorkingOn.Name))
		}
	} else {
		// We're working on something, so continue working on it.
		c.WorkingOn.Remaining--

		// Did we finish working on it?
		if c.WorkingOn.Remaining <= 0 {
			// Add the produced item to our storage.
			c.Storage.Add(c.WorkingOn.Item.newInstance(1234567)) // TODO: Generate unique id.
			log.Println(fmt.Sprintf("%q produced %q", c.Profession.Name, c.WorkingOn.Name))

			// Reset the current production task.
			c.WorkingOn = nil
		}
	}
}

// Profession represents a profession like "smith", "farmer", "miner", etc.
type Profession struct {
	ID            uint64       // Unique ID
	Name          string       // Name of the profession.
	Requires      LocationType // Requires location type.
	TypicalStart  int          // Typical time the works starts (hour)
	TypicalEnd    int          // Typical time the works ends (hour)
	TypicalDays   []DayOfWeek  // Typical days of the week the profession is performed.
	RequiresItems []*Item      // Required items (speer f. hunter, plow f. farmer, ...)
	Novice        int          // Requires number of (active) cycles to be "Novice"
	Skilled       int          // Requires number of (active) cycles to be "Skilled"
	Expert        int          // Requires number of (active) cycles to be "Expert"
	Skills        []*Skill     // Skills that can be learned by this profession.
	// TODO: Base Salary?
	// TODO: Required daily routines (e.g. "farmer" needs to water plants every day)
	// TODO: Expertise should also factor in talent.
}

// NewProfession creates a new profession.
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

// Skill represents a skill that can be used to produce items.
type Skill struct {
	ID            uint64      // Unique ID of the skill.
	Name          string      // Name of the skill.
	CanProduce    []*Item     // Items that can be produced by this skill.
	Requires      *Profession // Profession required.
	MinExperience int         // Minimum number of active cycles needed
}

// AddSkill adds a new skill to the profession.
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

// CanCraft returns all items that can be crafted by this profession at the given experience level.
func (p *Profession) CanCraft(activeTics int) []*Item {
	var res []*Item
	for _, skill := range p.Skills {
		// Check if we have enough experience to produce the items associated
		// with the given skill.
		if activeTics >= skill.MinExperience {
			res = append(res, skill.CanProduce...)
		}
	}
	return res
}
