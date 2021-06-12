package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Guitou mobile generator")

	server := NewHttpServer()

	if err := http.ListenAndServe(":8000", server); err != nil {
		log.Fatal("could not listen on port 5000 %w")
	}
}
