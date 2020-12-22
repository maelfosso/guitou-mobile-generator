package api

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"guitou.cm/mobile/generator/db"
	"guitou.cm/mobile/generator/repositories"
)

// NewHTTPServer creates an HTTP server
func NewHTTPServer(dbConnexion *db.MongoDBClient) http.Handler {

	projectRepository := repositories.NewProjectRepository(dbConnexion)
	// generateRepository, err := NewGenerateRepository(dbConnexion)

	svc := NewService(projectRepository) //, generateRepository)

	r := mux.NewRouter()

	generateHandler := httptransport.NewServer(
		makeGenerateEndpoint(svc),
		decodeGenerateRequest,
		encodeResponse,
	)

	r.Methods("POST").Path("/{id}/generate").Handler(generateHandler)

	return r
}
