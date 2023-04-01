package main

import (
	"log"
	"time"

	"github.com/Flokey82/go_gens/genfood"
)

func main() {
	for i := 0; i < 10; i++ {
		if res, err := genfood.FoodTextConfig.Generate(nil); err != nil {
			log.Println(i, err)
		} else {
			log.Println(i, res.Text)
		}
	}

	st := genfood.ExampleRules.NewStory(time.Now().UnixNano())
	if sto, err := st.Expand(); err != nil {
		log.Println(err)
	} else {
		log.Println(sto)
	}

	st2 := genfood.SandwichRules.NewStory(time.Now().UnixNano())
	if sto, err := st2.Expand(); err != nil {
		log.Println(err)
	} else {
		log.Println(sto)
	}
}
