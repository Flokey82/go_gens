package genstory

import (
	"math/rand"
	"strings"

	"github.com/Flokey82/go_gens/genlanguage"
	"github.com/Flokey82/go_gens/genreligion"
)

// NewWorld generates a new world creation mythos using the given seed.
func NewWorld(seed int64) (string, error) {
	rng := rand.New(rand.NewSource(seed))
	rlgGen := genreligion.NewGenerator(seed)

	// Generate a new language.
	lang := genlanguage.GenLanguage(rng.Int63())

	// Pick a random world name.
	worldName := strings.Title(lang.GetWord("world"))

	// Pick a random strategy for creation of the world.
	tokenReplacements := []TokenReplacement{{Token: StoryTokenWorld, Replacement: worldName}}
	if rand.Intn(2) == 0 {
		// "was created by flubwubb"
		god := rlgGen.GetDeity(lang, rlgGen.RandDeityGenMethod())
		tokenReplacements = append(tokenReplacements, TokenReplacement{Token: StoryTokenGod, Replacement: god.FullName()})
	}
	return StoryConfig.Generate(tokenReplacements)
}

var StoryConfig = &TextConfig{
	TokenPools: map[string][]string{
		StoryTokenIntro:     StoryIntros,
		StoryTokenCreation:  StoryCreationPool,
		StoryTokenAdjective: StoryAdjectivesPool,
		StoryTokenMaterial:  StoryMaterialsPool,
		StoryTokenShaping:   StoryShapingPool,
	},
	TokenIsMandatory: map[string]bool{
		StoryTokenWorld: true,
	},
	Tokens:    StoryTokens,
	Templates: StoryTemplates,
}

var StoryTemplates = []string{
	"[INTRO] [WORLD] was [CREATION] by [GOD].",
	"[INTRO] [WORLD] was [SHAPING] from a [ADJECTIVE] [MATERIAL].",
	"[INTRO] [WORLD] was [SHAPING] from a [MATERIAL] by [GOD].",
}

const (
	StoryTokenIntro     = "[INTRO]"
	StoryTokenWorld     = "[WORLD]"
	StoryTokenCreation  = "[CREATION]"
	StoryTokenGod       = "[GOD]"
	StoryTokenAdjective = "[ADJECTIVE]"
	StoryTokenMaterial  = "[MATERIAL]"
	StoryTokenShaping   = "[SHAPING]"
)

var StoryTokens = []string{
	StoryTokenIntro,
	StoryTokenWorld,
	StoryTokenCreation,
	StoryTokenGod,
	StoryTokenAdjective,
	StoryTokenMaterial,
	StoryTokenShaping,
}

// StoryIntros contains the intro lines for the world creation mythos.
var StoryIntros = []string{
	"Long ago,",
	"As it is written in the ancient texts,",
	"According to the legends,",
	"According to the ancient texts,",
	"In the time before time,",
	"In the beginning,",
	"During the spark of creation,",
}

var StoryCreationPool = []string{
	"created by",
	"shaped in a dream of",
	"given existence by",
	"brought into being by",
}

var StoryShapingPool = []string{
	"formed",
	"shaped",
	"created",
	"made",
}

var StoryAdjectivesPool = []string{
	"lone",
	"pure",
	"perfect",
	"perfectly round",
	"perfectly square",
	"precious",
	"beautiful",
	"beautifully round",
}

var StoryMaterialsPool = []string{
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
