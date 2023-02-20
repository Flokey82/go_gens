package genstory

import (
	"github.com/Flokey82/go_gens/genlanguage"
)

// NewSimpleTitleConfig returns a simple title configuration.
func NewSimpleTitleConfig(templates []string) *TextConfig {
	return &TextConfig{
		TokenPools:       DefaultTitleTokenPool,
		TokenIsMandatory: DefaultTitleTokenIsMandatory,
		Tokens:           DefaultTitleTokens,
		Templates:        templates,
		Title:            true,
	}
}

const (
	// Examples: The Tower, The Tree, Bronzemurdered, Likot Ubendeb, Animal Behaviours
	// Notably, these may or may not be plural and/or have an article. Name is directly related to the content of the book.
	TokenName = "[NAME]"
	// Examples: Tower, Tree, Bronzemurdered, Likot Ubendeb, Animal Behaviours
	// Notice that the inserted text may still be plural - limiting its usage. These likewise are related to the content of the book.
	// Examples: Despair, Roots, Scrolls, Wheel-and-axels
	TokenNoArtName = "[NO_ART_NAME]"
	// Examples: Despair, Roots, Scrolls, Wheel-and-axels
	// Again, may or may not be plural. These seem to have very little correlation to the book's topic. If it has any relationship to the content of the book, it could be words within the author's entity's vocabulary.
	TokenNoun = "[NOUN]"
	// Also seem to have little correlation to the book's topic. They are currently guessed to be from the civ, if not purely random.
	// Examples: Boyish, Inky, Angry, Bronzed
	TokenAdj = "[ADJ]"
	// Examples: The Age Of Legends, The Age Of Hill Titan and Dragon
	// These are pulled from the world's history, not from all possibilities.
	TokenAnyAge = "[ANY_AGE]"
	// Examples: He, She, We, They
	// It is believed it only generates subject pronouns (not "Us" or "Them"). Past tense makes this easy to use.
	TokenAnyPronoun = "[ANY_PRONOUN]"
	// Examples: His, Her, Our, Their
	TokenAnyPossessivePronoun = "[POSS_PRONOUN]"
	// Examples: The Fool Laughs, The Day Can Say It In The End, It Foretells Afterwards, The Day Mourns
	// Due to the wide variety, this is pretty hard to use.
	TokenPhrase = "[PHRASE]"
	// Examples: Riverwood, Paris, Eridu, The Great Library
	// Place names are pulled from the world's history, not from all possibilities.
	TokenPlace = "[PLACE]"
	// Examples: Exploring, discussing, understanding, learning, discovering,
	TokenVerb = "[VERB]"
	// Examples: Death, life, fire, fate...
	TokenGenitiveNoun = "[GENITIVE_NOUN]"
	// Examples: children, adults, toddlers, elders, babies, teenagers, youths, infants, adolescents, juve
	TokenReaderAge = "[READER_AGE]"
)

// DefaultTitleTokens is a list of all default title tokens.
var DefaultTitleTokens = []string{
	TokenName,
	TokenNoArtName,
	TokenNoun,
	TokenAdj,
	TokenAnyAge,
	TokenAnyPronoun,
	TokenAnyPossessivePronoun,
	TokenPhrase,
	TokenPlace,
	TokenVerb,
	TokenGenitiveNoun,
	TokenReaderAge,
}

// DefaultTitleTokenIsMandatory is a map of title tokens to whether or not they need
// to be provided for generating the title. Optional tokens are replaced with
// a random representation.
// Example: "[NAME] and the [NOUN]" will always have to be provided, but
// [ADJ] might be chosen randomly from a list of adjectives.
var DefaultTitleTokenIsMandatory = map[string]bool{
	TokenName:      true,
	TokenNoArtName: true,
	TokenNoun:      true,
	TokenAnyAge:    true,
}

