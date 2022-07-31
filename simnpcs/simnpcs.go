package simnpcs

import (
	"fmt"
	"log"
)

type Index struct {
	Entries   []*Character
	Locations []*Location
	Topics    *TopicPool
	TickCount uint64
	IDCount   uint64
}

func (idx *Index) NewProfession(name string, req LocationType) *Profession {
	prof := NewProfession(idx.GetID(), name, req)
	return prof
}

func (idx *Index) NewLocation(name string, parent *Location, t LocationType, s LocationScale) *Location {
	loc := NewLocation(idx.GetID(), name, t, s)
	idx.Locations = append(idx.Locations, loc)
	if parent != nil {
		parent.AssignChild(loc)
	}
	return loc
}

func (idx *Index) StartCareer(c *Character, p *Profession, l *Location) {
	// TODO: Account for change of workplace, retain experience.
	car := &Career{
		ID:         idx.GetID(),
		Start:      int(idx.TickCount),
		Profession: p,
		Storage:    newInventory(),
		Location:   l,
	}
	// Set up schedule based on typical working hours.
	for _, dow := range p.TypicalDays {
		for i := p.TypicalStart; i < p.TypicalEnd; i++ {
			car.WorkingHours[dow][i] = true
		}
	}
	c.SetCareer(car)
}

func New() *Index {
	idx := &Index{
		Topics: NewTopicPool(),
	}

	return idx
}

func (idx *Index) GetID() uint64 {
	idx.IDCount++
	return idx.IDCount
}

func (idx *Index) Tick() {
	// Move the simulation forward on tick.
	idx.TickCount++

	// Generate random world events.
	// TODO

	// Global plot points.
	// TODO

	log.Println(fmt.Sprintf("tick %d (Day %d Hour %d)", idx.TickCount, (idx.TickCount/24)%7, idx.TickCount%24))
	day := int(idx.TickCount / 24)
	hour := int(idx.TickCount % 24)
	locs := make(map[*Location][]*Character)

	// Get active routines and group them by location.
	for i := range idx.Entries {
		c := idx.Entries[i]
		// Update location based on routines and career.
		c.DoYourThing(day, hour)

		// Identify canidates for encounters.
		//// Routines that overlap in location and purpose
		//// have a high chance of an encounter.
		//// Based on existing social bonds and psychological
		//// profile, knowledge might be exchanged, and/or
		//// new social bonds being formed.
		if c.Location != nil {
			// Group routines by location.
			locs[c.Location] = append(locs[c.Location], c)
		}
	}

	// Match up possible encounters by location.
	for key := range locs {
		for i := range locs[key] {
			c1 := locs[key][i]
			for j := range locs[key] {
				if i == j {
					continue
				}
				c2 := locs[key][j]
				c1.Interact(c2, key)
			}
		}
	}
}

type Routine struct {
	DayOfWeek int
	Hour      int
	Location  *Location
	C         *Character
}

type AcquiredFact struct {
	ID   uint64
	Fact *Fact
	// Who / who else supplied information
	// When / during what interaction(s)
	// What was discussed
	//// Who else knows?
	//
	// Method of acquisition:
	//// Casual knowledge
	//// Education
	//// Legend
	//// First hand
}

type Interaction struct {
	ID uint64
	// Who was part of the interaction.
	// What was the type of interaction?
	// What was the outcome of the interaction
	// Was there an exchange of information
}

type Education struct {
	ID uint64
}

// SpeechModel describes the factors that has an impact on how a
// Character expresses themselves.
type SpeechModel struct {
	// Psychological factors
	//// High / low openness to experience
	////// More / less inquisitive
	//// High / low conscientiousness
	////// More / less detail and careful descriptions
	//// High / low extraversion
	////// More / less trust and willingness to share information
	//// High / low agreeableness
	////// More / less stubbornness and confrontation
	//// High / low neuroticism
	///// More / less trust
	// Place of residence
	//// The longer the stay, the stronger is the influence of local dialects.
	// Education (to a degree)
	// Social status
	// Profession (past and present)
}

type Impact struct {
	Emotional   float64 // Positive / negative emotion
	Monitary    float64 // Monitary gain
	Information float64 // Value of information exchanged
}

type Opinion struct {
	Count  int    // Number of values in average
	Impact        // Recent impact
	Total  Impact // Total impact
}

func (o *Opinion) Change(imp Impact) {
	o.Count++

	// Update emotional impact.
	o.Total.Emotional = incrementalAvrg(o.Total.Emotional, imp.Emotional, o.Count)

	// Update monitary impact.
	o.Total.Monitary = incrementalAvrg(o.Total.Monitary, imp.Monitary, o.Count)

	// Update information impact.
	o.Total.Information = incrementalAvrg(o.Total.Information, imp.Information, o.Count)

	// Update emotional impact with a weighted change.
	o.Emotional = weightedAvrg(o.Emotional, imp.Emotional, 0.5)

	// Update monitary impact with a weighted change.
	o.Monitary = weightedAvrg(o.Monitary, imp.Monitary, 0.5)

	// Update information impact with a weighted change.
	o.Information = weightedAvrg(o.Information, imp.Information, 0.5)
}

func (o *Opinion) String() string {
	if o.Emotional < 0 {
		return "dislikes"
	} else if o.Emotional == 0 {
		return "doesn't mind"
	}
	return "likes"
}

func incrementalAvrg(oldVal, newVal float64, count int) float64 {
	return oldVal + (newVal-oldVal)/(float64(count)-1)
}

func weightedAvrg(oldVal, newVal, weightFactor float64) float64 {
	return oldVal + weightFactor*(newVal-oldVal)
}
