package gameconstants

// SettlementType indicates what type of settlement a settlement is.
// (Thorpe, Hamlet, Village, Town, City, Metropolis)
type SettlementType int

// Settlement types (per population).
const (
	SettlementTypeUnset SettlementType = iota
	SettlementTypeThorpe
	SettlementTypeHamlet
	SettlementTypeVillage
	SettlementTypeTownSmall
	SettlementTypeTownMedium
	SettlementTypeTownLarge
	SettlementTypeCitySmall
	SettlementTypeCityMedium
	SettlementTypeCityLarge
	SettlementTypeMetropolis
)

// String returns a human readable string for the settlement type.
func (s SettlementType) String() string {
	switch s {
	case SettlementTypeUnset:
		return "Unset"
	case SettlementTypeThorpe:
		return "Thorpe"
	case SettlementTypeHamlet:
		return "Hamlet"
	case SettlementTypeVillage:
		return "Village"
	case SettlementTypeTownSmall:
		return "Small Town"
	case SettlementTypeTownMedium:
		return "Medium Town"
	case SettlementTypeTownLarge:
		return "Large Town"
	case SettlementTypeCitySmall:
		return "Small City"
	case SettlementTypeCityMedium:
		return "Medium City"
	case SettlementTypeCityLarge:
		return "Large City"
	case SettlementTypeMetropolis:
		return "Metropolis"
	default:
		return "Unknown"
	}
}

// Min population for each settlement type.
// NOTE: These values are pure guesswork.
// I chose the values so they'd roughly double with each step in the
// larger settlement types.
const (
	SettlementPopMinHamlet     = 50
	SettlementPopMinVillage    = 200
	SettlementPopMinTownSmall  = 2000
	SettlementPopMinTownMedium = 6000
	SettlementPopMinTownLarge  = 15000
	SettlementPopMinCitySmall  = 30000
	SettlementPopMinCityMedium = 60000
	SettlementPopMinCityLarge  = 125000
	SettlementPopMinMetropolis = 250000
)

// GetSettlementType returns the settlement type for a given population.
func GetSettlementType(pop int) SettlementType {
	// NOTE: The status of "city" and "town" had to be acquired from
	// the local lord or king. I'm not sure how this has originally
	// worked, but I assume that the local lord or king would have
	// granted the status to a settlement if it had a certain population.
	//
	// What I DO know is, that a large village could pay a fee to the
	// local lord or king to become a town, have a city wall constructed
	// and get extra protection. The protection would be bought with a
	// tithe or taxes.
	//
	// Also, the official title would influence how attractive the settlement
	// would be for people to move to, and the type of people that would
	// move to the settlement. Towns would attract merchants and craftsmen,
	// while cities would attract nobles and rich merchants.
	// Thorpes, hamlets, and villages would attract farmers and possibly peasants.
	switch {
	case pop < SettlementPopMinHamlet:
		return SettlementTypeThorpe
	case pop < SettlementPopMinVillage:
		return SettlementTypeHamlet
	case pop < SettlementPopMinTownSmall:
		return SettlementTypeVillage
	case pop < SettlementPopMinTownMedium:
		return SettlementTypeTownSmall
	case pop < SettlementPopMinTownLarge:
		return SettlementTypeTownMedium
	case pop < SettlementPopMinCitySmall:
		return SettlementTypeTownLarge
	case pop < SettlementPopMinCityMedium:
		return SettlementTypeCitySmall
	case pop < SettlementPopMinCityLarge:
		return SettlementTypeCityMedium
	case pop < SettlementPopMinMetropolis:
		return SettlementTypeCityLarge
	default:
		return SettlementTypeMetropolis
	}
}

// NomadCampType indicates what type of camp a camp is.
type NomadCampType int

// Camp types.
const (
	CampTypeUnset NomadCampType = iota
	CampTypeSmall
	CampTypeFamily
	CampTypeClan
	CampTypeTribe
)

// String returns a human readable string for the camp type.
func (c NomadCampType) String() string {
	switch c {
	case CampTypeUnset:
		return "Unset"
	case CampTypeSmall:
		return "Small"
	case CampTypeFamily:
		return "Family"
	case CampTypeClan:
		return "Clan"
	case CampTypeTribe:
		return "Tribe"
	default:
		return "Unknown"
	}
}

// Min population for each nomad camp type.
// NOTE: These values are pure guesswork.
const (
	CampPopMinFamily = 10
	CampPopMinClan   = 50
	CampPopMinTribe  = 200
)

// GetNomadCampType returns the camp type for a given population.
func GetNomadCampType(pop int) NomadCampType {
	switch {
	case pop < CampPopMinFamily:
		return CampTypeSmall
	case pop < CampPopMinClan:
		return CampTypeFamily
	case pop < CampPopMinTribe:
		return CampTypeClan
	default:
		return CampTypeTribe
	}
}
