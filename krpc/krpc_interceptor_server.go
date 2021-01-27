package krpc

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/util/gutil"
	"github.com/gogf/gf/util/gvalid"
	"github.com/gogf/katyusha/krpc/internal/grpctracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// UnaryError is the default unary interceptor for error converting from custom error to grpc error.
func (s *krpcServer) UnaryError(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, err := handler(ctx, req)
	if err != nil {
		code := gerror.Code(err)
		if code != -1 {
			err = status.Error(codes.Code(code), err.Error())
		}
	}
	return res, err
}

// UnaryRecover is the first interceptor that keep server not down from panics.
func (s *krpcServer) UnaryRecover(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	gutil.TryCatch(func() {
		res, err = handler(ctx, req)
	}, func(exception error) {
		err = gerror.WrapCode(int(codes.Internal), err, "panic recovered")
	})
	return
}

// Common validation unary interpreter.
func (s *krpcServer) UnaryValidate(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// It does nothing if there's no validation tag in the struct definition.
	if err := gvalid.CheckStruct(req, nil); err != nil {
		return nil, gerror.NewCode(
			int(codes.InvalidArgument),
			gerror.Current(err).Error(),
		)
	}
	return handler(ctx, req)
}

// UnaryTracing returns a grpc.UnaryServerInterceptor suitable for use in a grpc.NewServer call.
func (s *krpcServer) UnaryTracing(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return grpctracing.UnaryServerInterceptor(ctx, req, info, handler)
}

// StreamTracing returns a grpc.StreamServerInterceptor suitable for use in a grpc.NewServer call.
func (s *krpcServer) StreamTracing(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	return grpctracing.StreamServerInterceptor(srv, ss, info, handler)
}
