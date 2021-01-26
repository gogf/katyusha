package krpc

import (
	"github.com/gogf/katyusha/krpc/internal/tracing"
	"google.golang.org/grpc"
)

// StreamTracing returns a grpc.StreamServerInterceptor suitable for use in a grpc.NewServer call.
func (s *GrpcServer) StreamTracing(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	return tracing.StreamServerInterceptor(srv, ss, info, handler)
}
