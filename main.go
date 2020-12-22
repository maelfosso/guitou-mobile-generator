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

// TODO
// 1. GO-Kit example Cargo Shipping: Repository usage
// 2. SSH for Gitlab Clone
// 3. gRPC Client from golang to Request a Microservice
// 4. MongoDB Database integration
