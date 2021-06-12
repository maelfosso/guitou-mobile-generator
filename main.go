package main

import (
	"log"
	"net/http"
	"time"

	"guitou.cm/mobile/generator/api"
	"guitou.cm/mobile/generator/db"
)

func main() {
	log.Println("Guitou mobile generator")

	dbConnexion := db.Init()

	srv := &http.Server{
		Handler:      api.NewHTTPServer(dbConnexion),
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
