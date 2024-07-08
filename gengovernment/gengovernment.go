// Package gengovernment allows to generate various types of leadership or forms of government.
package gengovernment

import (
	"math/rand"
)

// LeadershipInfluence represents the minimum and maximum influence a type of leadership can have.
type LeadershipInfluence int

// The different levels of influence a leader can have.
const (
	LeadershipInfluenceUnset      LeadershipInfluence = iota // The influence level is not set.
	LeadershipInfluenceTribe                                 // Leadership is limited to a single tribe or clan.
	LeadershipInfluenceSettlement                            // Leadership extends to a single settlement or village.
	LeadershipInfluenceCityState                             // Leadership extends to a city-state or small kingdom.
	LeadershipInfluenceEmpire                                // Leadership extends to a large empire or nation.
)

type LeadershipForm int

// The different forms of leadership.
const (
	LeadershipFormChiefdom LeadershipForm = iota
	LeadershipFormCouncil
	LeadershipFormMonarchy
	LeadershipFormRepublic
	LeadershipFormDemocracy
	LeadershipFormDictatorship
	LeadershipFormMax
)

// Chiefdom can either be a single chief or a council of chiefs.
// So variations could be:
// - Single Chief
// - Council of Chiefs
// - Single Chief with Council of Advisors

// String returns the string representation of the leadership form.
func (l LeadershipForm) String() string {
	switch l {
	case LeadershipFormChiefdom:
		return "Chiefdom"
	case LeadershipFormMonarchy:
		return "Monarchy"
	case LeadershipFormCouncil:
		return "Council"
	case LeadershipFormRepublic:
		return "Republic"
	case LeadershipFormDemocracy:
		return "Democracy"
	case LeadershipFormDictatorship:
		return "Dictatorship"
	}
	return "Unknown"
}

// RangeInfluence returns the minimum and maximum influence a type of leadership can have.
func (l LeadershipForm) RangeInfluence() (min, max LeadershipInfluence) {
	switch l {
	case LeadershipFormChiefdom:
		return LeadershipInfluenceTribe, LeadershipInfluenceSettlement
	case LeadershipFormCouncil:
		return LeadershipInfluenceTribe, LeadershipInfluenceSettlement
	case LeadershipFormMonarchy:
		return LeadershipInfluenceCityState, LeadershipInfluenceEmpire
	case LeadershipFormRepublic:
		return LeadershipInfluenceCityState, LeadershipInfluenceEmpire
	case LeadershipFormDemocracy:
		return LeadershipInfluenceCityState, LeadershipInfluenceEmpire
	case LeadershipFormDictatorship:
		return LeadershipInfluenceSettlement, LeadershipInfluenceEmpire
	}
	return LeadershipInfluenceUnset, LeadershipInfluenceUnset
}

// NaturalProgression returns the possible forms of government that can evolve from the current form.
func (l LeadershipForm) NaturalProgression() []LeadershipForm {
	return GovTree[l]
}

// CoupProgression returns the possible forms of government that can evolve from the current form after a coup or revolution.
func (l LeadershipForm) CoupProgression() []LeadershipForm {
	return GovTreeCoup[l]
}

// GovTree represents the possible, natural paths when developing / evolving a government.
var GovTree = map[LeadershipForm][]LeadershipForm{
	LeadershipFormChiefdom: {
		LeadershipFormCouncil,
		LeadershipFormMonarchy,
		LeadershipFormDictatorship,
	},
	LeadershipFormCouncil: {
		LeadershipFormRepublic,
	},
	LeadershipFormMonarchy: {
		LeadershipFormRepublic,
		LeadershipFormDictatorship,
	},
	LeadershipFormRepublic: {
		LeadershipFormDemocracy,
	},
}

