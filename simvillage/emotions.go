package simvillage

import (
	"fmt"
	"math/rand"
)

// Mood manages daily moods and emotions
type Mood struct {
	mood         string
	is_depressed bool
	happy        int
	sad          int
	productivity float64
	mood_events  []*MoodEvent
	log          []string
}

func NewMood() *Mood {
	m := &Mood{}
	m.mood = ""
	m.is_depressed = false

	m.happy = 5 // 0-10 scale
	m.sad = 4   // 0-10 scale
	m.productivity = 1.0

	// Holds multi-day mood effects
	m.mood_events = nil

	m.update_mood()
	m.update_productivity()

	m.log = nil
	return m
}

func (m *Mood) Tick() []string {
	// Manage daily chance of feeling good or bad
	m.daily_mood()

	// Manage Mood events
	count_events := 0
	for count_events < len(m.mood_events) {
		a, b := m.mood_events[count_events].tick()
		delta_mood := [2]int{a, b}

		if delta_mood == [2]int{-1, -1} {
			//del m.mood_events[count_events]
			m.mood_events = append(m.mood_events[:count_events], m.mood_events[count_events+1:]...)
		} else {
			m.happy += delta_mood[0]
			m.sad += delta_mood[1]

			count_events++
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
	m.update_mood()
	m.update_productivity()

	cp_log := m.log
	m.log = nil
	return cp_log
}

func (m *Mood) death_event(rel_strength float64, s_name, o_name, txt string) {
	if rel_strength < 1.00 {
		m.__mood_event(1, 0, 2, fmt.Sprintf("%s is glad %s died.", s_name, o_name))
	} else if 1.00 < rel_strength && rel_strength < 2.00 {
		m.__mood_event(0, 1, 1, fmt.Sprintf("%s is indifferent to %ss death.", s_name, o_name))
	} else if 2.00 < rel_strength && rel_strength < 4.00 {
		m.__mood_event(0, 1, 3, fmt.Sprintf("%s is hurt over %ss death.", s_name, o_name))
	} else if 4.00 < rel_strength {
		m.__mood_event(-2, 2, 10, fmt.Sprintf("%s is profoundly damaged over %ss death.", s_name, o_name))
	}
}

// Crete a new mood event for this person
func (m *Mood) __mood_event(h_tot, s_tot, dur int, txt string) {
	m.mood_events = append(m.mood_events, NewMoodEvent(h_tot, s_tot, dur, txt))
}

func (m *Mood) mod_mood(happy, sad int) {
	// Modify the current moods
	m.happy += happy
	m.sad += sad

	if m.happy < 0 {
		m.happy = 0
	}
	if m.happy > 10 {
		m.happy = 10
	}

	if m.sad < 0 {
		m.sad = 0
	}
	if m.sad > 10 {
		m.sad = 10
	}
}

func (m *Mood) update_mood() {
	// Slap a label on the current emotional state
	if m.happy == m.sad {
		m.mood = "Indifferent"
	} else if m.happy > m.sad+3 {
		m.mood = "Joyous"
	} else if m.happy > m.sad {
		m.mood = "Happy"
	} else if m.sad > m.happy+2 {
		m.mood = "Sad"
	} else if m.sad > m.happy {
		m.mood = "Melancholic"
	}
	if (m.happy < 2) && (m.sad > 8) {
		m.mood = "Depressed"
		m.is_depressed = true
	} else {
		m.is_depressed = false
	}
}

func (m *Mood) update_productivity() {
	// Happy people are more productive
	if m.happy == m.sad {
		m.productivity = 1.0
	} else if m.happy > m.sad {
		m.productivity = 1.2
	} else {
		m.productivity = .75
	}
}

func (m *Mood) daily_mood() {
	// On any given day one can be happy or sad
	if rand.Float64() < AVG_HAPPY {
		// good day
		m.mod_mood(1, -1)
	} else {
		// bad day
		m.mod_mood(-1, 1)
	}
}

type MoodEvent struct {
	daily_happy int
	daily_sad   int
	duration    int
	elapsed     int
}

/*
   Moods can be effected by larger events like having a kid, losing
   a loved one, or getting a promotion at work. These last multiple
   days and effect sadness and happiness daily.
*/
func NewMoodEvent(daily_happy, daily_sad, duration int, text string) *MoodEvent {
	e := &MoodEvent{}
	e.daily_happy = daily_happy
	e.daily_sad = daily_sad
	e.duration = duration
	e.elapsed = 0
	return e
}

func (e *MoodEvent) tick() (int, int) {
	if e.elapsed <= e.duration {
		e.elapsed += 1
		return e.daily_happy, e.daily_sad
	}
	return -1, -1
}
