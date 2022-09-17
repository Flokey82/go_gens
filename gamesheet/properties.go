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