// GovTreeCoup represents the possible paths when a coup or revolution occurs.
var GovTreeCoup = map[LeadershipForm][]LeadershipForm{
	LeadershipFormChiefdom: {
		LeadershipFormDictatorship,
	},
	LeadershipFormCouncil: {
		LeadershipFormDictatorship,
		LeadershipFormMonarchy,
	},
	LeadershipFormMonarchy: {
		LeadershipFormDictatorship,
		LeadershipFormRepublic,
	},
	LeadershipFormRepublic: {
		LeadershipFormDictatorship,
		LeadershipFormMonarchy,
	},
	LeadershipFormDemocracy: {
		LeadershipFormDictatorship,
		LeadershipFormMonarchy,
	},
	LeadershipFormDictatorship: {
		LeadershipFormRepublic,
		LeadershipFormCouncil,
		LeadershipFormMonarchy,
	},
}

// TitleTag represents a tag that can be associated with a title to provide
// additional context or meaning.
type TitleTag int

const (
	TitleTagOther TitleTag = iota
	TitleTagMilitary
	TitleTagService
	TitleTagWisdom
	TitleTagMythical
	TitleTagSensuality
	TitleTagPower
	TitleTagDivine
	TitleTagGrandeur
	TitleTagMagic
	TitleTagCruelty
	TitleTagMax
)

// Title represents a potential title for a ruler or council member and
// includes the male and female variants, as well as tags for the context.
type Title struct {
	Variants [2]string
	Tags     []TitleTag
}

type TitleSet struct {
	Superior []Title // Titles for a superior or ruler like an emperor or dictator.
	Head     []Title // Titles a ruler like a king or queen, or a president.
}

var councilTitles = map[TitleTag][][2]string{
	TitleTagMilitary: {
		{"Strateg", "Strategin"},     // A Germanic term for a military leader, emphasizing tactical prowess.
		{"Dominar", "Dominatrix"},    // A Latin term meaning master or lord, ideal for a council leader who demands absolute obedience.
		{"Commander", "Commander"},   // For a council leader who leads through military might and conquest.
		{"Warlord", "Warmistress"},   // For a council leader who rules through martial strength and conquest.
		{"Warmaster", "Warmistress"}, // For a council leader who rules through martial strength and conquest.
		{"General", "Generalissima"}, // Classic titles for a military leader, with a sense of authority and command.
		{"Marshal", "Marshalin"},     // A military leader who commands armies and enforces order.
	},
	TitleTagService: {
		{"Protector", "Protectress"},     // For a council leader who sees themselves as a guardian of the people.
		{"Guardian", "Guardian"},         // A simple and straightforward title for a ruler who protects their subjects.
		{"Steward", "Stewardess"},        // For a ruler who manages the affairs of the realm with care and responsibility.
		{"Warden", "Wardess"},            // For a ruler who oversees the safety and security of the land.
		{"Defender", "Defendress"},       // For a ruler who protects their people from external threats.
		{"Sentinel", "Sentinella"},       // For a ruler who stands watch over the realm, protecting it from danger.
		{"Shieldbearer", "Shieldmaiden"}, // For a ruler who defends their people with strength and courage.
		{"Paladin", "Paladiness"},        // For a ruler who upholds justice and righteousness, protecting the weak and innocent.
		{"Champion", "Championess"},      // For a ruler who fights for their people and leads them to victory.
		{"Hero", "Heroine"},              // For a ruler who is celebrated for their bravery and valor in battle.
		{"Knight", "Dame"},               // For a ruler who embodies the ideals of chivalry and honor.
	},
	TitleTagWisdom: {
		{"Sage", "Sibyl"},            // For a council leader who emphasizes wisdom and guidance.
		{"Mentor", "Mentress"},       // For a ruler who teaches and guides their people.
		{"Visionary", "Visionary"},   // For a ruler who sees beyond the present and envisions a better future.
		{"Guide", "Guidess"},         // For a ruler who leads their people with wisdom and insight.
		{"Pathfinder", "Pathfinder"}, // For a ruler who discovers new paths and opportunities for their people.
		{"Wayfinder", "Wayfinder"},   // For a ruler who navigates the challenges of leadership with skill and grace.
		{"Luminary", "Luminaria"},    // For a ruler who shines with knowledge and enlightenment.
		{"Scholar", "Scholarina"},    // For a ruler who values learning and education, guiding their people to greater understanding.
	},
	TitleTagMythical: {
		{"Oracle", "Oracle"},         // For a ruler believed to have prophetic abilities.
		{"Prophet", "Prophetess"},    // For a council leader who claims to speak for the gods or foresee the future.
		{"Seer", "Seeress"},          // For a ruler who has the gift of foresight and prophecy.
		{"Diviner", "Diviner"},       // For a ruler who interprets signs and omens to guide their people.
		{"Augur", "Auguress"},        // For a ruler who reads the signs of nature and predicts the future.
		{"Harbinger", "Harbinger"},   // For a ruler who foretells the coming of great events and changes.
		{"Shaman", "Shamaness"},      // For a ruler who communes with the spirits and seeks guidance from the otherworldly.
		{"Druid", "Druidess"},        // For a ruler who draws power from the natural world and the spirits of the land.
		{"Soothsayer", "Soothsayer"}, // For a ruler who foretells the fate of individuals and nations.
	},
	TitleTagOther: {
		{"Elder", "Elderess"},          // For a council leader who is respected for their age and experience.
		{"Coucilman", "Councilwoman"},  // For a council leader who serves on a council or governing body.
		{"Councilor", "Councilress"},   // For a council leader who advises and guides the ruler.
		{"Advisor", "Advisor"},         // For a council leader who offers counsel and wisdom to the ruler.
		{"Consul", "Consul"},           // For a council leader who represents the interests of the people.
		{"Legate", "Legatress"},        // For a council leader who acts as an ambassador or envoy.
		{"Lictor", "Lictress"},         // For a council leader who enforces the laws and decrees of the ruler.
		{"Magistrate", "Magistratess"}, // For a council leader who administers justice and maintains order.
		{"Senator", "Senatress"},       // For a council leader who serves in a legislative body.
		{"Speaker", "Speaker"},         // For a council leader who speaks on behalf of the people.
		{"Burgess", "Burgess"},         // For a council leader who represents a district or constituency.
		{"Alderman", "Alderwoman"},     // For a council leader who serves as a municipal official.
	},
}

