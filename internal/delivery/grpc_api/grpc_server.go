package grpc

import (
	taskpbv1 "github.com/blr-coder/task-proto/gen/go/task/v1"
	"google.golang.org/grpc"
)

func NewGRPCServer(
	taskServer *TaskServiceServer,
	//someServer *someServiceServer,
) *grpc.Server {
	grpcServer := grpc.NewServer()

	// register grpcServerServices
	taskpbv1.RegisterTaskServiceServer(grpcServer, taskServer)
	// register other services...

	return grpcServer
}
