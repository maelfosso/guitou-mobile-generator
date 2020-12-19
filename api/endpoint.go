package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type serviceRequest struct {
	ProjectID string
}

type serviceResponse struct {
	Success bool  `json:"status,omitempty"`
	Err     error `json:"error,omitempty"`
}

func (r serviceResponse) error() error {
	return r.Err
}

func makeServiceEndpoint(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(serviceRequest)

		status, err := s.Generate(req.ProjectID)
		if err != nil {
			return serviceResponse{false, err}, nil
		}

		return serviceResponse{status, nil}, nil
	}
}
