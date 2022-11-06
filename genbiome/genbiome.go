// Package genbiome returns the biome for the given average precipitation and temperature.
package genbiome

import (
	"image/color"
)

// Biomes definition
// See: http://www-cs-students.stanford.edu/~amitp/game-programming/polygon-map-generation/#biomes
//
//
// Elevation ||| Moisture Zone ->
// Zone      V|| 6 (wet)    |    5    |     4    |     3    |     2     |   1 (dry)  |
// ==================================================================================
// 4 (high)   || SNOW ---------------------------| TUNDRA --| BARE -----|  SCORCHED -|
// ----------------------------------------------------------------------------------
// 3          || TAIGA ---------------| SHRUBLAND ----------| TEMPERATE DESERT ------|
// ----------------------------------------------------------------------------------
// 2          || TEMPERATE -| TEMPERATE ---------| GRASSLAND -----------| TEMPERATE -|
//            || RAIN FOREST| SEASONAL FOREST ---| ---------------------| DESERT ----|
// ----------------------------------------------------------------------------------
// 1 (low)    || TROPICAL RAIN FOREST | TROPICAL SEASONAL --| GRASSLAND | SUBTROPICAL|
//            || ---------------------| FOREST -------------| ----------| DESERT ----|
// ----------------------------------------------------------------------------------

var RedblobLookupTable = [][]int{
	{0x7, 0xA, 0x8, 0x8, 0x5, 0x5},
	{0x9, 0xA, 0xA, 0x8, 0x8, 0x5},
	{0x9, 0x9, 0x3, 0x3, 0x2, 0x2},
	{0xD, 0xC, 0x1, 0xB, 0xB, 0xB},
}

// The RedBlobGames biomes.
const (
	RedblobBiomeUnknown                 = 0x0
	RedblobBiomeTundra                  = 0x1
	RedblobBiomeTaiga                   = 0x2
	RedblobBiomeShrubland               = 0x3
	RedblobBiomeTemperateRainForest     = 0x4
	RedblobBiomeTropicalRainForest      = 0x5
	RedblobBiomeTropicalSeasonalForest  = 0x6
	RedblobBiomeSubtropicalDesert       = 0x7
	RedblobBiomeTemperateSeasonalForest = 0x8
	RedblobBiomeTemperateDesert         = 0x9 // combined with BioGrassland in whittaker
	RedblobBiomeGrassland               = 0xA // combined with BioTemperateDesert in whittaker
	RedblobBiomeSnow                    = 0xB // not in whittaker
	RedblobBiomeBare                    = 0xC // not in whittaker
	RedblobBiomeScorched                = 0xD // not in whittaker
)

// GetWhittakerBiome returns the RedBlobGames biome for the given height and moisture zones.
func GetRedblobBiome(height, moisture int) int {
	height = RedblobClampHeightToIndex(height)
	moisture = RedblobClampMoistureToIndex(moisture)
	return RedblobLookupTable[height][moisture]
}

// RedblobBiomeColor maps the (RedBlobGames) biome types to their color representation.
var RedblobBiomeColor = map[int]color.NRGBA{
	RedblobBiomeSnow:                    ColorSnowIce,
	RedblobBiomeTundra:                  {0xDD, 0xDD, 0xBB, 0xFF},
	RedblobBiomeBare:                    {0xBB, 0xBB, 0xBB, 0xFF},
	RedblobBiomeScorched:                {0x99, 0x99, 0x99, 0xFF},
	RedblobBiomeTaiga:                   {0xCC, 0xD4, 0xBB, 0xFF},
	RedblobBiomeShrubland:               {0xC4, 0xCC, 0xBB, 0xFF},
	RedblobBiomeTemperateDesert:         {0xE4, 0xE8, 0xCA, 0xFF},
	RedblobBiomeTemperateRainForest:     {0xA4, 0xC4, 0xA8, 0xFF},
	RedblobBiomeTemperateSeasonalForest: {0xB4, 0xC9, 0xA9, 0xFF},
	RedblobBiomeTropicalRainForest:      {0x9C, 0xBB, 0xA9, 0xFF},
	RedblobBiomeTropicalSeasonalForest:  {0xA9, 0xCC, 0xA4, 0xFF},
	RedblobBiomeGrassland:               {0xC4, 0xD4, 0xAA, 0xFF},
	RedblobBiomeSubtropicalDesert:       {0xE9, 0xDD, 0xC7, 0xFF},
}

