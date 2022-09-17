// Package gamesheet provides a minimal character sheet for agents, or players.
package gamesheet

// CharacterSheet represents a character sheet.
//
// TODO:
//   - Add conditions like poisoned, blinded, etc.
//   - Find a better way to handle max level (100).
//   - Handle stats.
//
// NOTE TO SELF:
// Do we allow the attribute values to change? If so, will retroactively
// the HP and AP increase or decrease? Would that even matter?
type CharacterSheet struct {
	CurrentXP   uint16 // Collected XP for the current level.
	Level       byte   // Current level.
	SkillPoints byte   // Skill points to distribute.
	BaseHP      byte   // Level 0 HP, will be used to calculate leveled HP.
	BaseAP      byte   // Level 0 AP, will be used to calculate leveled AP.
	HP          Slider // Hit points.
	AP          Slider // Action points.

	// Physical stats.
	StatExhaustion Attribute
	StatHunger     Attribute
	StatThirst     Attribute
	StatStress     Attribute

	// Physical attributes.
	AttrStrength     Attribute
	AttrIntelligence Attribute
	AttrDexterity    Attribute
	AttrResilience   Attribute
}

// New returns a new character sheet with the given base HP and AP.
//
// NOTE: The base HP and AP are the unleveled values. Depending on the
// character's stats, the HP and AP will increase as the character levels up
// or if a starting level > 0 has been set.
func New(baseHP, baseAP, level, str, itl, dex, res byte) *CharacterSheet {
	c := &CharacterSheet{
		BaseHP:           baseHP,
		BaseAP:           baseAP,
		Level:            level,
		SkillPoints:      levelUpSkillPoints * level,
		HP:               NewSlider(uint16(baseHP)),
		AP:               NewSlider(uint16(baseAP)),
		AttrStrength:     Attribute(str),
		AttrIntelligence: Attribute(itl),
		AttrDexterity:    Attribute(dex),
		AttrResilience:   Attribute(res),
	}
	c.Update()
	return c
}

// AddExperience adds experience to the character sheet.
func (c *CharacterSheet) AddExperience(xp uint16) {
	if c.Level >= maxLevel {
		return
	}
	c.CurrentXP += xp
	if nextLvlXP := c.NextLevelXP(); c.CurrentXP >= nextLvlXP {
		// Level up.
		c.Level++

		// Set new max HP and AP.
		c.Update()

		// Remove the XP required for the next level.
		c.CurrentXP -= nextLvlXP

		// Increase available skill points.
		c.SkillPoints += levelUpSkillPoints
	}
}

// Update recalculates stats like HP and AP based on the current
// level and attributes.
//
// Call this function if any of the attributes change to update
// the stats.
func (c *CharacterSheet) Update() {
	// Calculate new max HP and AP.
	//
	// Since resilience has an impact on both HP and AP,
	// we need to ensure that it has less of an impact
	// than dexterity and strength.
	//
	// We also need to keep in mind that the starting HP
	// and AP can range from 1 to 255, so we need to make
	// sure that by level 100, we are not somewhere in
	// crazy numbers like 65000.
	//
	//              l * l * v       p + (s/2)
	// newVal = v + --------- + l * ---------
	//              100 * 100         127.5
	//
	// v ......... starting value (0-255)
	// l ......... current level  (0-100)
	// p ......... primary stat   (0-255)
	// s ......... secondary stat (0-255)
	//
	// This would give us a max value of 810.
	//
	// In theory we could scale this up more if 810 is too low
	// for our game. But I think we should be fine with that.
	//
	// As a sidenote, Fallout has 320-440 HP and 70-110 AP.
	// We might reduce the AP if needed, but we can just make
	// actions more expensive. I think the main reason is to
	// make the ranges feel less "rule-ish" and more distinct.
	//
	// A quick Google gave me 460 HP as somewhat the max for
	// Dungeons and Dragons (not sure which class or edition).
	//
	// HP is influenced by strength and resilience.
	//   A character can take more damage if he is
	//   strong and resilient.
	//
	// AP is influenced by dexterity and resilience.
	//   A character can take more actions if he is
	//   dexterous and resilient (is less prone to
	//   exhaustion).

	// Calculate new max values with somewhat precise floating point.
	calcNewMax := func(baseVal, level, primaryStat, secondaryStat byte) uint16 {
		v := float32(baseVal)
		l := float32(level)
		p := float32(primaryStat)
		s := float32(secondaryStat)
		return uint16(v + (l * l * v / 10000) + l*(p+s/2)/127.5)
	}

	// Set new max values.
	c.HP.SetMax(calcNewMax(c.BaseHP, c.Level, byte(c.AttrStrength), byte(c.AttrResilience)))
	c.AP.SetMax(calcNewMax(c.BaseAP, c.Level, byte(c.AttrDexterity), byte(c.AttrResilience)))
}

const (
	levelUpSkillPoints = 2
	levelUpXPBase      = 6
	levelUpXPVariation = 51
	maxLevel           = 100
)

// NextLevelXP returns the XP required for the next level.
//
// NOTE: The required XP increases exponentially up to
// about 65100 XP for level 100. This way we make the
// best use of the unsigned 16 bit integer.
//
// Of course we could make the formula a bit more complicated
// so it doesn't irritate our more math-savvy players.
//
// See: https://pavcreations.com/level-systems-and-character-growth-in-rpg-games/
func (c *CharacterSheet) NextLevelXP() uint16 {
	nextLvl := uint16(c.Level + 1)
	return levelUpXPBase*nextLvl*nextLvl + levelUpXPVariation*nextLvl
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
