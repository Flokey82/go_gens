// Package genpotion implements a simple alchemy generator.
package genpotion

import (
	"fmt"
	"strings"
)

// Ingredient represents an alchemy ingredient with various effects.
//
// TODO:
// - Add the strength of each effect to the ingredient.
// - Add the price of the ingredient.
// - Add the rarity of the ingredient.
// - Use maybe ints instead of strings for effects?
type Ingredient struct {
	Name    string   // Name of the ingredient.
	Effects []string // Potential effects of the ingredient.
}

// NewIngredient returns a new ingredient with the given name and effects.
func NewIngredient(name string, effects ...string) *Ingredient {
	return &Ingredient{
		Name:    name,
		Effects: effects,
	}
}

// hasEffect returns true if the ingredient has the given effect.
func (in *Ingredient) hasEffect(effect string) bool {
	for _, e := range in.Effects {
		if e == effect {
			return true
		}
	}
	return false
}

// getSharedEffects returns array of effects present both in this ingredient and `other`
func (in *Ingredient) getSharedEffects(other *Ingredient) []string {
	var effects []string
	for _, e := range in.Effects {
		if !other.hasEffect(e) {
			effects = append(effects, e)
		}
	}
	return effects
}

// hasSomeEffects returns true if this ingredient has any of desired effects.
func (in *Ingredient) hasSomeEffects(effects []string) bool {
	for _, e := range effects {
		if in.hasEffect(e) {
			return true
		}
	}
	return false
}

// CanCraftPotion returns true if a potion can be successfully crafted from the
// given ingredients.
//
// NOTE: For a successful potion, each ingredient must share at least one effect
// with at least one other ingredient.
func CanCraftPotion(ingredients ...*Ingredient) bool {
	// Each ingredient must have at least one overlapping effect with another.
	for _, in1 := range ingredients {
		var valid bool
		for _, in2 := range ingredients {
			// Do not compare ingredient with itself.
			if in1.Name == in2.Name {
				continue
			}

			// If the ingredient has any overlapping effect, we can continue.
			if in1.hasSomeEffects(in2.Effects) {
				valid = true
				break
			}
		}

		// If no ingredient has any overlapping effect, the potion cannot be
		// crafted.
		if !valid {
			return false
		}
	}
	return true
}

// Potion represents an alchemical potion expressing different effects.
// TODO:
// - Add strength of each effect by summing up the effect strength of each ingredient.
// - Add price of potion based on ingredient prices.
type Potion struct {
	Name        string        // Name of the potion
	Ingredients []*Ingredient // List of ingredients
	Effects     []string      // Magical effects
}

// CraftPotion creates a new potion from the given ingredients.
// On failure, the function returns nil and false.
func CraftPotion(ingredients ...*Ingredient) (*Potion, bool) {
	// Check if we can craft a potion successfully given the ingredients.
	if !CanCraftPotion(ingredients...) {
		return nil, false
	}

	// Create a new potion.
	potion := &Potion{
		Ingredients: ingredients,
	}

	// Iterate through all ingredient effects and count the number
	// of occurrences.
	var effects []string
	effectCount := make(map[string]int)
	for _, in := range ingredients {
		for _, e := range in.Effects {
			// We collect all effects in a stable order.
			// Alternatively, we could just sort the result.
			if _, ok := effectCount[e]; !ok {
				effects = append(effects, e)
			}
			effectCount[e]++
		}
	}

	// Now iterate through the effects and assign all effects that were
	// encountered more than once.
	for _, e := range effects {
		// Skip unique effects.
		if effectCount[e] <= 1 {
			continue
		}
		potion.Effects = append(potion.Effects, e)
	}

	// Generate the name of the potion.
	// TODO: This sould be done a little bit more cleverly.
	potion.Name = fmt.Sprintf("Potion of %s", strings.Join(potion.Effects, ", "))
	return potion, true
}
