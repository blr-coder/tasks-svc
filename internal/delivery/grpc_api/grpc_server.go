package grpc

import (
	"fmt"
	taskpbv1 "github.com/blr-coder/task-proto/gen/go/task/v1"
	"github.com/blr-coder/tasks-svc/internal/delivery/grpc_api/handlers"
	"github.com/blr-coder/tasks-svc/internal/delivery/grpc_api/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	GRPCServer *grpc.Server
}

func NewGRPCServer(
	taskHandler *handlers.TaskGRPCHandler,
	// someServer1 *SomeServiceServer1,
	// someServer2 *SomeServiceServer2,
) *Server {
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(),
		interceptors.SimpleLoggingInterceptor,
	))

	// Даёт возможность клиенту обратиться к серверу за списком доступных методов
	reflection.Register(grpcServer)

	// register grpcServerServices
	taskpbv1.RegisterTaskServiceServer(grpcServer, taskHandler)
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
