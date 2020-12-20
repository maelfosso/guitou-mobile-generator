package main

import (
	"log"
	"net/http"
	"time"

	"guitou.cm/mobile/generator/api"
)

func main() {
	log.Println("Guitou mobile generator")

	srv := &http.Server{
		Handler:      api.NewHTTPServer(),
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
