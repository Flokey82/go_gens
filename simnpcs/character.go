package simnpcs

import (
	"fmt"
	"log"

	"github.com/Flokey82/aifiver"
)

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

const (
	DayTimeStart = 7
	DayTimeEnd   = 22
)

type CharacterStatus int

const (
	CharStatIdle CharacterStatus = iota
	CharStatWorking
	CharStatResting
	CharStatSleeping
)

const maxExhaustion = 12

type Character struct {
	ID         uint64          // Unique ID
	FirstName  string          // First name
	LastName   string          // Last name
	Title      string          // Title (Sir, Lady, ...), optional
	Exhaustion int             // Current exhaustion level (0-8)
	WakeAt     int             // Time of day to wake up
	SleepAt    int             // Time of day to go to sleep
	Status     CharacterStatus // Current status (idle, working, resting, sleeping, ...)

	aifiver.SmallModel // Personality

	Career      *Career     // Current career (maybe allow multiple careers?)
	Home        *Location   // Current home (where the character sleeps)
	Location    *Location   // Current location
	PastCareers []*Career   // TODO: Add reason for new career
	PastHomes   []*Location // TODO: Add reason for move

	// Social standing.
	// Superior / Underlings
	// Birthday / Gender
	// Passions
	// Hobby []*Career - Hobby gardening, chicken coop.
	// Beliefs
	// Affiliations
	// Social connections
	Knowledge map[*Topic][]AcquiredFact // Total knowledge
	Opinions  map[uint64]Opinion        // ID to opinion mapping
	Routines  [7][24]*Routine           // Fixed routines
	Sources   map[uint64][]*Location    // Where to find what
	Tasks                               // Current tasks
	Inventory *Inventory                // Personal inventory
}

// NewCharacter creates a new character.
func NewCharacter(id uint64, firstName, lastName string, p aifiver.SmallModel) *Character {
	return &Character{
		ID:         id,
		FirstName:  firstName,
		LastName:   lastName,
		SmallModel: p,
		Knowledge:  make(map[*Topic][]AcquiredFact),
		Opinions:   make(map[uint64]Opinion),
		Sources:    make(map[uint64][]*Location),
		WakeAt:     DayTimeStart,
		SleepAt:    DayTimeEnd,
		Inventory:  newInventory(),
	}
}

// Name of the character.
func (c *Character) Name() string {
	var prefix string
	if c.Title != "" {
		prefix = c.Title + " "
	}
	return prefix + c.FirstName + " " + c.LastName
}

// AddSources adds a list of locations to the sources map so the character can find the item.
func (c *Character) AddSources(item *Item, locs ...*Location) {
	knownLocs := make(map[*Location]bool)
	for _, loc := range c.Sources[item.ID] {
		knownLocs[loc] = true
	}
	for _, loc := range locs {
		if !knownLocs[loc] {
			c.Sources[item.ID] = append(c.Sources[item.ID], loc)
		}
	}
}

// Interact with another character.
func (c *Character) Interact(ct *Character, loc *Location) {
	// Determine compatibility between the two characters' personalities.
	compat := aifiver.Compatibility(&c.SmallModel, &ct.SmallModel)

	// TODO: Determine if an interaction is likely.
	log.Println(fmt.Sprintf("encounter between %q and %q: %d", c.Name(), ct.Name(), compat))

	// TODO: If the personalities allow, exchange topics of interest or gossip.

	// Determine the impact of the interaction.
	imp := Impact{
		Emotional: float64(compat),
	}

	// Change opinions of each other based on the impact of the interaction.
	op := c.ChangeOpinion(ct.ID, imp)
	log.Println(fmt.Sprintf("%q %s %q", c.Name(), op.String(), ct.Name()))

	opt := ct.ChangeOpinion(c.ID, imp)
	log.Println(fmt.Sprintf("%q %s %q", ct.Name(), opt.String(), c.Name()))

	// Buy any items we need.
	hasItems := make(map[*Item]bool)
	for _, item := range c.Career.SellsItems() {
		// Remember where stuff is sold.
		ct.AddSources(item, loc)
		hasItems[item] = true
	}
	// TODO: Reverse trade information!

	// Execute transactions.
	var completed []*Task
	for _, t := range ct.Tasks {
		// If the character has the item, sell it.
		if hasItems[t.Item] {
			if it := c.Career.Storage.Find(t.Item); it != nil {
				// There is still enough stock, so buy/sell.
				c.Career.Storage.Move(it, ct.Career.Storage)
				completed = append(completed, t)
				log.Println(fmt.Sprintf("%q sold %q to %q", c.Name(), t.Item.Name, ct.Name()))
			} else {
				// We might have sold out.
				log.Println(fmt.Sprintf("%q can not sell %q to %q", c.Name(), t.Item.Name, ct.Name()))
			}
		} else if len(c.Sources[t.Item.ID]) != 0 {
			// Exchange information about where to find the item.
			// TODO: Exchange of information should only occur if they like each other.
			ct.AddSources(t.Item, c.Sources[t.Item.ID]...)
		}
	}

	// Complete completed tasks.
	for _, t := range completed {
		ct.CompleteTask(t)
		log.Println(fmt.Sprintf("%q completed task %q", ct.Name(), t.String()))
	}
}

