package krpc

import (
	"context"
	"github.com/gogf/katyusha/krpc/internal/tracing"
	"google.golang.org/grpc"
)

// ChainUnary returns a ServerOption that specifies the chained interceptor
// for unary RPCs. The first interceptor will be the outer most,
// while the last interceptor will be the inner most wrapper around the real call.
// All unary interceptors added by this method will be chained.
func (s *krpcServer) ChainUnary(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(interceptors...)
}

// ChainStream returns a ServerOption that specifies the chained interceptor
// for stream RPCs. The first interceptor will be the outer most,
// while the last interceptor will be the inner most wrapper around the real call.
// All stream interceptors added by this method will be chained.
func (s *krpcServer) ChainStream(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc.ChainStreamInterceptor(interceptors...)
}

// UnaryTracing returns a grpc.UnaryServerInterceptor suitable for use in a grpc.NewServer call.
func (s *krpcServer) UnaryTracing(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return tracing.UnaryServerInterceptor(ctx, req, info, handler)
}

// StreamTracing returns a grpc.StreamServerInterceptor suitable for use in a grpc.NewServer call.
func (s *krpcServer) StreamTracing(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	return tracing.StreamServerInterceptor(srv, ss, info, handler)
}