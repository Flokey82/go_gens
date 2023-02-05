package simnpcs

type LocationType int

const (
	LocTypeNone LocationType = iota
	LocTypeShop
	LocTypeSmith
	LocTypeFarm
	LocTypeMine
	LocTypeInn
	LocTypeTown
	LocTypeHome
)

// LocationScale determines the likelyhood of an encounter.
type LocationScale int

const (
	LocScaleHouse LocationScale = iota
	LocScaleMerchantShopTavern
	LocScaleSquare
	LocScaleDistrict
	LocScaleTown
	LocScaleCity
	// etc
)

// TODO: Add features (farm plot, chicken coop) that allow passive production
// of goods... Features might require new routines (collecting eggs etc.)
type Location struct {
	ID       uint64        // Unique ID
	Name     string        // Name of the location
	Type     LocationType  // Type of location
	Scale    LocationScale // Scale of the location
	Parent   *Location     // Parent location
	Children []*Location   // Child locations
	Host     *Character    // Host of the location (merchant, innkeeper, etc.)
	// Storage  []*Item
	// Features []*LocationFeature
	// Available indoor slots
	// Available outdoor slots
	// Currently producing []*Item
}

// NewLocation creates a new location.
func NewLocation(id uint64, name string, t LocationType, s LocationScale) *Location {
	return &Location{
		ID:    id,
		Name:  name,
		Type:  t,
		Scale: s,
	}
}

// Visit is called if a character visits a location.
func (l *Location) Visit(p *Character) {
	if l.Host == nil {
		return
	}

	l.Host.Interact(p, l) // Interact with the host
}

// AssignChild adds a child location to the location.
func (l *Location) AssignChild(c *Location) {
	l.Children = append(l.Children, c)
	c.Parent = l
}

/*
type LocationFeature struct {
	ID   uint64
	Name string
	// Which skill is needed / improved?
	// What routines will be introduced by this feature? (collecting eggs etc.)
	// What attentiveness is requred for success?
	CanProduce      []*Item
	ProductionSlots int
	// Required indoor slots
	// Required outdoor slots
}
*/
