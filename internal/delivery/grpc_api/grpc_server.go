package grpc

import (
	"fmt"
	taskpbv1 "github.com/blr-coder/task-proto/gen/go/task/v1"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	GRPCServer *grpc.Server
}

func NewGRPCServer(
	taskServer *TaskServiceServer,
	// someServer1 *SomeServiceServer1,
	// someServer2 *SomeServiceServer2,
) *Server {
	grpcServer := grpc.NewServer()

	// register grpcServerServices
	taskpbv1.RegisterTaskServiceServer(grpcServer, taskServer)
	// register other services... for example:
	//taskpbv1.RegisterSomeServiceServer1(grpcServer, someServer1)
	//taskpbv1.RegisterSomeServiceServer2(grpcServer, someServer2)

	return &Server{
		GRPCServer: grpcServer,
	}
}

func (s *Server) Run(port string) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}

	return s.GRPCServer.Serve(listener)
}
