// Package simmemory provides a sample implementation of toughts/emotions and their storage in short, long, and core memory,
// based on the Dwarf Fortress wiki: http://dwarffortresswiki.org/index.php/DF2014:Memory_(thought)
package simmemory

import (
	"log"
	"math/rand"
)

const (
	daysToLongTerm = 10  // Number of days until thoughts are moved to long term memory
	daysToCore     = 255 // Number of days until thoughts are moved to core memory
)

// Thought represents a type of thought in the range of 0-255 (a single byte).
type Thought byte

// Predefined thoughts.
// NOTE: It'd be more useful to have this customizable, but for now we'll just
// use a few predefined thoughts.
const (
	ThoughtNone       Thought = 0  // No thought
	ThoughtNewFriend  Thought = 1  // Found a new friend
	ThoughtNewEnemy   Thought = 2  // Found a new enemy
	ThoughtNewPet     Thought = 3  // Found a new pet
	ThoughtLostFriend Thought = 4  // Lost a friend
	ThoughtLostEnemy  Thought = 5  // Lost an enemy
	ThoughtLostPet    Thought = 6  // Lost a pet
	ThoughtNewJob     Thought = 7  // Found a new job
	ThoughtLostJob    Thought = 8  // Lost a job
	ThoughtPromoted   Thought = 9  // Got promoted
	ThoughtDemoted    Thought = 10 // Got demoted
	ThoughtNewBaby    Thought = 11 // Had a new baby
	ThoughtNewSpouse  Thought = 12 // Got married
	ThoughtLostChild  Thought = 13 // Lost a child
	ThoughtLostSpouse Thought = 14 // Lost a spouse
	ThoughtSick       Thought = 15 // Got sick
	ThoughtHealed     Thought = 16 // Got healed
	ThoughtLast       Thought = 16 // Last thought (for random generation)
)

var ThoughtStrings = [256]string{
	ThoughtNone:       "None",
	ThoughtNewFriend:  "made a new friend",
	ThoughtNewEnemy:   "made a new enemy",
	ThoughtNewPet:     "got a new pet",
	ThoughtLostFriend: "lost a friend",
	ThoughtLostEnemy:  "lost an enemy",
	ThoughtLostPet:    "lost a pet",
	ThoughtNewJob:     "got a new job",
	ThoughtLostJob:    "lost a job",
	ThoughtPromoted:   "got promoted",
	ThoughtDemoted:    "got demoted",
	ThoughtNewBaby:    "had a new baby",
	ThoughtNewSpouse:  "got married",
	ThoughtLostChild:  "lost a child",
	ThoughtLostSpouse: "lost a spouse",
	ThoughtSick:       "got sick",
	ThoughtHealed:     "got healed",
}

// A range of intensities for thoughts.
const (
	IntensityVeryNegative     = -120
	IntensityNegative         = -50
	IntensitySlightlyNegative = -10
	IntensityNeutral          = 0
	IntensitySlightlyPositive = 10
	IntensityPositive         = 50
	IntensityVeryPositive     = 120
)

// ThoughtIntensity maps a thought to its intensity.
// NOTE: In dwarf fortress, the intensity is a positive value, and
// relies on the connected emotion (which is based on the personality)
// to determine the actual intensity (and if the actual value is positive
// or negative). We don't have emotions implemented (since they rely on
// character traits), so we just use the intensity directly and allow
// positive and negative values.
//
// This means though that we have to mask the intensity when we
// compare positive with negative values, which is a downside.
var ThoughtIntensity = [256]int8{
	IntensityNeutral,      // ThoughtNone
	IntensityPositive,     // ThoughtNewFriend
	IntensityNegative,     // ThoughtNewEnemy
	IntensityPositive,     // ThoughtNewPet
	IntensityNegative,     // ThoughtLostFriend
	IntensityPositive,     // ThoughtLostEnemy
	IntensityNegative,     // ThoughtLostPet
	IntensityPositive,     // ThoughtNewJob
	IntensityNegative,     // ThoughtLostJob
	IntensityPositive,     // ThoughtPromoted
	IntensityNegative,     // ThoughtDemoted
	IntensityVeryPositive, // ThoughtNewBaby
	IntensityVeryPositive, // ThoughtNewSpouse
	IntensityVeryNegative, // ThoughtLostChild
	IntensityVeryNegative, // ThoughtLostSpouse
	IntensityNegative,     // ThoughtSick
	IntensityPositive,     // ThoughtHealed
}

