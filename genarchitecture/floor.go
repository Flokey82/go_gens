package genarchitecture

type FloorStyle struct {
	BaseStyle
}

func generateFloorStyle(availableMaterials []string) FloorStyle {
	return FloorStyle{
		BaseStyle: generateBaseStyle(availableMaterials),
	}
}

type CeilingStyle struct {
	BaseStyle
}

func generateCeilingStyle(availableMaterials []string) CeilingStyle {
	return CeilingStyle{
		BaseStyle: generateBaseStyle(availableMaterials),
	}
}
