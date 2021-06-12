package services

import (
	"github.com/go-kit/kit/transport/grpc"
	"github.com/golang/protobuf/protoc-gen-go/grpc"
	"guitou.cm/mobile/generator/models"
)

type IProjectClient interface {
	isProjectExists(id string) (*models.Project, error)
}

type gRPCProjectClient struct {
	cc grpc.Client
}

func (g gRPCProjectClient) isProjectExists(id string) (*models.Project, error) {
	return nil
}

func NewGrpcProjectClient() IProjectClient {
	return &gRPCProjectClient{}
}
