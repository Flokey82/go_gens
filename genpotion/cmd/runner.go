package main

import (
	"log"

	"github.com/Flokey82/go_gens/genpotion"
)

func main() {
	// Set up some effects.
	const (
		effectStamina      = "stamina"
		effectHealth       = "health"
		effectMagica       = "magica"
		effectDiarrhea     = "diarrhea"
		effectInvisibility = "invisibility"
	)

	// Set up some ingredients.
	apple := genpotion.NewIngredient("Apple", effectStamina, effectHealth)
	sugarcane := genpotion.NewIngredient("Sugarcane", effectStamina, effectMagica)
	rottenEgg := genpotion.NewIngredient("Rotten Egg", effectMagica, effectDiarrhea)
	daffodilPetals := genpotion.NewIngredient("Daffodil Petals", effectMagica, effectDiarrhea)

	// Craft a potion.
	potion, success := genpotion.CraftPotion(apple, sugarcane, rottenEgg, daffodilPetals)
	if !success {
		log.Println("Potion failed to craft")
	} else {
		log.Printf("Potion crafted: %s", potion.Name)
	}
}