// DefaultTitleTokenPool is a map of title tokens to a list of possible replacements.
var DefaultTitleTokenPool = map[string][]string{
	TokenNoun:                 BookTitleNouns,
	TokenAdj:                  BookTitleAdjectives,
	TokenAnyAge:               BookTitleAge,
	TokenAnyPronoun:           BookTitlePronouns,
	TokenAnyPossessivePronoun: BookTitlePossessivePronouns,
	TokenPhrase:               BookTitlePhrases,
	TokenPlace:                BookTitlePlaces,
	TokenVerb:                 BookTitleVerbs,
	TokenGenitiveNoun:         BookTitleGenitiveNouns,
	TokenReaderAge:            BookTitleReaderAge,
}

var BookTitlePronouns = []string{
	"he",
	"she",
	"it",
	"they",
}

var BookTitlePossessivePronouns = []string{
	"his",
	"her",
	"its",
	"their",
}

var BookTitleGenitiveNouns = genlanguage.GenBase[genlanguage.GenBaseGenitive]

var BookTitleAdjectives = genlanguage.GenBase[genlanguage.GenBaseAdjective]

var BookTitleNouns = []string{
	"death",
	"life",
	"fire",
	"fate",
	"love",
	"hope",
	"despair",
	"roots",
	"scrolls",
	"wheel-and-axels",
	"nuts-and-bolts",
}

var BookTitleVerbs = []string{
	"exploring",
	"discussing",
	"understanding",
	"learning",
	"discovering",
	"searching",
	"seeking",
	"finding",
	"creating",
}

var BookTitlePhrases = []string{
	"the fool laughs",
	"the day can say it in the end",
	"it smells afterwards",
	"the day mourns",
	"all and none",
	"the day is born",
	"been there, done that",
	"funny enough for a tragedy",
}

var BookTitleHandbookSubtitles = []string{
	"the ultimate guide",
	"the complete guide",
	"the definitive guide",
	"the ultimate handbook",
	"the complete handbook",
	"the definitive handbook",
	"the ultimate reference",
	"the complete reference",
	"the definitive reference",
	"the ultimate manual",
	"the complete manual",
	"the definitive manual",
	"the ultimate book",
	"the complete book",
	"the definitive book",
}

var BookTitleReaderAge = []string{
	"children",
	"adults",
	"beginners",
	"experts",
	"everyone",
	"kids",
	"teens",
	"toddlers",
	"elders",
	"babies",
	"teenagers",
	"youths",
	"infants",
	"adolescents",
	"juveniles",
	"youngsters",
}

var BookTitleAge = []string{
	"the age of steam",
	"the age of the forgotten",
	"the age of the gods",
	"the age of the machine",
}

var BookTitlePlaces = []string{
	"cities",
	"villages",
	"towns",
	"hamlets",
	"city-states",
	"kingdoms",
	"empires",
	"nations",
	"countrys",
	"continents",
	"worlds",
	"universes",
	"pits",
}

