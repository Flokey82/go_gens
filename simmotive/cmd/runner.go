package main

import (
	"log"
	"time"

	"github.com/Flokey82/go_gens/simmotive"
)

func main() {
	m := simmotive.NewMotive()
	m.Init()
	m.Log("Your sim was born into the world")
	for {
		m.SimMotives()
		m.Clear()
		log.Printf("Day %d : [%02d:%02d]\n\n", m.ClockD, m.ClockH, m.ClockM)
		m.PrintMotives()
		log.Printf("\nLog")
		log.Printf("====")
		for _, str := range m.Logs {
			log.Println(str)
		}
		time.Sleep(50000000)
	}
}
