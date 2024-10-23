package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

func SimpleLoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s", info.FullMethod)

	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	return resp, err
}
