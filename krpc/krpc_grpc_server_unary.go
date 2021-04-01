// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package krpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

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
			s.Logger.Ctx(ctx).
				Stack(false).
				Stdout(s.config.LogStdout).
				File(s.config.ErrorLogPattern).
				Errorf(
					"%s, %.3fms, %+v, %+v, %d, %+v",
					info.FullMethod, float64(duration)/1e6, req, res, grpcCode, grpcMessage,
				)
		}
	} else {
		if s.config.AccessLogEnabled {
			s.Logger.Ctx(ctx).
				Stdout(s.config.LogStdout).
				File(s.config.AccessLogPattern).
				Printf(
					"%s, %.3fms, %+v, %+v",
					info.FullMethod, float64(duration)/1e6, req, res,
				)
		}
	}
	return res, err
}