var monarchyTitles = map[TitleTag][][2]string{
	TitleTagOther: {
		{"King", "Queen"},       // Classic titles for a supreme ruler, with a sense of grandeur and authority.
		{"Konungr", "Konung"},   // Old Norse words for king / queen, perfect for a harsh and brutal ruler.
		{"Rex", "Regina"},       // Latin terms for king / queen, with a sense of ancient authority and power.
		{"Regent", "Regentess"}, // While traditionally a temporary ruler, it can be twisted for a dictator who seized power and refuses to relinquish it.
	},
}

var democracyTitles = map[TitleTag][][2]string{
	TitleTagOther: {
		{"President", "President"},           // For a ruler elected by the people, with a sense of authority and responsibility.
		{"Chancellor", "Chancellor"},         // For a ruler who leads a council or governing body, with a sense of leadership and guidance.
		{"Prime Minister", "Prime Minister"}, // For a ruler who serves as the head of government, with a focus on administration and policy.
		{"Speaker", "Speaker"},               // For a ruler who represents the interests of the people, with a focus on communication and diplomacy.
	},
}

var republicTitles = map[TitleTag][][2]string{
	TitleTagOther: {
		{"Consul", "Consul"},           // For a ruler who represents the interests of the people, with a focus on diplomacy and negotiation.
		{"Senator", "Senatrix"},        // For a ruler who serves in a legislative body, with a focus on lawmaking and governance.
		{"Tribune", "Tribuness"},       // For a ruler who protects the rights and interests of the people, with a focus on justice and equality.
		{"Praetor", "Praetor"},         // For a ruler who administers justice and maintains order, with a focus on law and order.
		{"Aedile", "Aedilissa"},        // For a ruler who oversees public works and infrastructure, with a focus on development and progress.
		{"Quaestor", "Quaestor"},       // For a ruler who manages the finances and resources of the state, with a focus on economy and trade.
		{"Censor", "Censor"},           // For a ruler who enforces moral and ethical standards, with a focus on integrity and virtue.
		{"Magistrate", "Magistratess"}, // For a ruler who governs a city or region, with a focus on administration and governance.
		{"Legate", "Legatress"},        // For a ruler who acts as an ambassador or envoy, with a focus on diplomacy and negotiation.
		{"Praefect", "Praefecta"},      // For a ruler who commands the military forces of the state, with a focus on defense and security.
	},
}

