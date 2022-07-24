package main

import (
	"github.com/Flokey82/go_gens/simmotive"
	"log"
	"time"
)

func main() {
	m := simmotive.NewMotive()
	m.InitMotives()
	m.Log("Your sim was born into the world")
	for {
		m.SimMotives(1)
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
