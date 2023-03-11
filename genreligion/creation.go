package genreligion

import (
	"strings"

	"github.com/Flokey82/go_gens/genstory"
)

// NewCreation generates a new world creation mythos using the given seed.
func (g *Generator) NewCreation(god *Deity) (string, error) {
	// Pick a random world name.
	worldName := strings.Title(g.lang.GetWord("world"))

	// Pick a random strategy for creation of the world.
	tokenReplacements := []genstory.TokenReplacement{{Token: CreationTokenWorld, Replacement: worldName}}
	if god != nil {
		tokenReplacements = append(tokenReplacements, genstory.TokenReplacement{
			Token:       CreationTokenGod,
			Replacement: god.FullName(),
		})
	}
	return CreationConfig.Generate(tokenReplacements)
}

var CreationConfig = &genstory.TextConfig{
	TokenPools: map[string][]string{
		CreationTokenIntro:     CreationIntros,
		CreationTokenCreation:  WorldCreationPool,
		CreationTokenAdjective: WorldMaterialAdjectivesPool,
		CreationTokenMaterial:  WorldMaterialsPool,
		CreationTokenShaping:   WorldShapingPool,
	},
	TokenIsMandatory: map[string]bool{
		CreationTokenWorld: true,
		CreationTokenGod:   true,
	},
	Tokens:         CreationTokens,
	Templates:      StoryTemplates,
	UseAllProvided: true,
}

var StoryTemplates = []string{
	"[INTRO] [WORLD] was [CREATION] by [GOD].",
	"[INTRO] [WORLD] was [SHAPING] from a [ADJECTIVE] [MATERIAL].",
	"[INTRO] [WORLD] was [SHAPING] from a [MATERIAL] by [GOD].",
}

const (
	CreationTokenIntro     = "[INTRO]"
	CreationTokenWorld     = "[WORLD]"
	CreationTokenCreation  = "[CREATION]"
	CreationTokenGod       = "[GOD]"
	CreationTokenAdjective = "[ADJECTIVE]"
	CreationTokenMaterial  = "[MATERIAL]"
	CreationTokenShaping   = "[SHAPING]"
)

var CreationTokens = []string{
	CreationTokenIntro,
	CreationTokenWorld,
	CreationTokenCreation,
	CreationTokenGod,
	CreationTokenAdjective,
	CreationTokenMaterial,
	CreationTokenShaping,
}

// CreationIntros contains the intro lines for the world creation mythos.
var CreationIntros = []string{
	"Long ago,",
	"As it is written in the ancient texts,",
	"According to the legends,",
	"According to the ancient texts,",
	"In the time before time,",
	"In the beginning,",
	"During the spark of creation,",
}

var WorldCreationPool = []string{
	"created by",
	"shaped in a dream of",
	"given existence by",
	"brought into being by",
}

var WorldShapingPool = []string{
	"formed",
	"shaped",
	"created",
	"made",
}

var WorldMaterialAdjectivesPool = []string{
	"lone",
	"pure",
	"perfect",
	"perfectly round",
	"perfectly square",
	"precious",
	"beautiful",
	"beautifully round",
}

var WorldMaterialsPool = []string{
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
