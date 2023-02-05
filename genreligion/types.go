package genreligion

// Religions are logically organized like this:
//
// Group (Folk, Organized, ...)
// -> Form (Polytheism, Dualism, ...)
// -> Type (Church, Cult, ...)
//
// The type is more of a descriptive name for the religion, e.g. "Church" or "Cult".
const (
	// Religion groups.
	GroupFolk      = "Folk"
	GroupOrganized = "Organized"
	GroupCult      = "Cult"
	GroupHeresy    = "Heresy"

	// Religion forms.
	FormShamanism       = "Shamanism"
	FormAnimism         = "Animism"
	FormAncestorWorship = "Ancestor worship"
	FormPolytheism      = "Polytheism"
	FormDualism         = "Dualism"
	FormMonotheism      = "Monotheism"
	FormNontheism       = "Non-theism"
	FormCult            = "Cult"
	FormDarkCult        = "Dark Cult"
	FormHeresy          = "Heresy"
	// FormNature = "Nature"
)

// Forms maps a religion group to religion forms with a weighed probability.
var Forms = map[string]map[string]int{
	GroupFolk: {
		FormShamanism:       2,
		FormAnimism:         2,
		FormAncestorWorship: 1,
		FormPolytheism:      2,
	},
	GroupOrganized: {
		FormPolytheism: 5,
		FormDualism:    1,
		FormMonotheism: 4,
		FormNontheism:  1,
	},
	FormCult: {
		FormCult:     1,
		FormDarkCult: 1,
	},
	FormHeresy: {
		FormHeresy: 1,
	},
}

// Types maps a religion form to religion types with a weighed probability.
var Types = map[string]map[string]int{
	FormShamanism: {
		"Beliefs":   3,
		"Shamanism": 2,
		"Spirits":   1,
	},
	FormAnimism: {
		"Spirits": 1,
		"Beliefs": 1,
	},
	FormAncestorWorship: {
		"Beliefs":     1,
		"Forefathers": 2,
		"Ancestors":   2,
	},
	FormPolytheism: {
		"Deities":  3,
		"Faith":    1,
		"Gods":     1,
		"Pantheon": 1,
	},
	FormDualism: {
		"Religion": 3,
		"Faith":    1,
		"Cult":     1,
	},
	FormMonotheism: {
		"Religion": 1,
		"Church":   1,
	},
	FormNontheism: {
		"Beliefs": 3,
		"Spirits": 1,
	},
	FormCult: {
		"Cult":    4,
		"Sect":    4,
		"Arcanum": 1,
		"Coterie": 1,
		"Order":   1,
		"Worship": 1,
	},
	FormDarkCult: {
		"Cult":      2,
		"Sect":      2,
		"Blasphemy": 1,
		"Circle":    1,
		"Coven":     1,
		"Idols":     1,
		"Occultism": 1,
	},
	FormHeresy: {
		"Heresy":      3,
		"Sect":        2,
		"Apostates":   1,
		"Brotherhood": 1,
		"Circle":      1,
		"Dissent":     1,
		"Dissenters":  1,
		"Iconoclasm":  1,
		"Schism":      1,
		"Society":     1,
	},
}
