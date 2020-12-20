package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type generateRequest struct {
	ProjectID string
}

type generateResponse struct {
	Success bool  `json:"status,omitempty"`
	Err     error `json:"error,omitempty"`
}

func (r generateResponse) error() error {
	return r.Err
}

func makeGenerateEndpoint(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(generateRequest)

		status, err := s.Generate(req.ProjectID)
		if err != nil {
			return generateResponse{false, err}, nil
		}

		return generateResponse{status, nil}, nil
	}
}
