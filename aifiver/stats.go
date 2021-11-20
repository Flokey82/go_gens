package aifiver

// Opinion by type
// - Same trait
// - Opposite trait
// - Attraction
type Opinion int

// Fate that can befall an individual
// - Early death
// - Capture in battle
// - Death in battle
// - Accidental injury
// - Bad investment
// - Scheme fallacy
// - Scheme discovery
type Fate int

// AIMod influences a specific aspect of decision making
// - AI Boldness
// - AI Compassion
// - AI Energy
// - AI Greed
// - AI Honor
// - AI Rationality
// - AI Sociability
// - AI Vengefulness
type AIMod int

type Stats struct {
	Opinion map[Opinion]int // Opinion modifiers
	Fate    map[Fate]int    // Chance modifier for fates
	Skill   map[Skill]int   // Skill modifiers
	AI      map[AIMod]int   // AI modifiers
}

func NewStats() *Stats {
	return &Stats{
		Opinion: make(map[Opinion]int),
		Fate:    make(map[Fate]int),
		Skill:   make(map[Skill]int),
		AI:      make(map[AIMod]int),
	}
}

func (s *Stats) add(st *Stats) {
	for key, val := range st.Opinion {
		s.Opinion[key] += val
	}
	for key, val := range st.Fate {
		s.Fate[key] += val
	}
	for key, val := range st.Skill {
		s.Skill[key] += val
	}
	for key, val := range st.AI {
		s.AI[key] += val
	}
}
