package krpc

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

// ChainUnaryClient creates a single interceptor out of a chain of many interceptors.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainUnaryClient(one, two, three) will execute one before two before three.
func ChainUnaryClient(interceptors ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	return grpc_middleware.ChainUnaryClient(interceptors...)
}

// ChainStreamClient creates a single interceptor out of a chain of many interceptors.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainStreamClient(one, two, three) will execute one before two before three.
func ChainStreamClient(interceptors ...grpc.StreamClientInterceptor) grpc.StreamClientInterceptor {
	return grpc_middleware.ChainStreamClient(interceptors...)
}
