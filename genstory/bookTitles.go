package genstory

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/Flokey82/go_gens/genlanguage"
)

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
	"The [GENITIVE_NOUN] of [ADJ] [NOUN] of [NAME]",
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

type TokenReplacement struct {
	Token       string // The token to replace.
	Replacement string // The replacement text.
}

// GenerateTitle generates a title from the provided tokens and a list of
// possible titles.
//   - The provided tokens are used to replace the tokens in the possible titles.
//   - If a token is not provided, it is replaced with a random value.
//   - If a token is not provided and is not optional, all possible titles are
//     excluded that require that token.
//
// TODO: Also return the selected title template, and the individual replacements,
// so that the caller can use them for the description or the generation of content.
func GenerateTitle(provided []TokenReplacement, titleOptions []string) (string, error) {
	// Count how many token replacements we have for each token.
	tokenReplacements := map[string]int{}
	for _, replacement := range provided {
		tokenReplacements[replacement.Token]++
	}

	// Loop over all titles and find the ones where we have all required tokens.
	possibleTitles := []string{}
	for _, i := range rand.Perm(len(titleOptions)) {
		title := titleOptions[i]
		// Check if we have all required tokens the required number of times.
		var missingToken bool
		for _, token := range titleTokens {
			if tokenReplacements[token] < strings.Count(title, token) {
				if titleTokenIsMandatory[token] {
					missingToken = true
					break
				}
			}
		}

		// Something is missing, skip this title.
		if missingToken {
			continue
		}

		// Also make sure all tokens we have provided are available in the title,
		// since we want to pick a complete title, referencing all provided tokens.
		for _, replacement := range provided {
			if strings.Count(title, replacement.Token) < tokenReplacements[replacement.Token] {
				missingToken = true
				break
			}
		}

		// Something is missing, skip this title.
		if missingToken {
			continue
		}

		// We have all required tokens, add the title to the list of possible titles.
		possibleTitles = append(possibleTitles, title)
	}

	// If we have no possible titles, return an error.
	if len(possibleTitles) == 0 {
		return "", errors.New("no possible titles satisfying the provided tokens")
	}

	// Pick a random title.
	title := randArrayString(possibleTitles)

	// Replace all tokens with the provided replacements.
	for _, replacement := range provided {
		title = strings.Replace(title, replacement.Token, replacement.Replacement, 1)
	}

	// Replace all optional tokens with random replacements.
	remainingTokens, err := ExtractTokens(title)
	if err != nil {
		return "", fmt.Errorf("failed to extract tokens from title: %v", err)
	}

	// Relplace each token one by one until we can't find any more.
	for _, token := range remainingTokens {
		if !titleTokenIsMandatory[token] && strings.Contains(title, token) {
			// Pick a random replacement.
			// TODO: What to do if we don't have any replacements for a token?
			replacement := randArrayString(bookTokenToRandom[token])

			// Replace the token.
			title = strings.Replace(title, token, replacement, 1)
		}
	}
	return strings.Title(title), nil
}

func randArrayString(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	return arr[rand.Intn(len(arr))]
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

var titleTokens = []string{
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

// titleTokenIsMandatory is a map of title tokens to whether or not they need
// to be provided for generating the title. Optional tokens are replaced with
// a random representation.
// Example: "[NAME] and the [NOUN]" will always have to be provided, but
// [ADJ] might be chosen randomly from a list of adjectives.
var titleTokenIsMandatory = map[string]bool{
	TokenName:      true,
	TokenNoArtName: true,
	TokenNoun:      true,
	TokenAnyAge:    true,
}

var bookTokenToRandom = map[string][]string{
	TokenNoun:                 bookTitleNouns,
	TokenAdj:                  bookTitleAdjectives,
	TokenAnyAge:               bookTitleAge,
	TokenAnyPronoun:           bookTitlePronouns,
	TokenAnyPossessivePronoun: bookTitlePossessivePronouns,
	TokenPhrase:               bookTitlePhrases,
	TokenPlace:                bookTitlePlaces,
	TokenVerb:                 bookTitleVerbs,
	TokenGenitiveNoun:         bookTitleGenitiveNouns,
	TokenReaderAge:            bookTitleReaderAge,
}

var bookTitlePronouns = []string{
	"he",
	"she",
	"it",
	"they",
}

var bookTitlePossessivePronouns = []string{
	"his",
	"her",
	"its",
	"their",
}

var bookTitleGenitiveNouns = genlanguage.GenBase[genlanguage.GenBaseGenitive]

var bookTitleAdjectives = genlanguage.GenBase[genlanguage.GenBaseAdjective]

var bookTitleNouns = []string{
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

var bookTitleVerbs = []string{
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

var bookTitlePhrases = []string{
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

var bookTitleReaderAge = []string{
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

var bookTitleAge = []string{
	"the age of steam",
	"the age of the forgotten",
	"the age of the gods",
	"the age of the machine",
}

var bookTitlePlaces = []string{
	"city",
	"village",
	"town",
	"hamlet",
	"city-state",
	"kingdom",
	"empire",
	"nation",
	"country",
	"continent",
	"world",
	"universe",
	"pit",
}

// ExtractTokens extracts all tokens from a string.
// A token is a string surrounded by square brackets ('[]').
func ExtractTokens(s string) ([]string, error) {
	// Scan through the string, looking for tokens.
	var tokens []string
	var tokenStart int
	var inToken bool
	for i, c := range s {
		if c == '[' {
			if inToken {
				return nil, fmt.Errorf("unexpected token start at %d", i)
			}
			inToken = true
			tokenStart = i
		} else if c == ']' {
			if !inToken {
				return nil, fmt.Errorf("unexpected token end at %d", i)
			}
			inToken = false
			tokens = append(tokens, s[tokenStart:i+1])
		}
	}
	if inToken {
		return nil, fmt.Errorf("unexpected end of string")
	}
	return tokens, nil
}
