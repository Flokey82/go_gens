package genstory

// NOTE: These categories are based on https://dwarffortresswiki.org/index.php/DF2014:Book
const (
	// This will write down a randomly-selected form of knowledge the adventurer
	// is aware of, to be learned by future readers. Most commonly this will be
	// musical, poetic, and dance forms the adventurer knows or composed.
	// This can also include scientific research the adventurer has learned,
	// and necromancer adventurers can spread the secrets of life and death by
	// writing manuals about them.
	BookTypeManual = iota
	// General writing about a specific site, generally described as "concerning"
	// that town, dark pit, etc. without going into detail.
	BookTypeGuide
	// In-depth writing about a particular site, group, or civilization. This will
	// be presented as multiple chapters, each chapter relating to a historical
	// event related to the writing's subject. It teaches histfigs about the group.
	BookTypeChronicle
	// Generic prose, typically described as having no particular subject.
	BookTypeShortStory
	// Generic prose, typically described as having no particular subject.
	BookTypeNovel
	// Teaches about a person and several events happening to that person, each of
	// which is represented as a separate chapter.
	//
	// Requires Historian's biography topic.
	BookTypeBiography
	// Teaches about the author and several events happening to the author, each of
	// which is represented as a separate chapter.
	//
	// Requires Historian's autobiography topic.
	BookTypeAutobiography
	// Writes a new poem, like the Musical Composition writes a musical composition
	//
	// Requires knowing any poetry forms.
	BookTypePoem
	// Generic prose, typically described as having no particular subject
	BookTypePlay
	// Generic prose, typically described as having no particular subject. These
	// often have no title.
	BookTypeLetter
	// Might be writing about events, people, places or values.
	BookTypeEssay
	// Concerns and teaches a value
	//
	// Requires the Philosopher's dialectic reasoning topic.
	BookTypeDialog
	// Writes new songs. This functions similarly to composing new songs, with the
	// added benefit of writing it down for others to learn. However, unlike normal
	// composition, you do not get to select which musical form to base the song on.
	//
	// Requires knowing any musical forms.
	BookTypeMusicalComposition
	// Writes a new dance, like the Musical Composition.
	//
	// Requires knowing any dance forms.
	BookTypeChoreography
	// Concerns (and teaches about) two historical figures, may emphasize a value too.
	//
	// Requires Historian's comparative biography topic.
	BookTypeComparativeBiography
	// Concerns a list of historical figures
	//
	// Requires Historian's biographic dictionary topic.
	BookTypeBiographicDictionary
	// Concerns the lineage of a specific historical figure. Does not mention anyone
	// besides the main figure.
	//
	// Requires Historian's genealogy topic.
	BookTypeGenealogy
	// Teaches about several notable historical objects in a world, so artifacts,
	// sites, people.
	//
	// Requires Historian's encyclopedia topic.
	BookTypeEncyclopedia
	// Teaches about a culture, and several events happening to that entity.
	//
	// Requires Historian's cultural history topic.
	BookTypeCulturalHistory
	// Teaches about two cultures/groups and may emphasize a value too
	//
	// Requires Historian's cultural comparison topic.
	BookTypeCulturalComparison
	// Teaches an event. Description suggests the work in question is an exploration
	// of what would have happened had this event not played out as it did.
	// May emphasize a value.
	//
	// Requires Historian's alternate history topic.
	BookTypeAlternateHistory
	// Concerns the history of an engineering topic, teaches the topic in question.
	//
	// Requires Historian's treatise on technological advancement topic.
	BookTypeTreatiseOnTechnologicalAdvancement
	// Concerns and teaches about a language.
	//
	// Requires Philosopher's dictionary topic.
	BookTypeDictionary
	// Nothing at the moment, but can be a 'good resource of information' or
	// 'badly compiled'.
	//
	// Requires Astronomer's star chart topic.
	BookTypeStarChart
	// Nothing at the moment, but can be a 'good resource of information' or
	// 'badly compiled'.
	//
	// Requires one of the Astronomer's star catalogue topics.
	BookTypeStarCatalogue
	// This is regarding, and teaches about, a region.
	//
	// Requires Geographer's atlas topic.
	BookTypeAtlas
	// This is a holy book, and teaches about a religion.
	//
	// Requires Cleric's theology topic.
	BookTypeTheology
	// This is a sacred text, and teaches about a religion.
	//
	// Requires Cleric's scripture topic.
	BookTypeScripture
	// This is a legal text, and describes a legal system.
	//
	// Requires Clerks's law topic.
	BookTypeLaw
)
