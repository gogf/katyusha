package krpc

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryError is the default unary interceptor for error converting from custom error to grpc error.
func (s *GrpcServer) UnaryError(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, err := handler(ctx, req)
	if err != nil {
		code := gerror.Code(err)
		if code != -1 {
			err = status.Error(codes.Code(code), err.Error())
		}
	}
	return res, err
}