const intensityMask = 0x7f

// Group represents a group of thoughts.
type Group byte

const (
	GroupNone Group = iota
	GroupSocial
	GroupWork
	GroupFamily
	GroupHealth
)

// ThoughtGroup maps a thought to its group.
//
// A groupw would determine what emotions are triggered by a thought,
// and determine what personality traits are used to determine the
// intensity of the thought.
var ThoughtGroup = [256]Group{
	GroupNone,   // ThoughtNone
	GroupSocial, // ThoughtNewFriend
	GroupSocial, // ThoughtNewEnemy
	GroupSocial, // ThoughtNewPet
	GroupSocial, // ThoughtLostFriend
	GroupSocial, // ThoughtLostEnemy
	GroupSocial, // ThoughtLostPet
	GroupWork,   // ThoughtNewJob
	GroupWork,   // ThoughtLostJob
	GroupWork,   // ThoughtPromoted
	GroupWork,   // ThoughtDemoted
	GroupFamily, // ThoughtNewBaby
	GroupFamily, // ThoughtNewSpouse
	GroupFamily, // ThoughtLostChild
	GroupFamily, // ThoughtLostSpouse
	GroupHealth, // ThoughtSick
	GroupHealth, // ThoughtHealed
}

// Memory represents the short, long, and core memory of a creature.
// Short-, and long term memory have 8 slots, each of which can hold a thought,
// while core memory has 32 slots.
//
// We also track the number of days a thought has been in memory, as
// we promote thoughts from short to long memory after 10 days, and
// from long to core memory after 255 days.
//
// There are certain conditions that need to be fulfilled for a thought
// to be promoted. If promotion is not possible, the thought is
// discarded.
//
// TODO: In theory we could store thoughts and their age in the same
// array and use a step size of 2. This might be more efficient, but
// I don't really know if that's true.
type Memory struct {
	Short    [8]Thought
	Long     [8]Thought
	Core     [32]Thought
	AgeShort [8]byte
	AgeLong  [8]byte
}

// NewMemory returns a new memory.
func NewMemory() *Memory {
	return &Memory{}
}

// Log logs the state of the memory.
func (m *Memory) Log() {
	log.Println("Short term memory:")
	for i, thought := range m.Short {
		if thought == ThoughtNone {
			continue
		}
		log.Printf("  %d: %s (%d)", i, ThoughtStrings[thought], m.AgeShort[i])
	}
	log.Println("Long term memory:")
	for i, thought := range m.Long {
		if thought == ThoughtNone {
			continue
		}
		log.Printf("  %d: %s (%d)", i, ThoughtStrings[thought], m.AgeLong[i])
	}
	log.Println("Core memory:")
	for i, thought := range m.Core {
		if thought == ThoughtNone {
			continue
		}
		log.Printf("  %d: %s", i, ThoughtStrings[thought])
	}
}

// Tick advances the memory by one day.
func (m *Memory) Tick() {
	for i := 0; i < 8; i++ {
		if m.Short[i] != 0 {
			m.AgeShort[i]++
		}
		if m.Long[i] != 0 {
			m.AgeLong[i]++
		}
	}
	// First we check if we can promote any long-term thoughts to core memory,
	// which will free up a slot in long term memory.
	for i, age := range m.AgeLong {
		if age >= daysToCore {
			m.PromoteToCore(m.Long[i])
			m.Long[i] = 0
			m.AgeLong[i] = 0
		}
	}
	// Then we check if we can promote any short-term thoughts to long memory,
	// which will free up a slot in short term memory.
	for i, age := range m.AgeShort {
		if age >= daysToLongTerm {
			m.PromoteToLong(m.Short[i])
			m.Short[i] = 0
			m.AgeShort[i] = 0
		}
	}
}

