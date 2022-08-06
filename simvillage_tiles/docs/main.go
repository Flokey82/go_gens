package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Println("Point browser at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
