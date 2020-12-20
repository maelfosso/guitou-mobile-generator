package api

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// NewHTTPServer creates an HTTP server
func NewHTTPServer() http.Handler {
	svc := NewService()

	r := mux.NewRouter()

	generateHandler := httptransport.NewServer(
		makeGenerateEndpoint(svc),
		decodeGenerateRequest,
		encodeResponse,
	)

	r.Methods("POST").Path("/{id}/generate").Handler(generateHandler)

	return r
}