// PromoteToLong promotes a thought to long memory.
// Once a memory has remained in a short-term memory slot for 10 days it will attempt
// to be promoted to a long-term memory slot.
//
// There are 8 long-term memory slots, and the procedure works similarly to short-term
// memory allocations, with one important difference.
//
// When the attempt to promote is made, a check is first made to see if there is an
// empty slot, if there is an empty slot the memory will be promoted to that slot even
// if a memory of that group already exists in another long-term memory slot.
//
// It is possible (but very rare) to have more than one memory of the same group in
// long term memory.
//
// This cannot happen in short term memory. If there are no empty slots a check is
// made to see if an existing memory of the same group exists.
//
// The promotion will fail if the existing long term memory is stronger, or will overwrite
// if the existing long term memory is weaker.
//
// If there are no empty slots and no existing memory of the same group, then the weakest
// of the other existing memories in long term will be overwritten.
//
// When a short-term memory is promoted (or possibly fails to promote) to long term memory
// it leaves an empty slot in the dwarf's short-term memory.
//
// Due to the cycling of the weakest short-term memories, it tends to be the stronger emotions
// that cause memories to remain in short-term memory for long enough to be promoted.
//
// The effect of the promotion on the dwarf's short-term memory is that it 'purges' a slot,
// allowing for a relatively weaker emotion to stick around without being overwritten by the
// cycling.
//
// Long-term memories are important and can be particularly impactful on a dwarf's mood because:
// 1) if a dwarf is frequently experiencing the same thing, good or bad, the same emotion
// can easily exist in both short term and long term, effectively doubling its impact
// 2) long term memories are often revisited long after an experience has ceased to occur
// 3) long term memories can become clogged with thoughts that can't be promoted further
func (m *Memory) PromoteToLong(t Thought) {
	// First we check if we can promote to an empty slot.
	for i, lt := range m.Long {
		if lt == 0 {
			m.Long[i] = t
			m.AgeLong[i] = 0
			return
		}
	}
	// Then we check if we can promote to an existing slot.
	for i, lt := range m.Long {
		if ThoughtGroup[lt] == ThoughtGroup[t] {
			if ThoughtIntensity[lt]&intensityMask < ThoughtIntensity[t]&intensityMask {
				m.Long[i] = t
				m.AgeLong[i] = 0
			}
			return
		}
	}
	// Finally we overwrite the weakest thought.
	weakest := 0
	for i := 1; i < 8; i++ {
		if ThoughtIntensity[m.Long[i]]&intensityMask < ThoughtIntensity[m.Long[weakest]]&intensityMask {
			weakest = i
		}
	}
	m.Long[weakest] = t
	m.AgeLong[weakest] = 0
}

// PromoteToCore promotes a thought to core memory.
// Once a memory has remained in long-term memory for 255 days, it has a 1:3 chance
// of being promoted to core memory and causing one or more personality changes, if
// the memory is of a group that can be promoted to core memory.
//
// There are 32 core memory slots, and thoughts are promoted to the first empty
// slot. If there are no empty slots, the weakest thought is overwritten.
//
// Core memories are less impactful than long-term memories, as they are rarely,
// if ever, revisited. However, the change that is made to the personality of the
// dwarf is permanent, and this can be for good or bad.
func (m *Memory) PromoteToCore(t Thought) {
	// There is a 1:3 chance of promoting a long-term memory to core memory.
	if rand.Intn(3) != 0 {
		return
	}

	// First we check if we can promote to an empty slot.
	for i, ct := range m.Core {
		if ct == 0 {
			m.Core[i] = t
			return
		}
	}
	// Then we check if we can promote to an existing slot.
	for i, ct := range m.Core {
		if ThoughtGroup[ct] == ThoughtGroup[t] {
			if ThoughtIntensity[ct]&intensityMask < ThoughtIntensity[t]&intensityMask {
				m.Core[i] = t
			}
			return
		}
	}
	// Finally we overwrite the weakest thought.
	weakest := 0
	for i := 1; i < 32; i++ {
		if ThoughtIntensity[m.Core[i]]&intensityMask < ThoughtIntensity[m.Core[weakest]]&intensityMask {
			weakest = i
		}
	}
	m.Core[weakest] = t
}

