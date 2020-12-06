package krpc

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

// UnaryLogger is the default unary interpreter for logging purpose.
func (s *GrpcServer) UnaryLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var (
		start    = time.Now()
		res, err = handler(ctx, req)
		duration = time.Since(start)
	)
	if err != nil {
		if s.config.ErrorLogEnabled {
			s.Logger.Ctx(ctx).Stack(false).Stdout(s.config.LogStdout).File(s.config.ErrorLogPattern).Errorf(
				"%s, %.3fms, %+v, %+v, %+v",
				info.FullMethod, float64(duration)/1e6, req, res, err,
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
