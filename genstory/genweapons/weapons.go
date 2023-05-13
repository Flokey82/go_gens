package genweapons

import (
	"github.com/Flokey82/go_gens/genstory"
)

const (
	WeaponTokenName            = "[WEAPON_NAME]"
	WeaponTokenType            = "[WEAPON_TYPE]"
	WeaponTokenSmith           = "[WEAPON_SMITH]"
	WeaponTokenCreationMethod  = "[WEAPON_CREATION_METHOD]"
	WeaponTokenMaterial        = "[WEAPON_MATERIAL]"
	WeaponTokenMaterialQuality = "[WEAPON_MATERIAL_QUALITY]"
	WeaponTokenAnyAge          = "[ANY_AGE]"
)

var WeaponTokens = []string{
	WeaponTokenName,
	WeaponTokenType,
	WeaponTokenSmith,
	WeaponTokenCreationMethod,
	WeaponTokenMaterial,
	WeaponTokenMaterialQuality,
	WeaponTokenAnyAge,
}

var WeaponsTextConfig = &genstory.TextConfig{
	TokenPools: map[string][]string{
		WeaponTokenName:            WeaponNames,
		WeaponTokenType:            WeaponTypes,
		WeaponTokenSmith:           WeaponSmiths,
		WeaponTokenCreationMethod:  WeaponCreationMethods,
		WeaponTokenMaterial:        WeaponMaterials,
		WeaponTokenMaterialQuality: WeaponMaterialQualities,
		WeaponTokenAnyAge:          WeaponAge,
	},
	TokenIsMandatory: map[string]bool{
		WeaponTokenName: true,
	},
	Tokens:         WeaponTokens,
	Templates:      WeaponIntros,
	UseAllProvided: true,
}

var WeaponIntros = []string{
	"[WEAPON_NAME:quote] is a [WEAPON_TYPE] [WEAPON_CREATION_METHOD] by [WEAPON_SMITH] from [WEAPON_MATERIAL_QUALITY] [WEAPON_MATERIAL].",
	"During the [ANY_AGE], [WEAPON_SMITH] [WEAPON_CREATION_METHOD] [WEAPON_NAME:quote], a [WEAPON_TYPE] made from [WEAPON_MATERIAL_QUALITY] [WEAPON_MATERIAL].",
	"[WEAPON_CREATION_METHOD] by [WEAPON_SMITH], [WEAPON_NAME:quote] is a [WEAPON_TYPE] made from [WEAPON_MATERIAL_QUALITY] [WEAPON_MATERIAL].",
	"[WEAPON_CREATION_METHOD] from [WEAPON_MATERIAL_QUALITY] [WEAPON_MATERIAL], [WEAPON_NAME:quote] is a [WEAPON_TYPE] [WEAPON_CREATION_METHOD] by [WEAPON_SMITH].",
	"This magnificent [WEAPON_TYPE] is called [WEAPON_NAME:quote]. [WEAPON_SMITH] [WEAPON_CREATION_METHOD] it from [WEAPON_MATERIAL_QUALITY] [WEAPON_MATERIAL].",
	"Striking fear into the hearts of the [ANY_AGE], [WEAPON_NAME:quote] is a [WEAPON_TYPE] [WEAPON_CREATION_METHOD] by [WEAPON_SMITH] from [WEAPON_MATERIAL_QUALITY] [WEAPON_MATERIAL].",
}

var WeaponTypes = []string{
	"axe",
	"bow",
	"club",
	"dagger",
	"flail",
	"hammer",
	"mace",
	"pike",
	"rapier",
	"spear",
	"sword",
	"whip",
}

var WeaponCreationMethods = []string{
	"crafted",
	"created",
	"shaped",
	"formed",
	"made",
	"forged",
}

var WeaponMaterials = []string{
	"adamantine",
	"aluminum",
	"bone",
	"bronze",
	"copper",
	"crystal",
	"diamond",
	"ebony",
	"gold",
	"iron",
	"ivory",
	"silver",
	"steel",
	"stone",
	"wood",
}

var WeaponMaterialQualities = []string{
	"fine",
	"pure",
	"perfect",
	"refined",
	"purified",
	"poor",
	"crude",
	"rough",
	"unrefined",
	"impure",
	"dirty",
	"filthy",
}

var WeaponSmiths = []string{
	"John Smith",
	"Jane Smith",
	"John Doe",
	"Jane Doe",
	"John Smithson",
	"Jane Smithson",
	"John Doe-Smith",
	"Jane Doe-Smith",
	"John Smith-Doe",
	"Jane Smith-Doe",
}

var WeaponNames = []string{
	"Flame",
	"Splinter",
	"Razor",
	"Edge",
	"Blade",
	"Scythe",
	"Rake",
	"Scourge",
	"Ender",
	"Slayer",
	"Lifedrinker",
	"Deathbringer",
	"Executioner",
	"Destroyer",
}

var WeaponAge = []string{
	"the age of steam",
	"the age of the forgotten",
	"the age of the gods",
	"the age of the machine",
}
