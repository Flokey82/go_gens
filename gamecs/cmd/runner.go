package main

import (
	"fmt"
	"github.com/Flokey82/go_gens/gamecs"
	"github.com/Flokey82/go_gens/gameloop"
	"time"
)

func main() {
	w := gamecs.New()

	itGrain := gamecs.NewItemType("grain", "grain")
	itBread := gamecs.NewItemType("bread", "food")
	itBread.Requires = []*gamecs.ItemType{itGrain}
	pBaker := gamecs.NewProfessionType("baker")
	pBaker.CanCraft = []*gamecs.ItemType{itBread}
	pFarmer := gamecs.NewProfessionType("farmer")
	pFarmer.CanCraft = []*gamecs.ItemType{itGrain}

	a1 := w.NewChar()
	a1.SetProfession(w, pBaker)

	a2 := w.NewChar()
	a2.SetProfession(w, pFarmer)

	// Add a number of characters.
	for i := 0; i < 5; i++ {
		w.NewChar()
	}

	gl := gameloop.New(200*time.Millisecond, w.Update)
	gl.Start()
	fmt.Scanln()
	fmt.Println("done")
	w.ExportGif("rgb.gif")
}
