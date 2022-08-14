package main

import (
	"log"

	"github.com/Flokey82/go_gens/genpotions"
)

func main() {
	const (
		effectStamina      = "stamina"
		effectHealth       = "health"
		effectMagica       = "magica"
		effectDiarrhea     = "diarrhea"
		effectInvisibility = "invisibility"
	)
	apple := genpotions.NewIngredient("Apple", effectStamina, effectHealth)
	sugarcane := genpotions.NewIngredient("Sugarcane", effectStamina, effectMagica)
	rottenEgg := genpotions.NewIngredient("Rotten Egg", effectMagica, effectDiarrhea)
	daffodilPetals := genpotions.NewIngredient("Daffodil Petals", effectMagica, effectDiarrhea)
	potion, success := genpotions.CraftPotion(apple, sugarcane, rottenEgg, daffodilPetals)
	if !success {
		log.Println("Potion failed to craft")
	} else {
		log.Printf("Potion crafted: %s", potion.Name)
	}
}
