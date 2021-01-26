package krpc

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/util/gutil"
	"github.com/gogf/katyusha/krpc/internal/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
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

// UnaryLogger is the default unary interceptor for logging purpose.
func (s *GrpcServer) UnaryLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var (
		start    = time.Now()
		res, err = handler(ctx, req)
		duration = time.Since(start)
	)
	if err != nil {
		var (
			grpcCode    codes.Code
			grpcMessage string
		)
		grpcStatus, ok := status.FromError(err)
		if ok {
			grpcCode = grpcStatus.Code()
			grpcMessage = grpcStatus.Message()
		} else {
			grpcMessage = err.Error()
		}
		if s.config.ErrorLogEnabled {
			s.Logger.Ctx(ctx).Stack(false).Stdout(s.config.LogStdout).File(s.config.ErrorLogPattern).Errorf(
				"%s, %.3fms, %+v, %+v, %d, %+v",
				info.FullMethod, float64(duration)/1e6, req, res, grpcCode, grpcMessage,
			)
		}
	} else {
		if s.config.AccessLogEnabled {
			s.Logger.Ctx(ctx).Stdout(s.config.LogStdout).File(s.config.AccessLogPattern).Printf(
				"%s, %.3fms, %+v, %+v",
				info.FullMethod, float64(duration)/1e6, req, res,
			)
		}
	}
	return res, err
}

// UnaryRecover is the first interceptor that keep server not down from panics.
func (s *GrpcServer) UnaryRecover(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	gutil.TryCatch(func() {
		res, err = handler(ctx, req)
	}, func(exception error) {
		err = gerror.WrapCode(int(codes.Internal), err, "panic recovered")
	})
	return
}

// UnaryTracing returns a grpc.UnaryServerInterceptor suitable for use in a grpc.NewServer call.
func (s *GrpcServer) UnaryTracing(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return tracing.UnaryServerInterceptor(ctx, req, info, handler)
}