// NOTE: This is based on "book_instruction.txt" from Dwarf Fortress.
var BookInstructionTitles = []string{
	"[NAME]",
	"A Course on [NAME]",
	"A Humble Offering to [NAME]",
	"A Meditation on [NAME]",
	"A Record of [NAME]",
	"A Treatise on [NAME]",
	"A World of [NAME]",
	"About [NAME]",
	"After [NAME]",
	"Against [NAME]",
	"An Exploration of [NAME]",
	"An Introduction to [NAME]",
	"An Offering to [NAME]",
	"At One With [NAME]",
	"Before [NAME]",
	"Better [NO_ART_NAME]",
	"Beyond [NAME]",
	"Book of [NAME]",
	"Captivated by [NAME]",
	"Choose [NAME]",
	"Classic [NO_ART_NAME]",
	"Commentary on [NAME]",
	"Common Sense [NO_ART_NAME]",
	"Concerning [NAME]",
	"Discourse on [NAME]",
	"Discovering [NAME]",
	"Doubts About [NAME]",
	"Dreams of [NAME]",
	"Elements of [NAME]",
	"Errors In [NAME]",
	"Exploring [NAME]",
	"Explorations of [NAME]",
	"Factual [NO_ART_NAME]",
	"For the Love of [NAME]",
	"Foundations of [NAME]",
	"Give Me [NAME]",
	"Great [NO_ART_NAME]",
	"In Pursuit of [NAME]",
	"Inquiries on [NAME]",
	"Interpretations of [NAME]",
	"Introduction to [NAME]",
	"It All Begins With [NAME]",
	"It Is [NAME]",
	"Journey to [NAME]",
	"Life With [NAME]",
	"Master of [NAME]",
	"Mastering [NAME]",
	"Meditations on [NAME]",
	"Misconceptions About [NAME]",
	"More [NO_ART_NAME]",
	"Musings on [NAME]",
	"My Thoughts on [NAME]",
	"Mysteries of [NAME]",
	"Never Underestimate [NAME]",
	"New [NO_ART_NAME]",
	"Of [NAME]",
	"On [NAME]",
	"Pathways to [NAME]",
	"Principles of [NAME]",
	"Question [NAME]",
	"Questions About [NAME]",
	"Records of [NAME]",
	"Reflections on [NAME]",
	"Secret [NO_ART_NAME]",
	"Start Your Day With [NAME]",
	"Strange [NO_ART_NAME]",
	"The Art of [NAME]",
	"The Book of [NAME]",
	"The Future of [NAME]",
	"The Great [NO_ART_NAME]",
	"The Hidden Meaning of [NAME]",
	"The History of [NAME]",
	"The Interpretation of [NAME]",
	"The Journey into [NAME]",
	"The Knowledge of [NAME]",
	"The Meaning of [NAME]",
	"The Mysteries of [NAME]",
	"The Mystery of [NAME]",
	"The Nuanced [NO_ART_NAME]",
	"The Possibilities of [NAME]",
	"The Pursuit of [NAME]",
	"The Secret of [NAME]",
	"The Student's [NO_ART_NAME]",
	"The Study of [NAME]",
	"The True [NO_ART_NAME]",
	"The Truth About [NAME]",
	"The Unabridged [NO_ART_NAME]",
	"The Way With [NAME]",
	"The Wizard's Guide to [NAME]",
	"The World of [NAME]",
	"The World Without [NAME]",
	"Thoughts on [NAME]",
	"Time Spent With [NAME]",
	"Traditional [NO_ART_NAME]",
	"Treatise on [NAME]",
	"True [NO_ART_NAME]",
	"Uncanny [NO_ART_NAME]",
	"Uncovering [NAME]",
	"Understanding [NAME]",
	"Unknown [NO_ART_NAME]",
	"Unusual [NO_ART_NAME]",
	"Useful [NO_ART_NAME]",
	"Victory By [NAME]",
	"[NAME] After The End",
	"[NAME] and Beyond",
	"[NAME] and More",
	"[NAME] and Other Topics",
	"[NAME] and Other Travesties",
	"[NAME] and The Coming Troubles",
	"[NAME] and The Universe",
	"[NAME] Came Full Circle",
	"[NAME] Explained",
	"[NAME] Exposed",
	"[NAME] For Everyone",
	"[NAME] For Students",
	"[NAME] For The Beginning Practitioner",
	"[NAME] In Practice",
	"[NAME] In The Modern Era",
	"[NAME] In The Time of My Ancestors",
	"[NAME] In Theory",
	"[NAME] In [ANY_AGE]",
	"[NAME] Interpreted",
	"[NAME] Might Help",
	"[NAME] Questioned",
	"[NAME] The Easy Way",
	"[NAME] Uncovered",
	"[NAME] Understood",
	"[NAME] When It Counts",
	"[NAME] Within Reason",
	"[NAME] Without Limits",
	"[NAME], Abridged",
	"[NAME], My Life",
	"[NAME], My Love",
	"[NAME]: A Brief History",
	"[NAME]: A Brief Introduction",
	"[NAME]: A New Approach",
	"[NAME]: A Quandary",
	"[NAME]: Before and After",
	"[NAME]: Common Practice",
	"[NAME]: Fact or Fiction?",
	"[NAME]: Further Musings",
	"[NAME]: My Only Mistake",
	"[NAME]: Natural or Supernatural?",
	"[NAME]: Principles and Practice",
	"[NAME]: Problems And Solutions",
	"[NAME]: The Definitive Guide",
	"[NAME]: The Truth",
	"Can [NAME] Save The World?",
	"Could It Be [NAME]?",
	"Did [NAME] Falter?",
	"Do We Understand [NAME]?",
	"Do You Know [NAME]?",
	"First [NAME], Then The World!",
	"The [NO_ART_NAME] Book",
	"To [NAME] and Glory!",
}