var dictatorialTitles = map[TitleTag][][2]string{
	TitleTagMilitary: {
		{"Dictator", "Dictatrix"},                  // For a dictator who demands absolute obedience.
		{"Strategos", "Strategessa"},               // Greek terms for military leader, emphasizing tactical prowess.
		{"Supreme Commander", "Supreme Commander"}, // For a dictator who leads through military might and conquest.
		{"Warlord", "Warmistress"},                 // For a dictator who rules through martial strength and conquest.
		{"Warmaster", "Warmistress"},               // For a dictator who rules through martial strength and conquest.
	},
	TitleTagPower: {
		{"Dominus", "Domina"},        // A Latin term meaning master or lord / mistress or lady, ideal for a dictator who demands absolute obedience.
		{"Imperator", "Imperatrix"},  // The Latin word for emperor / empress, offering a strong and historical feel.
		{"Imperator", "Imperatrice"}, // The French/Latin word for emperor / empress, with a touch of elegance and power.
		{"Kaiser", "Kaiserin"},       // The German word for emperor / empress, with a sense of authority and command.
		{"Tsar", "Tsarina"},          // The Russian word for emperor / empress, with a sense of autocratic power.
		{"Czar", "Czarina"},          // The Russian term for emperor / empress, with a sense of autocratic power.
		{"Emperor", "Empress"},       // Classic titles for a supreme ruler, with a sense of grandeur and authority.
		{"Basileus", "Basilissa"},    // Greek terms for king / queen, with a touch of ancient power and wisdom.
		{"Konung", "Drottning"},      // Old Norse words for king / queen, perfect for a harsh and brutal ruler.
		{"Rex", "Regina"},            // Latin terms for king / queen, with a sense of ancient authority and power.
		{"Regent", "Regentess"},      // While traditionally a temporary ruler, it can be twisted for a dictator who seized power and refuses to relinquish it.
	},
	TitleTagService: serviceTitles, // TODO: maybe use a prefix to distinguish from the council titles.
	TitleTagWisdom: {
		{"Father", "Mother"},       // For a dictator who sees themselves as a nurturing and protective figure.
		{"Patriarch", "Matriarch"}, // For a ruler who emphasizes tradition and family control.
	},
	TitleTagDivine: {
		{"High Priest", "High Priestess"}, // For a dictator who claims divine guidance or spiritual authority.
		{"Oracle", "Oracle"},              // For a ruler believed to have prophetic abilities.
		{"Prophet", "Prophetess"},         // For a dictator who claims to speak for the gods or foresee the future.
	},
	TitleTagSensuality: {
		{"Seducer", "Seductress"},          // For a dictator who rules through charm and manipulation.
		{"Enchanter", "Enchantress"},       // For a ruler who uses magic or mystique to control others.
		{"Tempter", "Temptress"},           // For a dictator who lures others into their power through temptation.
		{"Whisperer", "Whisperess"},        // For a ruler who manipulates events from the shadows, using whispers and secrets.
		{"Shadowmaster", "Shadowmistress"}, // For a dictator who rules from the shadows, unseen and mysterious.
		{"Trickster", "Trickstress"},       // For a dictator who rules through deception and cunning.
	},
	TitleTagGrandeur: {
		{"Archon", "Archessa"},              // A powerful ruler from ancient Greece, implying absolute control.
		{"Sovereign", "Sovereign"},          // A classic term for a supreme ruler, it takes on a more mysterious air in a fantasy world.
		{"Overlord", "Overlady"},            // Evokes a sense of dominion and control over a vast territory.
		{"Autarch", "Auctrix"},              // An independent ruler with absolute authority.
		{"Hegemon", "Hegemone"},             // A leader of a group of states, implying dominance.
		{"Suzerain", "Suzeraine"},           // A feudal lord with control over lesser lords.
		{"Regent", "Regentess"},             // While traditionally a temporary ruler, it can be twisted for a dictator who seized power and refuses to relinquish it.
		{"Grand Vizier", "Grand Vizieress"}, // For a dictator who rules with cunning and intrigue.
	},
	TitleTagMagic: {
		{"Archmage", "Archmagess"},      // For a dictator who wields powerful magic.
		{"Wizard-King", "Wizard-Queen"}, // For a ruler who commands mystical forces.
		{"Warlock", "Witch-Queen"},      // For a dictator who uses dark magic to control their subjects.
	},
}

