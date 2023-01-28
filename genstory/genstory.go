package genstory

import (
	"math/rand"
	"strings"

	"github.com/Flokey82/go_gens/genlanguage"
	"github.com/Flokey82/go_gens/genreligion"
)

// NewWorld generates a new world creation mythos using the given seed.
func NewWorld(seed int64) string {
	rng := rand.New(rand.NewSource(seed))
	rlgGen := genreligion.NewGenerator(seed)

	// Generate a new language.
	lang := genlanguage.GenLanguage(rng.Int63())

	// Pick a random intro line.
	// "Long ago,"
	intro := intros[rng.Intn(len(intros))]

	// Pick a random world name.
	// "Flubwubbworld"
	worldName := strings.Title(lang.GetWord("world"))

	// Pick a random strategy for creation of the world.
	switch strategies[rng.Intn(len(strategies))] {
	case StratCreationGod:
		// "was created by flubwubb"
		creationProcess := creation[rng.Intn(len(creation))]
		god := rlgGen.GetDeity(lang, rlgGen.RandDeityGenMethod())
		creationGod := god.FullName()
		return intro + " " + worldName + " was " + creationProcess + " by " + creationGod + "."
	case StratShapingAdjectiveMaterial:
		// "formed from a lone pearl"
		shapingProcess := shaping[rng.Intn(len(shaping))]
		adjective := adjectives[rng.Intn(len(adjectives))]
		material := materials[rng.Intn(len(materials))]
		return intro + " " + worldName + " was " + shapingProcess + " from a " + adjective + " " + material + "."
	case StratShapingMaterialGod:
		// "was shaped from clay by flubwubb"
		shapingProcess := shaping[rng.Intn(len(shaping))]
		material := materials[rng.Intn(len(materials))]
		god := rlgGen.GetDeity(lang, rlgGen.RandDeityGenMethod())
		creationGod := god.FullName()
		return intro + " " + worldName + " was " + shapingProcess + " from a " + material + " by " + creationGod + "."
	}
	return ""
}

// intros contains the intro lines for the world creation mythos.
var intros = []string{
	"Long ago,",
	"As it is written in the ancient texts,",
	"According to the legends,",
	"According to the ancient texts,",
	"In the time before time,",
	"In the beginning,",
	"During the spark of creation,",
}

const (
	StratCreationGod              = "creation + god"
	StratShapingAdjectiveMaterial = "shaping + adjective + material"
	StratShapingMaterialGod       = "shaping + material + god"
)

var strategies = []string{
	StratCreationGod,              // "was created by flubwubb"
	StratShapingAdjectiveMaterial, // "formed from a lone pearl"
	StratShapingMaterialGod,       // "was shaped from clay by flubwubb"
}

var creation = []string{
	"created by",
	"shaped in a dream of",
	"given existence by",
	"brought into being by",
}

var shaping = []string{
	"formed",
	"shaped",
	"created",
	"made",
}

var adjectives = []string{
	"lone",
	"pure",
	"perfect",
	"perfectly round",
	"perfectly square",
	"precious",
	"beautiful",
	"beautifully round",
}

var materials = []string{
	"pearl",
	"gem",
	"crystal",
	"piece of stone",
	"rock",
	"chunk of clay",
	"speck of dust",
	"drop of water",
	"gust of air",
	"mote of fire",
	"ball of earth",
}