// NOTE: This is based on "book_art.txt" from Dwarf Fortress.
var BookArtTitles = []string{
	"[NAME]",
	"[PHRASE]",
	"[ADJ] [NO_ART_NAME]",
	"The [ADJ] [NO_ART_NAME]",
	"[NAME] [ADJ]",
	"[NAME] and the [NOUN]",
	"[NAME] and the [ADJ] [NOUN]",
	"The [NOUN] and [NAME]",
	"The [ADJ] [NOUN] and [NAME]",
	"[NAME]: [PHRASE]",
	"It Must Have Been [NAME]",
	"My Friend [NAME]",
	"The Birth of [NAME]",
	"The Sun Sets on [NAME]",
	"We See [NAME]",
	"[NAME] Ever Onward",
	"[NAME] and Nothing More",
	"And [ANY_PRONOUN] Sang '[NAME]!'",
}

// BookVariantTitles is a list of very flexible book titles.
var BookVariantTitles = []string{
	"[NAME]",
	"[NAME]: [PHRASE]",
	"[NAME] and [NOUN]",
	"[NAME] and the [ADJ] [NOUN]",
	"[NOUN] and [NAME]",
	"The [ADJ] [NOUN] and [NAME]",
	"The [ADJ] [NOUN] of [NAME]",
	"[ADJ] [PLACE]",
	"[ADJ] [PLACE] and [NAME]",
	"[ADJ] [PLACE] and the [NOUN]",
	"[ADJ] [PLACE] and the [ADJ] [NOUN]",
	"[ADJ] [PLACE] of the [NOUN]",
	"[ADJ] [PLACE] of the [ADJ] [NOUN]",
	"[ADJ] [NOUN] of [PLACE]",
	"[VERB] [NAME]",
	"[VERB] [PLACE] and [NAME]",
	"[VERB] [PLACE] and the [NOUN]",
	"[VERB] [PLACE] and the [ADJ] [NOUN]",
	"[VERB] [PLACE] of the [NOUN]",
	"[VERB] [PLACE] of the [ADJ] [NOUN]",
	"[VERB] [NOUN] of [PLACE]",
	"[VERB] [NOUN] of [NAME]",
	"[VERB] [NOUN] of [ADJ] [PLACE]",
	"The [GENITIVE_NOUN] of [NAME]",
	"The [GENITIVE_NOUN] of [PLACE]",
	"The [GENITIVE_NOUN] of [ADJ] [PLACE]",
	"The [GENITIVE_NOUN] of [ADJ] [NOUN]",
	"The [GENITIVE_NOUN] of [ADJ] [NOUN] in [PLACE]",
	"The [GENITIVE_NOUN] of [ADJ] [NOUN] with [NAME]",
	"The [GENITIVE_NOUN] of [ADJ] [NOUN] in [ADJ] [PLACE]",
	"The [GENITIVE_NOUN] of [READER_AGE]",
	"[VERB] [NOUN] of [READER_AGE]",
	"[VERB] [NOUN] for [READER_AGE]",
	"[VERB] [NOUN] of [ADJ] [NOUN] for [READER_AGE]",
	"[NOUN] for [READER_AGE]",
	"[READER_AGE] and [NOUN]",
	"[POSS_PRONOUN] [NOUN]",
	"[POSS_PRONOUN] [ADJ] [NOUN]",
}
