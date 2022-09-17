package gamesheet

// Status is a character sheet status value that changes at a
// set rate and has a defined limit.
// TODO: Reduce the memory footprint of this struct.
// - Limit does not have to be a float.
type Status struct {
	Val   float32 // Current value.
	Rate  float32 // Increase / decrease per second.
	Limit float32 // Maximum value.
}

// NewStatus returns a new status with the given rate and limit.
func NewStatus() Status {
	return Status{Limit: 100.0}
}

// Tick advances the stat simulation by 'delta' (fraction of seconds)
// and returns true if it reaches the set limit.
func (s *Status) Tick(delta float64) bool {
	s.Val += s.Rate * float32(delta)
	if s.Val > s.Limit {
		s.Val = s.Limit
		return true
	}
	if s.Val < 0 {
		s.Val = 0
		return false
	}
	return false
}

// Add adds the given value to the status.
func (s *Status) Add(val float32) {
	s.Val += val
	if s.Val > s.Limit {
		s.Val = s.Limit
	} else if s.Val < 0 {
		s.Val = 0
	}
}

// Max returns if the value has reached its upper limit.
func (s *Status) Max() bool {
	return s.Val == s.Limit
}

const dayToSecond = 24 * 60 * 60

// State represents a specific physical state (awake, asleep, ...)
type State struct {
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
		Exhaustion: 100.0 / (4 * dayToSecond),  // We die after 4 days without rest.
		Hunger:     100.0 / (10 * dayToSecond), // We die after 10 days without food (it might be way longer, but meh).
		Thirst:     100.0 / (3 * dayToSecond),  // We die after 3 days without water.
		Stress:     100.0 / (2 * dayToSecond),  // We die after 2 days of stress.
	}
)