var chiefdomTitles = map[TitleTag][][2]string{
	TitleTagMilitary: {
		{"Chieftain", "Chieftainess"}, // For a chief who leads through military might and conquest.
		{"Warlord", "Warmistress"},    // For a chief who rules through martial strength and conquest.
		{"Warmaster", "Warmistress"},  // For a chief who rules through martial strength and conquest.
	},
	TitleTagOther: {
		{"Chief", "Chieftess"},     // For a chief who leads their people with wisdom and strength.
		{"Elder", "Elderess"},      // For a chief who is respected for their age and experience.
		{"Patriarch", "Matriarch"}, // For a chief who emphasizes tradition and family control.
		{"Pater", "Mater"},         // For a chief who sees themselves as a nurturing and protective figure.
		{"Sire", "Dame"},           // For a chief who embodies the ideals of chivalry and honor.
		{"Guide", "Guidess"},       // For a chief who leads their people with wisdom and insight.
	},
}

// There could be prefixes or suffixes that are added to the title.
// Latin suffixes:
var latinSuffixes = [][2]string{
	{"Ultimus", "Ultima"},                // The last or final ruler of a dynasty or empire.
	{"Primus", "Prima"},                  // The first ruler of a new dynasty or empire.
	{"Maximus", "Maxima"},                // The greatest or most powerful ruler of all time.
	{"Invictus", "Invicta"},              // The unconquered ruler, who has never been defeated.
	{"Magnus", "Magna"},                  // The great ruler, known for their wisdom and power.
	{"Terribilis", "Terribilissa"},       // The fearsome ruler, who strikes terror into the hearts of their subjects.
	{"Gloriosus", "Gloriosa"},            // The glorious ruler, who is celebrated and revered by all.
	{"Infelix", "Infelix"},               // The unlucky ruler, who brings misfortune and disaster to their realm.
	{"Audax", "Audax"},                   // The bold ruler, who takes risks and challenges the status quo.
	{"Fortis", "Fortis"},                 // The strong ruler, who is known for their physical and mental strength.
	{"Ferox", "Ferocis"},                 // The fierce ruler, who rules with aggression and violence.
	{"Pius", "Pia"},                      // The pious ruler, who is devoted to the gods and their people.
	{"Felix", "Felicis"},                 // The lucky ruler, who brings prosperity and good fortune to their realm.
	{"Aeternus", "Aeterna"},              // The eternal ruler, who is said to live forever and never age.
	{"Immortalis", "Immortalis"},         // The immortal ruler, who is said to be a god or demigod.
	{"Divinus", "Divina"},                // The divine ruler, who is believed to be a god or chosen by the gods.
	{"Sacrosanctus", "Sacrosancta"},      // The sacred ruler, who is revered as a holy figure by their people.
	{"Regalis", "Regalis"},               // The royal ruler, who embodies the ideals of kingship and nobility.
	{"Nobilis", "Nobilis"},               // The noble ruler, who is known for their honor and integrity.
	{"Justus", "Justa"},                  // The just ruler, who rules with fairness and impartiality.
	{"Magnanimus", "Magnanima"},          // The magnanimous ruler, who is generous and forgiving to their subjects.
	{"Victor", "Victrix"},                // The victorious ruler, who has won great battles and conquered their enemies.
	{"Triumphator", "Triumphatrix"},      // The triumphant ruler, who is celebrated for their victories and achievements.
	{"Triumphantalis", "Triumphantalis"}, // The triumphant ruler, who is celebrated for their victories and achievements.
}

