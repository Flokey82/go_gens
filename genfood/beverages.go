package genfood

// Types of beverages.
const (
	BeverageTea      = "tea"
	BeverageBeer     = "beer"
	BeverageWine     = "wine"
	BeverageCoffee   = "coffee"
	BeverageJuice    = "juice"
	BeverageCocktail = "cocktail"
	BeverageWater    = "water"
)

// Type of refinement.
const (
	RefinementFermented = "fermented" // Beer, wine
	RefinementDestilled = "destilled" // Spirits
	RefinementBrewed    = "brewed"    // Coffee, tea
	RefinementPressed   = "pressed"   // Juice
	RefinementFiltered  = "filtered"  // Water
)

// Beverage.
//
// TODO: Indicate if beverage is ...
// - carbonated
// - hot, room temperature or cold
// - alcoholic (can be derived from ingredients and refinement)
// - sweet, sour, bitter, salty, umami (can be derived from ingredients)
type Beverage struct {
	Type       string   // Type of beverage
	Refinement string   // Fermented, destilled, brewed, pressed, filtered
	Primary    string   // Primary ingredient
	Secondary  []string // Secondary ingredients
}

// Description returns a description of the beverage.
func (b Beverage) Description() string {
	return "A " + b.Type + " " + b.Refinement + " from " + b.Primary + "."
}
