package village

// Building represents an instance of BuildingType.
type Building struct {
	*BuildingType
}

// BuildingType represents a class of building requiring and/or providing
// specific resources.
type BuildingType struct {
	Name string
	*Production
}

// NewBuildingType returns a new building type with the given name.
func NewBuildingType(name string) *BuildingType {
	return &BuildingType{
		Name:       name,
		Production: NewProduction(),
	}
}

// String implements the stringer function for this building type.
func (bt *BuildingType) String() string {
	return bt.Name
}

// NewBuilding returns a new building instance of this type.
func (bt *BuildingType) NewBuilding() *Building {
	return &Building{
		BuildingType: bt,
	}
}

// BuildingPool represents all known / registered building types.
type BuildingPool struct {
	Types    []*BuildingType
	Provides map[string][]*BuildingType
	Requires map[string][]*BuildingType
}

// NewBuildingPool returns a new, empty building pool.
func NewBuildingPool() *BuildingPool {
	return &BuildingPool{
		Provides: make(map[string][]*BuildingType),
		Requires: make(map[string][]*BuildingType),
	}
}

// FindType returns the given building type from the pool.
func (bp *BuildingPool) FindType(name string) *BuildingType {
	for _, bt := range bp.Types {
		if bt.Name == name {
			return bt
		}
	}
	return nil
}

// AddType registers the given building type in the pool.
func (bp *BuildingPool) AddType(b *BuildingType) {
	// TODO: Prevent duplicate entries.
	bp.Types = append(bp.Types, b)
	bp.Update()
}

// Update rebuilds the internal mappings for looking up building types
// by resource ID.
func (bp *BuildingPool) Update() {
	bp.Provides = make(map[string][]*BuildingType)
	bp.Requires = make(map[string][]*BuildingType)
	for _, b := range bp.Types {
		for key := range b.GetExcess() {
			bp.Provides[key] = append(bp.Provides[key], b)
		}
		for key := range b.GetMissing() {
			bp.Requires[key] = append(bp.Requires[key], b)
		}
	}
}