// Determine overlap in Topics.
func (c *Character) FindTopics(ct *Character) []*Topic {
	var res []*Topic
	for t := range c.Knowledge {
		if _, ok := ct.Knowledge[t]; ok {
			res = append(res, t)
		}
	}
	// TODO: Rank by likelyhood.
	return res
}

// Change opinion on an entity with the given ID.
func (c *Character) ChangeOpinion(id uint64, imp Impact) Opinion {
	op := c.Opinions[id]
	op.Change(imp)
	c.Opinions[id] = op
	return op
}

// WakeUp wakes up the character.
func (c *Character) WakeUp() {
	log.Println(fmt.Sprintf("%q woke up", c.Name()))
	c.Status = CharStatIdle
}

// Sleep puts the character to sleep.
func (c *Character) Sleep() {
	c.GoTo(c.Home)
	log.Println(fmt.Sprintf("%q fell asleep", c.Name()))
	c.Status = CharStatSleeping
}

// Work causes the character to work.
func (c *Character) Work() {
	// Go to work location.
	c.GoTo(c.Career.Location)
	if c.Status != CharStatWorking {
		log.Println(fmt.Sprintf("%q started work", c.Name()))
		c.Status = CharStatWorking
	}
	c.Career.Update()
}

// Idle sets the character to idle.
func (c *Character) Idle() {
	// TODO: Use "Passions" to determine spare time activity.
	// TODO: What if they do not have a home?
	c.GoTo(c.Home)
	if c.Status != CharStatIdle {
		log.Println(fmt.Sprintf("%q is now idle", c.Name()))
		c.Status = CharStatIdle
	}
}

// DoYourThing causes the character to do watever is expected for the given day and hour.
func (c *Character) DoYourThing(day int, hour int) {
	dayOfWeek := day % 7
	if hour == c.WakeAt {
		c.WakeUp()
		c.Plan()
	} else if hour == c.SleepAt {
		c.Sleep()
	} else if r := c.GetRoutine(dayOfWeek, hour); r != nil {
		c.GoTo(r.Location)
		r.Location.Visit(c)
	} else if c.Career.IsWorkTime(dayOfWeek, hour) {
		c.Work()
	} else if c.Status != CharStatSleeping { // Fix this.
		c.Idle()
	}
}

// SetCareer sets a new active career for the character.
func (c *Character) SetCareer(car *Career) {
	// TODO: Account for change of workplace, retain experience.
	if c.Career != nil {
		c.PastCareers = append(c.PastCareers, c.Career)
	}
	c.Career = car

	// Set host of location
	if car.Location != nil {
		car.Location.Host = c
	}
}

// AddRoutine adds a specific routine for the character.
func (c *Character) AddRoutine(r *Routine) {
	c.Routines[r.DayOfWeek][r.Hour] = r
}

// GetRoutine gets any routine for the given day and hour.
func (c *Character) GetRoutine(dayOfWeek int, hour int) *Routine {
	return c.Routines[dayOfWeek][hour]
}

// GoTo moves the character to the given location.
func (c *Character) GoTo(loc *Location) {
	if loc != nil && loc != c.Location {
		log.Println(fmt.Sprintf("%q goes to %q", c.Name(), loc.Name))
		c.Location = loc
	}
}

// Plan creates a number of tasks for the character on the current day.
func (c *Character) Plan() {
	// Wake up at planned time.
	// What do I have to do tomorrow?
	// If I have a job, when do I start?
	// If I have an errand, when will it be?
	// When do I have to go to bed the latest to make it in time?
	//// Set planned wake time / Set planned bed time.

	// Is there something I need to get for my job?
	// TODO: Group by source and plan to visit.
	c.Tasks = nil
	for _, item := range c.Career.NeedsItems() {
		c.AddTask(0, item, TaskFind)
	}
}