// GetRedblobBiomeColor returns a color representing the (RedBlobGames) biome
// with the given height and moisture zones.
//
// The intensity (ranging from 0.0-1.0) allows darkening / brightening for (for example)
// indicating elevation. ... lower might be darker (0.0), higher might be brighter (1.0).
func GetRedblobBiomeColor(height, moisture int, intensity float64) color.NRGBA {
	c := RedblobBiomeColor[GetRedblobBiome(height, moisture)]
	return color.NRGBA{
		R: uint8(intensity * float64(c.R)),
		G: uint8(intensity * float64(c.G)),
		B: uint8(intensity * float64(c.B)),
		A: 255,
	}
}

func RedblobClampHeightToIndex(height int) int {
	if height < 1 {
		return 0
	}
	if height > 4 {
		return 3
	}
	return height - 1
}

func RedblobClampMoistureToIndex(moisture int) int {
	if moisture < 1 {
		return 0
	}
	if moisture > 6 {
		return 5
	}
	return moisture - 1
}

// Whittaker biomes based on:
// https://github.com/JoeyR/FastBiome
var WhittakerLookupTable = [][]int{
	// <- +30°C                                                                                                                                         | 0 °C                                                                     |<- -15°C
	{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x9, 0x9, 0x9, 0x9, 0x9, 0x9}, // 0dm precipitation
	{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x6, 0x6, 0x6, 0x6, 0x9, 0x9, 0x9, 0x9, 0x9, 0x9, 0x9},
	{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x6, 0x6, 0x6, 0x6, 0x8, 0x9, 0x9, 0x9, 0x9, 0x9, 0x9, 0x9, 0x9},
	{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x8, 0x8, 0x9, 0x9, 0x9, 0x9, 0x9, 0x0, 0x0, 0x0, 0x0},
	{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x7, 0x7, 0x7, 0x7, 0x7, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x8, 0x8, 0x8, 0x8, 0x9, 0x9, 0x9, 0x9, 0x9, 0x0, 0x0, 0x0, 0x0},
	{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x2, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x9, 0x9, 0x9, 0x9, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x2, 0x2, 0x2, 0x2, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x9, 0x9, 0x9, 0x9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x5, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x9, 0x9, 0x9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x3, 0x3, 0x3, 0x3, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x9, 0x9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x3, 0x3, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x6, 0x6, 0x6, 0x6, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x6, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x4, 0x8, 0x8, 0x8, 0x8, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x4, 0x4, 0x8, 0x8, 0x8, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x4, 0x8, 0x8, 0x8, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x4, 0x4, 0x8, 0x8, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x4, 0x4, 0x4, 0x8, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x4, 0x4, 0x4, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x2, 0x2, 0x2, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x2, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x4, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x0, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x0, 0x0, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x0, 0x0, 0x1, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x0, 0x0, 0x0, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{0x0, 0x0, 0x0, 0x0, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, // 45dm precipitation
}

// The Whittaker biomes.
const (
	WhittakerBiomeUnknown                       = 0x0
	WhittakerBiomeTropicalRainForest            = 0x1
	WhittakerBiomeTropicalSeasonalForestSavanna = 0x2
	WhittakerBiomeSubtropicalDesert             = 0x3
	WhittakerBiomeTemperateRainForest           = 0x4
	WhittakerBiomeTemperateSeasonalForest       = 0x5
	WhittakerBiomeWoodlandShrubland             = 0x6
	WhittakerBiomeTemperateGrasslandDesert      = 0x7
	WhittakerBiomeBorealForestTaiga             = 0x8
	WhittakerBiomeTundra                        = 0x9
)

// GetWhittakerBiome returns the Whittaker biome for the given temperature and precipitation.
func GetWhittakerBiome(temperature, precipitation int) int {
	temperature = WhittakerClampTemperatureToIndex(temperature)
	precipitation = WhittakerClampPrecipitationToIndex(precipitation)
	return WhittakerLookupTable[precipitation][temperature]
}

// WhittakerBiomeColor maps the (Whittaker) biome types to their color representation.
var WhittakerBiomeColor = map[int]color.NRGBA{
	WhittakerBiomeTropicalRainForest:            {0x9C, 0xBB, 0xA9, 0xFF},
	WhittakerBiomeTropicalSeasonalForestSavanna: {0xA9, 0xCC, 0xA4, 0xFF},
	WhittakerBiomeSubtropicalDesert:             {0xE9, 0xDD, 0xC7, 0xFF},
	WhittakerBiomeTemperateRainForest:           {0xA4, 0xC4, 0xA8, 0xFF},
	WhittakerBiomeTemperateSeasonalForest:       {0xB4, 0xC9, 0xA9, 0xFF},
	WhittakerBiomeWoodlandShrubland:             {0xC4, 0xCC, 0xBB, 0xFF},
	WhittakerBiomeTemperateGrasslandDesert:      {0xE4, 0xE8, 0xCA, 0xFF},
	WhittakerBiomeBorealForestTaiga:             {0xCC, 0xD4, 0xBB, 0xFF},
	WhittakerBiomeTundra:                        {0xDD, 0xDD, 0xBB, 0xFF},
}

// GetWhittakerModBiomeColor returns a color representing the (Whittaker) biome
// with the given temperature and precipitation.
//
// The intensity (ranging from 0.0-1.0) allows darkening / brightening for (for example)
// indicating elevation. ... lower might be darker (0.0), higher might be brighter (1.0).
func GetWhittakerBiomeColor(temperature, precipitation int, intensity float64) color.NRGBA {
	c := WhittakerBiomeColor[GetWhittakerBiome(temperature, precipitation)]
	return color.NRGBA{
		R: uint8(intensity * float64(c.R)),
		G: uint8(intensity * float64(c.G)),
		B: uint8(intensity * float64(c.B)),
		A: 255,
	}
}

var WhittakerLookupTableModified = [][]int{
	// <- +30°C                                                                                                                                         | 0 °C                                                                     |<- -15°C
	{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0x9, 0x9, 0x9, 0x9, 0x9, 0x9}, // 0dm precipitation
	{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x6, 0x6, 0x6, 0x6, 0x9, 0x9, 0x9, 0x9, 0x9, 0x9, 0x9},
	{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x6, 0x6, 0x6, 0x6, 0x8, 0x9, 0x9, 0x9, 0x9, 0x9, 0x9, 0x9, 0x9},
	{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x7, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x8, 0x8, 0x9, 0x9, 0x9, 0x9, 0x9, 0xB, 0xB, 0xB, 0xB},
	{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x7, 0x7, 0x7, 0x7, 0x7, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x8, 0x8, 0x8, 0x8, 0x9, 0x9, 0x9, 0x9, 0x9, 0xB, 0xB, 0xB, 0xB},
	{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0xE, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x9, 0x9, 0x9, 0x9, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0xE, 0xE, 0xE, 0xE, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x9, 0x9, 0x9, 0x9, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0xE, 0xE, 0xE, 0xE, 0xE, 0xE, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x5, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x9, 0x9, 0x9, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x3, 0x3, 0x3, 0x3, 0xE, 0xE, 0xE, 0xE, 0xE, 0xE, 0xE, 0xE, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x9, 0x9, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x3, 0x3, 0xE, 0xE, 0xE, 0xE, 0xE, 0xE, 0xE, 0xE, 0xE, 0xE, 0x6, 0x6, 0x6, 0x6, 0x6, 0x6, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x9, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0xE, 0xE, 0xE, 0xE, 0xE, 0xE, 0xE, 0xE, 0x2, 0x2, 0x2, 0x2, 0x6, 0x6, 0x6, 0x6, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0xE, 0xE, 0xE, 0xE, 0xE, 0xE, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x6, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0xE, 0xE, 0xE, 0xE, 0xE, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0xE, 0xE, 0xE, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0xE, 0xE, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0xE, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x4, 0x8, 0x8, 0x8, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x4, 0x4, 0x8, 0x8, 0x8, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x4, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x4, 0x4, 0x8, 0x8, 0x8, 0x8, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x4, 0x4, 0x4, 0x8, 0x8, 0xD, 0x8, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x4, 0x4, 0x4, 0x4, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x5, 0x5, 0x5, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x5, 0x5, 0x5, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x2, 0x2, 0x2, 0x2, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x2, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0x4, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x4, 0x4, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x4, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x4, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0x4, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0x4, 0x4, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x4, 0xC, 0xC, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0xC, 0xC, 0xC, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0xC, 0xC, 0xC, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0xC, 0xC, 0xC, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0xC, 0xC, 0xC, 0xC, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0xC, 0xC, 0xC, 0xC, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0xC, 0xC, 0xC, 0xC, 0xC, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0xC, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0xC, 0xC, 0xC, 0xC, 0xC, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0xC, 0xC, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0xC, 0xC, 0xC, 0xC, 0xC, 0xC, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0xC, 0xC, 0x1, 0x1, 0x1, 0x1, 0x1, 0xC, 0xC, 0xC, 0xC, 0xC, 0xC, 0xC, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0xC, 0xC, 0xC, 0x1, 0x1, 0x1, 0x1, 0xC, 0xC, 0xC, 0xC, 0xC, 0xC, 0xC, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB},
	{0xC, 0xC, 0xC, 0xC, 0x1, 0x1, 0xC, 0xC, 0xC, 0xC, 0xC, 0xC, 0xC, 0xC, 0xC, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0xD, 0x8, 0x8, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB, 0xB}, // 45dm precipitation
}

// The modified Whittaker biomes.
const (
	WhittakerModBiomeUnknown                 = 0x0
	WhittakerModBiomeTropicalRainForest      = 0x1
	WhittakerModBiomeTropicalSeasonalForest  = 0x2
	WhittakerModBiomeSubtropicalDesert       = 0x3
	WhittakerModBiomeTemperateRainForest     = 0x4
	WhittakerModBiomeTemperateSeasonalForest = 0x5
	WhittakerModBiomeWoodlandShrubland       = 0x6
	WhittakerModBiomeTemperateGrassland      = 0x7
	WhittakerModBiomeBorealForestTaiga       = 0x8
	WhittakerModBiomeTundra                  = 0x9
	WhittakerModBiomeColdDesert              = 0xA // Cold, dry winters, possibly warm, dry summers.
	WhittakerModBiomeSavannah                = 0xE // Hot, dry, grassland.

	// Not covered by Whittaker:
	WhittakerModBiomeSnow     = 0xB // Precipitation (any), temperature below freezing (includes glaciers).
	WhittakerModBiomeHotSwamp = 0xC // Hot, lots of precipitation, tropical swamps!
	WhittakerModBiomeWetlands = 0xD // Temperate but lots of precipitation (marshes, cold swamps, etc.).
)

// GetWhittakerModBiome returns the (modified) Whittaker biome for the given temperature and precipitation.
func GetWhittakerModBiome(temperature, precipitation int) int {
	temperature = WhittakerClampTemperatureToIndex(temperature)
	precipitation = WhittakerClampPrecipitationToIndex(precipitation)
	return WhittakerLookupTableModified[precipitation][temperature]
}

var (
	ColorRainForest     = color.NRGBA{R: 0x1d, G: 0xBA, B: 0x37, A: 0xFF}
	ColorSeasonalForest = color.NRGBA{R: 0x28, G: 0xAC, B: 0x2A, A: 0xFF}
	ColorSnowIce        = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
)

// WhittakerModBiomeColor maps the (modified Whittaker) biome types to their color representation.
var WhittakerModBiomeColor = map[int]color.NRGBA{
	WhittakerModBiomeUnknown:                 {0, 0, 0, 0xFF}, // {0xA9, 0, 0, 0xFF}, // TODO: Find better solution. It's wet and with relative cool temperature.
	WhittakerModBiomeTropicalRainForest:      ColorRainForest,
	WhittakerModBiomeTropicalSeasonalForest:  ColorSeasonalForest,
	WhittakerModBiomeSubtropicalDesert:       {0xfd, 0xe3, 0x8d, 0xFF}, // was 0xE9DDC7
	WhittakerModBiomeTemperateRainForest:     ColorRainForest,
	WhittakerModBiomeTemperateSeasonalForest: ColorSeasonalForest,
	WhittakerModBiomeWoodlandShrubland:       {0xC4, 0xCC, 0xBB, 0xFF},
	WhittakerModBiomeTemperateGrassland:      {0x97, 0xE9, 0x87, 0xFF},
	WhittakerModBiomeBorealForestTaiga:       {0x58, 0x9F, 0x5A, 0xFF},
	WhittakerModBiomeTundra:                  {0x69, 0xB0, 0x0B, 0xFF},
	WhittakerModBiomeColdDesert:              {0xDD, 0xDD, 0xC0, 0xFF},
	WhittakerModBiomeHotSwamp:                {0x96, 0x4B, 0x00, 0xFF},
	WhittakerModBiomeWetlands:                {0x0B, 0xB0, 0x1B, 0xFF}, // TODO: Find a better color
	WhittakerModBiomeSavannah:                {0xff, 0xd1, 0x45, 0xFF},
	WhittakerModBiomeSnow:                    ColorSnowIce,
}

// GetWhittakerModBiomeColor returns a color representing the (modified Whittaker) biome
// with the given temperature and precipitation.
//
// The intensity (ranging from 0.0-1.0) allows darkening / brightening for (for example)
// indicating elevation. ... lower might be darker (0.0), higher might be brighter (1.0).
func GetWhittakerModBiomeColor(temperature, precipitation int, intensity float64) color.NRGBA {
	c := WhittakerModBiomeColor[GetWhittakerModBiome(temperature, precipitation)]
	return color.NRGBA{
		R: uint8(intensity * float64(c.R)),
		G: uint8(intensity * float64(c.G)),
		B: uint8(intensity * float64(c.B)),
		A: 255,
	}
}

// Whittaker constants.
const (
	// Temperature range.
	MinTemperatureC = -15
	MaxTemperatureC = 30

	// Precipitation range.
	MinPrecipitationDM = 0
	MaxPrecipitationDM = 45
)

// WhittakerClampTemperatureToIndex clamps the given temperature (yearly average) to the range of the lookup table.
func WhittakerClampTemperatureToIndex(temp int) int {
	if temp <= MinTemperatureC {
		return 44
	}
	if temp >= MaxTemperatureC {
		return 0
	}
	return 44 - (temp + 15)
}

// WhittakerClampPrecipitationToIndex clamps the precipitation (in dm/year) to the range of the lookup table.
func WhittakerClampPrecipitationToIndex(prec int) int {
	if prec >= MaxPrecipitationDM {
		return 44
	}
	if prec <= MinPrecipitationDM {
		return 0
	}
	return prec
}

// AzgaarLookupTable maps temperature and moisture to biome.
// See: https://github.com/Azgaar/Fantasy-Map-Generator/blob/master/main.js
var AzgaarLookupTable = [5][26]int{
	// hot ↔ cold [>19°C; <-4°C]; dry ↕ wet
	{1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 10},
	{3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 9, 9, 9, 9, 10, 10, 10},
	{5, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 9, 9, 9, 9, 9, 10, 10, 10},
	{5, 6, 6, 6, 6, 6, 6, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 9, 9, 9, 9, 9, 9, 10, 10, 10},
	{7, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 9, 9, 9, 9, 9, 9, 9, 10, 10},
}

var azgaarBiomeName = [13]string{
	"Marine",
	"Hot desert",
	"Cold desert",
	"Savanna",
	"Grassland",                  //4
	"Tropical seasonal forest",   //5
	"Temperate deciduous forest", //6
	"Tropical rainforest",        //7
	"Temperate rainforest",       //8
	"Taiga",                      //9
	"Tundra",                     //10
	"Glacier",                    //11
	"Wetland",                    //12
}

// Biomes as per Azgaar's map generator.
const (
	AzgaarBiomeMarine = iota
	AzgaarBiomeHotDesert
	AzgaarBiomeColdDesert
	AzgaarBiomeSavanna
	AzgaarBiomeGrassland
	AzgaarBiomeTropicalSeasonalForest
	AzgaarBiomeTemperateDeciduousForest
	AzgaarBiomeTropicalRainforest
	AzgaarBiomeTemperateRainforest
	AzgaarBiomeTaiga
	AzgaarBiomeTundra
	AzgaarBiomeGlacier
	AzgaarBiomeWetland
)

// AzgaarBiomeColor maps Azgaar's biome IDs to their color representation.
var AzgaarBiomeColor = [13]color.NRGBA{
	{R: 0x46, G: 0x6e, B: 0xab, A: 0xFF},
	{R: 0xfb, G: 0xe7, B: 0x9f, A: 0xFF},
	{R: 0xb5, G: 0xb8, B: 0x87, A: 0xFF},
	{R: 0xd2, G: 0xd0, B: 0x82, A: 0xFF},
	{R: 0xc8, G: 0xd6, B: 0x8f, A: 0xFF},
	{R: 0xb6, G: 0xd9, B: 0x5d, A: 0xFF},
	{R: 0x29, G: 0xbc, B: 0x56, A: 0xFF},
	{R: 0x7d, G: 0xcb, B: 0x35, A: 0xFF},
	{R: 0x40, G: 0x9c, B: 0x43, A: 0xFF},
	{R: 0x4b, G: 0x6b, B: 0x32, A: 0xFF},
	{R: 0x96, G: 0x78, B: 0x4b, A: 0xFF},
	{R: 0xd5, G: 0xe7, B: 0xeb, A: 0xFF},
	{R: 0x0b, G: 0x91, B: 0x31, A: 0xFF},
}

// Habitability as per Azgaar's map generator.
var AzgaarBiomeHabitability = [13]int{0, 4, 10, 22, 30, 50, 100, 80, 90, 12, 4, 0, 12}

// Biome movement cost as per Azgaar's map generator.
var AzgaarBiomeMovementCost = [13]int{10, 200, 150, 60, 50, 70, 70, 80, 90, 200, 1000, 5000, 150}

// GetAzgaarBiome returns the biome for the given temperature, elevation, and precipitation.
//
// NOTE:
// - Elevation ranges from 0 to 100. (everything below 20 is considered water)
// - Temperature ranges from -5 to 20.
// - Moisture ranges from 0 to 20.
func GetAzgaarBiome(moisture, temperature, height int) int {
	if height < 0 {
		return 0 // marine biome: all water cells
	}
	if temperature < -5 {
		return 11 // permafrost biome
	}
	if isAzgaarWetland(moisture, temperature, height) {
		return 12 // wetland biome
	}
	moistureBand := min((moisture / 5), 4)             // [0-4]
	temperatureBand := min(max(20-temperature, 0), 25) // [0-25]
	return AzgaarLookupTable[moistureBand][temperatureBand]
}

// GetAzgaarBiomeColor returns the biome color for the given temperature, elevation, and precipitation.
//
// NOTE:
// - Elevation ranges from 0 to 100. (everything below 20 is considered water)
// - Temperature ranges from -5 to 20.
// - Moisture ranges from 0 to 20.
func GetAzgaarBiomeColor(moisture, temperature, height int, intensity float64) color.NRGBA {
	c := AzgaarBiomeColor[GetAzgaarBiome(moisture, temperature, height)]
	return color.NRGBA{
		R: uint8(intensity * float64(c.R)),
		G: uint8(intensity * float64(c.G)),
		B: uint8(intensity * float64(c.B)),
		A: 255,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// isAzgaarWetland returns true if the region is considered a wetland.
func isAzgaarWetland(moisture, temperature, height int) bool {
	if moisture > 40 && temperature > -2 && height < 25 {
		return true //near coast
	}
	if moisture > 24 && temperature > -2 && height > 24 && height < 60 {
		return true //off coast
	}
	return false
}
