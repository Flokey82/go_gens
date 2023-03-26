package genarchitecture

import (
	"math/rand"
	"os"

	svg "github.com/ajstarks/svgo"
)

func GenerateStyle(availableMaterials []string) Style {
	return Style{
		OuterDoorStyle: generateDoorStyle(availableMaterials),
		InnerDoorStyle: generateDoorStyle(availableMaterials),
		WindowStyle:    generateWindowStyle(availableMaterials),
		InnerWallStyle: generateWallStyle(availableMaterials),
		OuterWallStyle: generateWallStyle(availableMaterials),
		FloorStyle:     generateFloorStyle(availableMaterials),
		CeilingStyle:   generateCeilingStyle(availableMaterials),
		RoofStyle:      generateRoofStyle(availableMaterials),
	}
}

type Style struct {
	OuterDoorStyle DoorStyle
	InnerDoorStyle DoorStyle
	WindowStyle    WindowStyle
	InnerWallStyle WallStyle
	OuterWallStyle WallStyle
	FloorStyle     FloorStyle
	CeilingStyle   CeilingStyle
	RoofStyle      RoofStyle
}

func (s *Style) ExportSvg(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sv := svg.New(f)
	sv.Start(1000, 1000)
	s.OuterDoorStyle.DrawToSVG(sv, 0, 0, 100, 200)
	s.InnerDoorStyle.DrawToSVG(sv, 200, 0, 100, 200)
	sv.End()
	return nil
}

func (s Style) Description() string {
	// Describe the outside of the building.
	var description string
	description += "The outside of the building is made of " + s.OuterWallStyle.Description() + ". "
	description += "The roof is made of " + s.RoofStyle.Description() + ". "
	description += "The floor is made of " + s.FloorStyle.Description() + ". "
	description += "The ceiling is made of " + s.CeilingStyle.Description() + ". "
	description += "The outer door is made of " + s.OuterDoorStyle.Description() + ". "
	description += "The windows are made of " + s.WindowStyle.Description() + ". "

	// Describe the inside of the building.
	description += "The inside of the building is made of " + s.InnerWallStyle.Description() + ". "
	description += "The inner door is made of " + s.InnerDoorStyle.Description() + ". "

	return description
}

// BaseStyle is a struct that contains the basic style information for a
// building component. It can be embedded in other structs to provide
// a common set of style information.
type BaseStyle struct {
	Ornate   *Decoration // decoration of the component
	Material string      // main material of the component
	Finish   string      // finish of the component
}

// Description returns a string describing the style of the component.
func (b BaseStyle) Description() string {
	var leader string
	if b.Finish != FinishNone {
		leader = b.Finish + " "
	}
	leader += b.Material
	if b.Ornate == nil || b.Ornate.Type == DecorationTypeNone {
		return leader
	}
	return leader + " " + b.Ornate.Description()
}

func generateBaseStyle(availableMaterials []string) BaseStyle {
	// Generate a random decoration
	var ornate *Decoration
	if rand.Intn(4) == 0 {
		ornate = genDecoration()
	}

	return BaseStyle{
		Ornate:   ornate,
		Material: randomString(availableMaterials),
		Finish:   randomString(finishes),
	}
}

const (
	SizeNone   = 0
	SizeTiny   = 1
	SizeSmall  = 2
	SizeMedium = 3
	SizeLarge  = 4
	SizeHuge   = 5
)

const (
	MaterialStone   = "stone"
	MaterialWood    = "wood"
	MaterialBrick   = "brick"
	MaterialGlass   = "glass"
	MaterialMetal   = "metal"
	MaterialMarble  = "marble"
	MaterialPlaster = "plaster"
	MaterialPaper   = "paper"
	MaterialClay    = "clay"
	MaterialCeramic = "ceramic"
	MaterialPlastic = "plastic"
	MaterialLeather = "leather"
	MaterialFur     = "fur"
	MaterialBone    = "bone"
	MaterialShell   = "shell"
)

var Materials = []string{
	MaterialStone,
	MaterialWood,
	MaterialBrick,
	MaterialGlass,
	MaterialMetal,
	MaterialMarble,
	MaterialPlaster,
	MaterialPaper,
	MaterialClay,
	MaterialCeramic,
	MaterialPlastic,
	MaterialLeather,
	MaterialFur,
	MaterialBone,
	MaterialShell,
}

const (
	FinishNone      = "none"
	FinishPolished  = "polished"
	FinishRough     = "rough"
	FinishBurnished = "burnished"
	FinishPainted   = "painted"
	FinishPlastered = "plastered"
)

var finishes = []string{
	FinishNone,
	FinishPolished,
	FinishRough,
	FinishBurnished,
	FinishPainted,
	FinishPlastered,
}

func randomString(options []string) string {
	if len(options) == 0 {
		return ""
	}
	return options[rand.Intn(len(options)-1)]
}

func randomInt(min, max int) int {
	return rand.Intn(max-min) + min
}
