package gamesheet

// Status is a character sheet status value that changes at a
// set rate.
// TODO: Reduce the memory footprint of this struct.
type Status struct {
	Val  float32 // Current value.
	Rate float32 // Increase / decrease per second.
}

// NewStatus returns a new status with the given rate and limit.
func NewStatus() Status {
	return Status{}
}

// Tick advances the stat simulation by 'delta' (fraction of seconds)
// and returns true if it reaches the set limit.
func (s *Status) Tick(delta float64) bool {
	if s.Val == 100.0 {
		return true
	}
	s.Add(s.Rate * float32(delta))
	return s.Val >= 100.0
}

// Add adds the given value to the status.
func (s *Status) Add(val float32) {
	s.Val += val
	if s.Val > 100.0 {
		s.Val = 100.0
	} else if s.Val < 0 {
		s.Val = 0
	}
}

// Max returns if the value has reached its upper limit.
func (s *Status) Max() bool {
	return s.Val >= 100.0
}

const (
	hourToSecond = 60 * 60
	dayToSecond  = 24 * hourToSecond
)

// State represents a specific physical state (awake, asleep, ...)
type State struct {
	Name       string
	Exhaustion float32
	Hunger     float32
	Thirst     float32
	Stress     float32
}

// Some constants related to stats.
//
// TODO: This should be on a per-creature basis.
// - A camel needs less water than a human.
// - A humpback whale survives 6 MONTHS without food!
//
// TODO: There should also be the recovery rate.
//   - After 7 hours of sleep, the exhaustion stat should be reduced by
//     a day's worth of exhaustion.
//   - One hour of rest should reduce stress significantly.
//   - While sleeping, hunger and thirst should increase much slower.
//   - During strenuous activity, hunger, thirst, and exhaustion should
//     increase much faster.
//   - When in combat and in danger, stress should increase.
var (
	StateAwake = &State{
		Name:       "awake",
		Exhaustion: 100.0 / (4 * dayToSecond),  // We die after 4 days without rest.
		Hunger:     100.0 / (10 * dayToSecond), // We die after 10 days without food (it might be way longer, but meh).
		Thirst:     100.0 / (3 * dayToSecond),  // We die after 3 days without water.
		Stress:     100.0 / (2 * dayToSecond),  // We die after 2 days of stress.
	}
	StateAsleep = &State{
		Name:       "asleep",
		Exhaustion: -100.0 / (8 * hourToSecond), // We should be fully rested after 8 hours of sleep.
		Hunger:     100.0 / (20 * dayToSecond),  // While asleep, we starve slower.
		Thirst:     100.0 / (20 * dayToSecond),  // While asleep, we dehydrate slower.
		Stress:     -100.0 / (8 * hourToSecond), // While asleep, we recover from stress faster.
	}
)
