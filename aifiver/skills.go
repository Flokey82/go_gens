package aifiver

// Skill for various aspects
type Skill int

// Various skill types.
const (
	TSkillDiplomacy Skill = iota // Negotiation
	TSkillMartial                // Tactical skills
	TSkillStewardship
	TSkillIntrigue
	TSkillLearning
	TSkillProwess
	TSkillHealth
	TSkillParenting
)

// skillToString is a mapping from skill to string.
var skillToString = map[Skill]string{
	TSkillDiplomacy:   "diplomacy",
	TSkillMartial:     "martial",
	TSkillStewardship: "stewardship",
	TSkillIntrigue:    "intrigue",
	TSkillLearning:    "learning",
	TSkillProwess:     "prowess",
	TSkillHealth:      "health",
	TSkillParenting:   "parenting",
}

// String implements the stringer function for a skill.
func (s *Skill) String() string {
	if str, ok := skillToString[*s]; ok {
		return str
	}
	return "unknown"
}
