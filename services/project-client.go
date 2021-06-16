package services

import (
	"fmt"

	"google.golang.org/grpc"
	"guitou.cm/mobile/generator/protos"
)

// type IProjectClient interface {
// 	isProjectExists(id string) (*models.Project, error)
// }

// type gRPCProjectClient struct {
// 	client protos.ProjectsClient
// }

// func (g gRPCProjectClient) isProjectExists(id string) (*models.Project, error) {
// 	return nil, nil
// }

const (
	PROJECT_MSVC_GRPC = "projects-api:50051"
)

func NewGrpcProjectClient() (protos.ProjectsClient, func() error, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial(PROJECT_MSVC_GRPC, opts...)
	if err != nil {
		return nil, func() error { return nil }, fmt.Errorf("fail to dial [project-api] grpc server: %v", err)
	}

	return protos.NewProjectsClient(conn), conn.Close, nil
	// &gRPCProjectClient{
	// 	client: protos.NewProjectsClient(conn),
	// }
}