var serviceTitles = [][2]string{
	{"Protector", "Protectress"},     // For a council leader who sees themselves as a guardian of the people.
	{"Guardian", "Guardian"},         // A simple and straightforward title for a ruler who protects their subjects.
	{"Steward", "Stewardess"},        // For a ruler who manages the affairs of the realm with care and responsibility.
	{"Warden", "Wardess"},            // For a ruler who oversees the safety and security of the land.
	{"Defender", "Defendress"},       // For a ruler who protects their people from external threats.
	{"Sentinel", "Sentinella"},       // For a ruler who stands watch over the realm, protecting it from danger.
	{"Shieldbearer", "Shieldmaiden"}, // For a ruler who defends their people with strength and courage.
	{"Paladin", "Paladiness"},        // For a ruler who upholds justice and righteousness, protecting the weak and innocent.
	{"Champion", "Championess"},      // For a ruler who fights for their people and leads them to victory.
	{"Hero", "Heroine"},              // For a ruler who is celebrated for their bravery and valor in battle.
	{"Knight", "Dame"},               // For a ruler who embodies the ideals of chivalry and honor.
}

// ...
// Like for example:
// - Domina Invicta: The unconquered lady, who has never been defeated.
// - Imperator Magnus: The great emperor, known for their wisdom and power.
// - King Maximus: The greatest or most powerful king of all time.
// - Queen Gloriosa: The glorious queen, who is celebrated and revered by all.
// - Archon Terribilis: The fearsome ruler, who strikes terror into the hearts of their subjects.
//
// In English, we could use:
// - The Great: For a ruler who is known for their wisdom and power.
// - The Conqueror: For a ruler who has won great battles and conquered their enemies.
// - The Wise: For a ruler who is known for their intelligence and insight.
// - The Just: For a ruler who rules with fairness and impartiality.
// - The Magnificent: For a ruler who is celebrated for their victories and achievements.
// - The Terrible: For a ruler who strikes terror into the hearts of their subjects.
// - The Bold: For a ruler who takes risks and challenges the status quo.
// - The Fierce: For a ruler who rules with aggression and violence.
// - The Pious: For a ruler who is devoted to the gods and their people.
// - The unbroken: For a ruler who has never been defeated.
// - The Eternal: For a ruler who is said to live forever and never age.
// Note that with latin we would add the suffixes to the title, while in English we would add them as a prefix
// or as a suffix to the name.

// Archmage: For a dictator who wields powerful magic.
// Dragon-Speaker: For a ruler who commands dragons or other mythical creatures.
// Herald of the Apocalypse: For a truly terrifying dictator who claims to usher in a new age.
// Oracle-King/Queen: For a dictator who claims divine guidance or prophetic abilities.

// Titles with a sense of grandeur:
// Grand Vizier: For a dictator who rules with cunning and intrigue.
// High Inquisitor: For a dictator who enforces strict religious or ideological conformity.
// Hierophant: A religious leader who interprets sacred mysteries, for a dictator with a claim to divine authority.

// Titles with a sense of mystery:
// Shadow Lord/Lady: For a dictator who rules from the shadows, manipulating events behind the scenes.
// Veiled Sovereign: For a dictator who conceals their true identity or motives.
// Whispering Tyrant: For a dictator who rules through fear and whispers of dark power.

