package genlanguage

const (
	GenBaseNumber      = "number"
	GenBaseBeing       = "being"
	GenBaseAnimal      = "animal"
	GenBaseAdjective   = "adjective"
	GenBaseGenitive    = "genitive"
	GenBaseTheGenitive = "theGenitive"
	GenBaseColor       = "color"
)

// GenBase contains a map of word lists used for name generation.
// TODO: Group individual entries by logical categories.
// So we can build up a pantheon of gods, each associated with different domains.
// For example:
// {North, South, East, West} -> {Direction}
// {Bride, Groom, Widow, Widower, Wife, Husband} -> {Marriage}
// {Giver, Taker, Destroyer, Creator, Maker, Breaker} -> {Action}
// {Sky, Earth, Water, Fire, Air, Spirit} -> {Elements}
// {Light, Dark, Bright, Shining, Shadow, Darkness} -> {Light}
var GenBase = map[string][]string{
	GenBaseNumber: {"Zero", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Eleven", "Twelve"},
	GenBaseBeing: {
		"Ancestor",
		"Ancient",
		"Angel",
		"Betrayer",
		"Bride",
		"Brother",
		"Chief",
		"Child",
		"Council",
		"Creator",
		"Deity",
		"Elder",
		"Father",
		"Forebearer",
		"Forefather",
		"Foremother",
		"Forgiven",
		"Forgotten",
		"Giver",
		"God",
		"Goddess",
		"Golem",
		"Groom",
		"Guardian",
		"Guide",
		"Keeper",
		"King",
		"Knight",
		"Lady",
		"Lord",
		"Lover",
		"Maker",
		"Master",
		"Mistress",
		"Mother",
		"Numen",
		"Orphan",
		"Overlord",
		"Reaper",
		"Ruler",
		"Seducer",
		"Seductress",
		"Servant",
		"Sister",
		"Spirit",
		"Thief",
		"Virgin",
		"Warrior",
		"Watcher",
		"Widow",
		"Widower",
		"Wife",
		"Witch",
		"Wizard",
	},
	GenBaseAnimal: {
		"Antelope",
		"Ape",
		"Badger",
		"Basilisk",
		"Bear",
		"Beaver",
		"Bison",
		"Boar",
		"Buffalo",
		"Camel",
		"Cat",
		"Centaur",
		"Chimera",
		"Cobra",
		"Crane",
		"Crocodile",
		"Crow",
		"Cyclope",
		"Deer",
		"Dog",
		"Dragon",
		"Eagle",
		"Elk",
		"Falcon",
		"Fox",
		"Goat",
		"Goose",
		"Hare",
		"Hawk",
		"Heron",
		"Horse",
		"Hound",
		"Hyena",
		"Ibis",
		"Jackal",
		"Jaguar",
		"Kraken",
		"Lark",
		"Leopard",
		"Lion",
		"Mantis",
		"Marten",
		"Moose",
		"Mule",
		"Narwhal",
		"Owl",
		"Ox",
		"Panther",
		"Pegasus",
		"Phoenix",
		"Rat",
		"Raven",
		"Rook",
		"Scorpion",
		"Serpent",
		"Shark",
		"Sheep",
		"Snake",
		"Sphinx",
		"Spider",
		"Swan",
		"Tiger",
		"Turtle",
		"Unicorn",
		"Viper",
		"Vulture",
		"Walrus",
		"Wolf",
		"Wolverine",
		"Worm",
		"Wyvern",
	},
	GenBaseAdjective: {
		"Aggressive",
		"Almighty",
		"Ancient",
		"Angry",
		"Anxious",
		"Awful",
		"Beautiful",
		"Benevolent",
		"Big",
		"Blind",
		"Blond",
		"Bloody",
		"Brave",
		"Broken",
		"Brutal",
		"Burning",
		"Calm",
		"Careful",
		"Charming",
		"Cheerful",
		"Chosen",
		"Clean",
		"Crazy",
		"Cruel",
		"Dead",
		"Deadly",
		"Deaf",
		"Deathless",
		"Deep",
		"Defiant",
		"Delicate",
		"Delightful",
		"Desperate",
		"Devastating",
		"Distant",
		"Disturbing",
		"Divine",
		"Dying",
		"Enchanting",
		"Ephemeral",
		"Eternal",
		"Evil",
		"Explicit",
		"Fair",
		"Far",
		"Fat",
		"Fatal",
		"Favorable",
		"Flying",
		"Friendly",
		"Frozen",
		"Giant",
		"Good",
		"Grateful",
		"Great",
		"Greedy",
		"Happy",
		"High",
		"Holy",
		"Honest",
		"Huge",
		"Hungry",
		"Immutable",
		"Infallible",
		"Inherent",
		"Last",
		"Latter",
		"Lonely",
		"Lost",
		"Loud",
		"Loving",
		"Lucky",
		"Lustful",
		"Mad",
		"Magical",
		"Main",
		"Major",
		"Marine",
		"Naval",
		"New",
		"Old",
		"Patient",
		"Peaceful",
		"Pregnant",
		"Prime",
		"Proud",
		"Pure",
		"Sacred",
		"Sad",
		"Scary",
		"Secret",
		"Selected",
		"Severe",
		"Silent",
		"Sleeping",
		"Slumbering",
		"Spiteful",
		"Strong",
		"Sunny",
		"Superior",
		"Sustainable",
		"Troubled",
		"Undying",
		"Unhappy",
		"Unknown",
		"Waking",
		"Wild",
		"Wise",
		"Worried",
		"Young",
	},
	GenBaseGenitive: {
		"Cold",
		"Darkness",
		"Dawn",
		"Day",
		"Death",
		"Doom",
		"Dreams",
		"Dusk",
		"Fate",
		"Fire",
		"Fog",
		"Frost",
		"Gates",
		"Heaven",
		"Home",
		"Hope",
		"Ice",
		"Justice",
		"Kings",
		"Life",
		"Light",
		"Lightning",
		"Love",
		"Nature",
		"Night",
		"Pain",
		"Snow",
		"Springs",
		"Stars",
		"Summer",
		"Sun",
		"Sunset",
		"Thunder",
		"Time",
		"Victory",
		"War",
		"Wealth",
		"Winter",
		"Wisdom",
	},
	GenBaseTheGenitive: {
		"Abyss",
		"Blood",
		"Dawn",
		"Earth",
		"East",
		"Eclipse",
		"Fall",
		"Harvest",
		"Moon",
		"North",
		"Peak",
		"Rainbow",
		"Sea",
		"Sky",
		"South",
		"Stars",
		"Storm",
		"Sun",
		"Tree",
		"Underworld",
		"West",
		"Wild",
		"Word",
		"World",
	},
	GenBaseColor: {
		"Amber",
		"Black",
		"Blue",
		"Bright",
		"Brown",
		"Dark",
		"Golden",
		"Green",
		"Grey",
		"Light",
		"Orange",
		"Pink",
		"Purple",
		"Red",
		"White",
		"Yellow",
	},
}