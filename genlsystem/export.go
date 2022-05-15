package genlsystem

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

func ExportToPNG(filePath string, m image.Image) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	b := bufio.NewWriter(f)
	if err := png.Encode(b, m); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err := b.Flush(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Wrote %s\n", filePath)
}
