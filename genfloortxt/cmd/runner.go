package main

import (
	"log"
	"os"

	"github.com/Flokey82/go_gens/genfloortxt"
)

func main() {
	file, err := os.Open("sample.plan")
	if err != nil {
		log.Fatal(err)
	}
	p := genfloortxt.ReadPlan(file)
	for _, line := range p.Render() {
		log.Println(line)
	}
}