// AddThought adds a thought to the dwarf's memory.
// A dwarf has 8 short-term memory slots. When a dwarf has a thought, a check is
// made to see if a memory of that group already exists in a short-term memory slot.
//
// If the thought doesn't fall into an existing group, the new thought will fill an
// empty slot, or if no slots are empty, overwrite the weakest memory (the one with
// the weakest emotion) in the 8 short-term memory slots - even if the overwritten
// memory is stronger than the new one.
//
// If the thought already has one of its group in a memory slot the strongest memory
// of the two, the new thought and the existing memory, will be kept and the other
// discarded.
//
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// NOTE: Below are examples that only apply once "emotion" is implemented, which
// effectively changes the intensity of a thought.
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//
// So if a dwarf is in a particularly bad mood, getting shat on by a bird might
// be felt more intensely than if they were in a good mood.
//
// For example, a dwarf gets caught in the rain and is dejected (intensity 1/4).
// If they haven't seen rain in the last year, this will be written to their memory,
// overwriting the weakest existing memory of a different group.
//
// If they have seen rain in the last year, it will overwrite the previously weaker
// emotion of being annoyed by the rain (intensity 1/8), or ignore the new experience
// if the old emotion was the stronger being dismayed by the rain (intensity 1/2).
//
// This mostly leads to a constant cycling of the weakest of the 8 memory slots as new
// thoughts overwrite each other. An overwritten thought is a forgotten thought.
//
// You can see this in the dwarf's thoughts and preferences screen,
// they'll mention more than 8 things that they have recently experienced,
// but they are only being affected by 8 of those things.
//
// The list is not a reflection of what is in their current memory, but is a reflection
// of what has been in their memory recently. The consequence of this cycling is that
// short-term memories are mostly fleeting, with a maximum of 7 short-term memories
// having a lasting effect on a dwarf's mood.
func (m *Memory) AddThought(t Thought) {
	// First we check if we have a slot with a memory of the same group.
	for i := 0; i < 8; i++ {
		if ThoughtGroup[m.Short[i]] == ThoughtGroup[t] {
			if ThoughtIntensity[m.Short[i]]&intensityMask < ThoughtIntensity[t]&intensityMask {
				m.Short[i] = t
				m.AgeShort[i] = 0
			}
			return
		}
	}
	// Then we check if we can add to an empty slot.
	for i := 0; i < 8; i++ {
		if m.Short[i] == 0 {
			m.Short[i] = t
			m.AgeShort[i] = 0
			return
		}
	}
	// Finally we overwrite the weakest thought.
	weakest := 0
	for i := 1; i < 8; i++ {
		if ThoughtIntensity[m.Short[i]]&intensityMask < ThoughtIntensity[m.Short[weakest]]&intensityMask {
			weakest = i
		}
	}
	m.Short[weakest] = t
	m.AgeShort[weakest] = 0
}

/*
// MemoryManager represents a customizable index for thought types,
// groups, and intensities.
// TODO: Make the intensity a int8 and mask the values when checking
// the intensity, so we can compare positive and negative values in
// their absolute value.
type MemoryManager struct {
	Thoughts     [256]Thought
	ThoughtGroup [256]byte // Maps a thought to its group.
	ThoughtInt   [256]byte // Maps a thought to its intensity.
}

// NewMemoryManager returns a new memory manager.
func NewMemoryManager() *MemoryManager {
	return &MemoryManager{}
}

// AddThought adds a new thought to the memory manager.
// 't': The thought (ID) to add.
// 'g': The group of the thought.
// 'i': The intensity of the thought.
func (m *MemoryManager) AddThought(t Thought, g byte, i byte) {
	m.Thoughts[t] = t
	m.ThoughtGroup[t] = g
	m.ThoughtInt[t] = i
}
*/
