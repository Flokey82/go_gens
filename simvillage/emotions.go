package simvillage

import (
	"fmt"
	"math/rand"
)

// Mood manages daily moods and emotions
type Mood struct {
	mood         string
	isDepressed  bool
	happy        int
	sad          int
	productivity float64
	moodEvents   []*MoodEvent
	log          []string
}

func NewMood() *Mood {
	m := &Mood{
		isDepressed:  false,
		happy:        5, // 0-10 scale
		sad:          4, // 0-10 scale
		productivity: 1.0,
		moodEvents:   nil, // Holds multi-day mood effects
	}

	m.updateMood()
	m.updateProductivity()
	return m
}

func (m *Mood) Tick() []string {
	// Manage daily chance of feeling good or bad
	m.dailyMood()

	// Manage Mood events
	countEvents := 0
	for countEvents < len(m.moodEvents) {
		a, b := m.moodEvents[countEvents].tick()
		delta_mood := [2]int{a, b}
		if delta_mood == [2]int{-1, -1} {
			//del m.mood_events[count_events]
			m.moodEvents = append(m.moodEvents[:countEvents], m.moodEvents[countEvents+1:]...)
		} else {
			m.happy += delta_mood[0]
			m.sad += delta_mood[1]
			countEvents++
		}
	}

	// People will gradually stabalize to 5/3 by default
	if rand.Intn(2) < 1 {
		if m.happy > 5 {
			m.happy -= 1
		}
		if m.happy < 5 {
			m.happy += 1
		}
		if m.sad > 3 {
			m.sad -= 1
		}
		if m.sad < 3 {
			m.sad += 1
		}
	}
	m.updateMood()
	m.updateProductivity()

	cp_log := m.log
	m.log = nil
	return cp_log
}

func (m *Mood) deathEvent(rel_strength float64, s_name, o_name, txt string) {
	if rel_strength < 1.00 {
		m._moodEvent(1, 0, 2, fmt.Sprintf("%s is glad %s died.", s_name, o_name))
	} else if 1.00 < rel_strength && rel_strength < 2.00 {
		m._moodEvent(0, 1, 1, fmt.Sprintf("%s is indifferent to %ss death.", s_name, o_name))
	} else if 2.00 < rel_strength && rel_strength < 4.00 {
		m._moodEvent(0, 1, 3, fmt.Sprintf("%s is hurt over %ss death.", s_name, o_name))
	} else if 4.00 < rel_strength {
		m._moodEvent(-2, 2, 10, fmt.Sprintf("%s is profoundly damaged over %ss death.", s_name, o_name))
	}
}

// Crete a new mood event for this person
func (m *Mood) _moodEvent(h_tot, s_tot, dur int, txt string) {
	m.moodEvents = append(m.moodEvents, NewMoodEvent(h_tot, s_tot, dur, txt))
}

func (m *Mood) modMood(happy, sad int) {
	// Modify the current moods
	m.happy += happy
	if m.happy < 0 {
		m.happy = 0
	} else if m.happy > 10 {
		m.happy = 10
	}

	m.sad += sad
	if m.sad < 0 {
		m.sad = 0
	} else if m.sad > 10 {
		m.sad = 10
	}
}

const (
	MoodIndifferent = "Indifferent"
	MoodJoyous      = "Joyous"
	MoodHappy       = "Happy"
	MoodSad         = "Sad"
	MoodMelancholic = "Melancholic"
	MoodDepressed   = "Depressed"
)

func (m *Mood) updateMood() {
	// Slap a label on the current emotional state
	if m.happy == m.sad {
		m.mood = MoodIndifferent
	} else if m.happy > m.sad+3 {
		m.mood = MoodJoyous
	} else if m.happy > m.sad {
		m.mood = MoodHappy
	} else if m.sad > m.happy+2 {
		m.mood = MoodSad
	} else if m.sad > m.happy {
		m.mood = MoodMelancholic
	}

	if (m.happy < 2) && (m.sad > 8) {
		m.mood = MoodDepressed
		m.isDepressed = true
	} else {
		m.isDepressed = false
	}
}

func (m *Mood) updateProductivity() {
	// Happy people are more productive
	if m.happy == m.sad {
		m.productivity = 1.0
	} else if m.happy > m.sad {
		m.productivity = 1.2
	} else {
		m.productivity = .75
	}
}

func (m *Mood) dailyMood() {
	// On any given day one can be happy or sad
	if rand.Float64() < AVG_HAPPY {
		m.modMood(1, -1) // good day
	} else {
		m.modMood(-1, 1) // bad day
	}
}

type MoodEvent struct {
	dailyHappy int
	dailySad   int
	duration   int
	elapsed    int
}

// Moods can be effected by larger events like having a kid, losing
// a loved one, or getting a promotion at work. These last multiple
// days and effect sadness and happiness daily.
func NewMoodEvent(dailyHappy, dailySad, duration int, text string) *MoodEvent {
	return &MoodEvent{
		dailyHappy: dailyHappy,
		dailySad:   dailySad,
		duration:   duration,
		elapsed:    0,
	}
}

func (e *MoodEvent) tick() (int, int) {
	if e.elapsed <= e.duration {
		e.elapsed += 1
		return e.dailyHappy, e.dailySad
	}
	return -1, -1
}
