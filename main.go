package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Guitou mobile generator")

	server := NewHttpServer()
	log.Fatal(http.ListenAndServe(":8000", server))
	// if err := http.ListenAndServe(":8000", server); err != nil {
	// 	log.Fatalf("could not listen on port 8000 %w", err)
	// }
}
