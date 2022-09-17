// Package gamesheet provides a minimal character sheet for agents, or players.
package gamesheet

// CharacterSheet represents a character sheet.
//
// TODO:
//   - Add conditions like poisoned, blinded, etc.
//   - Find a better way to handle max level (100).
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
func New(baseHP, baseAP byte) *CharacterSheet {
	return &CharacterSheet{
		BaseHP: baseHP,
		BaseAP: baseAP,
		HP:     NewSlider(uint16(baseHP)),
		AP:     NewSlider(uint16(baseAP)),
	}
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
		c.update()

		// Remove the XP required for the next level.
		c.CurrentXP -= nextLvlXP

		// Increase available skill points.
		c.SkillPoints += levelUpSkillPoints
	}
}

func (c *CharacterSheet) update() {
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
	//                    lvl * lvl * baseVal         primStat + (secStat/2)
	// newVal = baseVal + ------------------- + lvl * ----------------------
	//                        100 * 100                       127.5
	//
	// baseVal ......... starting value (0-255)
	// lvl ............. current level  (0-100)
	// primStat ........ primary stat   (0-255)
	// secStat ......... secondary stat (0-255)
	//
	// This would give us a max value of 810.
	//
	// HP is influenced by strength and resilience.
	//   A character can take more damage if he is
	//   strong and resilient.
	//
	// AP is influenced by dexterity and resilience.
	//   A character can take more actions if he is
	//   dexterous and resilient (is less prone to
	//   exhaustion).
	//
	//               lvl * lvl * baseAp         dex + (res/2)
	// ap = baseAp + ------------------ + lvl * -------------
	//                   100 * 100 			       127.5
	//
	// This would give us a max AP of 810
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

	// Calculate new max values with somewhat precise floating point.
	calcNewMax := func(baseVal, level, primaryStat, secondaryStat byte) uint16 {
		fBv := float32(baseVal)
		fLvl := float32(level)
		fMax := fBv + (fLvl * fLvl * fBv / 10000) + fLvl*(float32(primaryStat)+float32(secondaryStat)/2)/127.5
		return uint16(fMax)
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