// Title returns the title of the leader.
// TODO: Role as a parameter, so we can get the title of a council member, etc.
// TODO: Also separate by language, so we can always get the titles for all roles
// with a consistent theme and feel.
func (l LeadershipForm) Title(preferredTags ...TitleTag) [2]string {
	// Pick a list of possible titles based on the leadership form.
	var titles map[TitleTag][][2]string
	switch l {
	case LeadershipFormMonarchy:
		titles = monarchyTitles
	case LeadershipFormCouncil:
		titles = councilTitles
	case LeadershipFormRepublic:
		titles = republicTitles
	case LeadershipFormDemocracy:
		titles = democracyTitles
	case LeadershipFormDictatorship:
		titles = dictatorialTitles
	case LeadershipFormChiefdom:
		titles = chiefdomTitles
	default:
		return [2]string{"Leader", "Leader"}
	}

	// Make a map of preferred tags.
	// We will add the preferred tags twice, so that they have a higher chance of being picked.
	// TODO: Improve this.
	preferred := make(map[TitleTag]bool)
	for _, tag := range preferredTags {
		preferred[tag] = true
	}

	// Pick a category randomly.
	var category TitleTag
	var possible []TitleTag
	for k := range titles {
		possible = append(possible, k)
		if preferred[k] {
			possible = append(possible, k)
		}
	}
	category = possible[rand.Intn(len(possible))]
	// Pick a title randomly from the category.
	title := titles[category][rand.Intn(len(titles[category]))]
	return title
}

type SuccessionType int

const (
	SuccessionTypeInherited SuccessionType = iota
	SuccessionTypeElected
	SuccessionTypeAppointed
	SuccessionTypeMerit
	SuccessionTypeChosen
	SuccessionTypeMax
)

func (l LeadershipForm) GetSuccessionType() SuccessionType {
	switch l {
	case LeadershipFormMonarchy:
		return SuccessionTypeInherited
	case LeadershipFormCouncil:
		return SuccessionTypeElected
	case LeadershipFormRepublic:
		return SuccessionTypeElected
	case LeadershipFormDemocracy:
		return SuccessionTypeElected
	case LeadershipFormDictatorship:
		return SuccessionTypeAppointed
	case LeadershipFormChiefdom:
		return SuccessionTypeChosen
	}
	return SuccessionTypeChosen
}

type GenderRestriction int

const (
	GenderRestrictionNone GenderRestriction = iota
	GenderRestrictionMale
	GenderRestrictionFemale
	// GenderRestrictionMalePreferred
	// GenderRestrictionFemalePreferred
	GenderRestrictionMax
)

func (l LeadershipForm) GetGenderRestriction() GenderRestriction {
	return GenderRestriction(rand.Intn(int(GenderRestrictionMax)))
}

func (l LeadershipForm) Generate(preferredTags ...TitleTag) *Leadership {
	return &Leadership{
		Form:              l,
		Title:             l.Title(preferredTags...),
		Succession:        l.GetSuccessionType(),
		GenderRestriction: l.GetGenderRestriction(),
	}
}

type Leadership struct {
	Form              LeadershipForm    // The form of leadership.
	Title             [2]string         // Primary leader title (male and female variants)
	Succession        SuccessionType    // The type of succession.
	GenderRestriction GenderRestriction // Restrictions on gender for leadership (if any).
}

// ChangeForm changes the form of leadership.
func (l *Leadership) ChangeForm(form LeadershipForm, preferredTags ...TitleTag) {
	l.Form = form
	l.Title = form.Title(preferredTags...)
	l.Succession = form.GetSuccessionType()
	// l.GenderRestriction = form.GetGenderRestriction()
}

// Genders
// TODO: Expand
const (
	GenderMale = iota
	GenderFemale
)

// GetTitle returns the title of the leader given the gender.
func (l *Leadership) GetTitle(gender int) string {
	return l.Title[gender]
}
