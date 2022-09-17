package gamesheet

// Slider represents a variable value with a variable upper bound.
// Example: Health, mana, etc.
type Slider [2]uint16

// NewSlider returns a new slider with the given value and maximum.
func NewSlider(max uint16) Slider {
	return Slider{max, max}
}

// Add adds the given value to the slider.
func (s *Slider) Add(val int) {
	s.SetValue(int(s[0]) + val)
}

// Value returns the slider's value.
func (s *Slider) Value() uint16 {
	return s[0]
}

// SetValue sets the slider's value.
func (s *Slider) SetValue(val int) {
	// Protect against overflow and underflow.
	if val > int(s[1]) {
		s[0] = s[1]
		return
	}
	if val < 0 {
		s[0] = 0
		return
	}
	s[0] = uint16(val)
}

// Max returns the slider's maximum value.
func (s *Slider) Max() uint16 {
	return s[1]
}

// SetMax sets the slider's maximum value.
func (s *Slider) SetMax(val uint16) {
	s[1] = val
}

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
func NewStatus(rate, limit float32) Status {
	return Status{Rate: rate, Limit: limit}
}

// Tick advances the stat simulation by 'delta' milliseconds.
func (s *Status) Tick(delta int64) {
	s.Val += s.Rate * float32(delta) / 1000
	if s.Val > s.Limit {
		s.Val = s.Limit
	}
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

// Attribute represents a character attribute.
type Attribute byte

// Add adds the given value to the attribute.
func (a *Attribute) Add(val int) {
	// Protect against overflow and underflow.
	res := int(*a) + val
	if res > 255 {
		*a = 255
		return
	}
	if res < 0 {
		*a = 0
		return
	}
	*a = Attribute(res)
}
