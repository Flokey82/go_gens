// Package gamesheet provides a minimal character sheet for agents, or players.
package gamesheet

// CharacterSheet represents a character sheet.
type CharacterSheet struct {
	CurrentXP uint16             // Collected XP for the current level.
	Level     byte               // Current level.
	HP        Slider             // Hit points.
	AP        Slider             // Action points.
	Attrs     [AttrMax]Attribute // Hunger, thirst, etc.
}

// New returns a new character sheet.
func New(maxHP, maxAP byte) *CharacterSheet {
	// TODO: Instead of setting a fixed start max for the sliders,
	// set it by level and a base value.
	return &CharacterSheet{
		HP: NewSlider(maxHP),
		AP: NewSlider(maxAP),
	}
}

// AddExperience adds experience to the character sheet.
func (c *CharacterSheet) AddExperience(xp uint16) {
	c.CurrentXP += xp
	if nextLvlXP := c.NextLevelXP(); c.CurrentXP >= nextLvlXP {
		// TODO: Set new max HP and AP.
		c.Level++
		c.CurrentXP -= nextLvlXP
	}
}

const levelUpXPBase = 255

// NextLevelXP returns the XP required for the next level.
// NOTE: The required XP double with each level.
// Sure, this could be a bit better balanced, a table
// or whatever, but this just works too well :)
//
// Of course this will overflow at some point, but YOLO.
func (c *CharacterSheet) NextLevelXP() uint16 {
	return levelUpXPBase << uint16(c.Level)
}

// Various attributes.
const (
	AttrExhaustion = iota
	AttrHunger
	AttrThirst
	AttrStress
	AttrMax
)

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

// Slider represents a variable value with a variable upper bound.
// Example: Health, mana, etc.
type Slider [2]byte

// NewSlider returns a new slider with the given value and maximum.
func NewSlider(max byte) Slider {
	return Slider{max, max}
}

// Add adds the given value to the slider.
func (s *Slider) Add(val int) {
	s.SetValue(int(s[0]) + val)
}

// Value returns the slider's value.
func (s *Slider) Value() byte {
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
	s[0] = byte(val)
}

// Max returns the slider's maximum value.
func (s *Slider) Max() byte {
	return s[1]
}

// SetMax sets the slider's maximum value.
func (s *Slider) SetMax(val byte) {
	s[1] = val
}
