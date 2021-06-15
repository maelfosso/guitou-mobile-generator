package services

import (
	"log"

	"google.golang.org/grpc"
	"guitou.cm/mobile/generator/models"
	"guitou.cm/mobile/generator/protos"
)

type IProjectClient interface {
	isProjectExists(id string) (*models.Project, error)
}

type gRPCProjectClient struct {
	conn grpc.Client
}

func (g gRPCProjectClient) isProjectExists(id string) (*models.Project, error) {
	return nil, nil
}

const (
	PROJECT_MSVC_GRPC = "project-api:50051"
)

func NewGrpcProjectClient() IProjectClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBloc())

	conn, err := grpc.Dial(PROJECT_MSVC_GRPC, opts...)
	if err != nil {
		log.Fatalf("fail to dial [project-api] grpc server: %v", err)
	}

	return &gRPCProjectClient{
		conn: protos.NewProjectsClient(conn),
	}
}
