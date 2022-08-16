package main

import (
	"log"

	"github.com/Flokey82/go_gens/genpotion"
)

func main() {
	const (
		effectStamina      = "stamina"
		effectHealth       = "health"
		effectMagica       = "magica"
		effectDiarrhea     = "diarrhea"
		effectInvisibility = "invisibility"
	)
	apple := genpotion.NewIngredient("Apple", effectStamina, effectHealth)
	sugarcane := genpotion.NewIngredient("Sugarcane", effectStamina, effectMagica)
	rottenEgg := genpotion.NewIngredient("Rotten Egg", effectMagica, effectDiarrhea)
	daffodilPetals := genpotion.NewIngredient("Daffodil Petals", effectMagica, effectDiarrhea)
	potion, success := genpotion.CraftPotion(apple, sugarcane, rottenEgg, daffodilPetals)
	if !success {
		log.Println("Potion failed to craft")
	} else {
		log.Printf("Potion crafted: %s", potion.Name)
	}
}
