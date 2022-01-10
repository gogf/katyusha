// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package krpc

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// internalUnaryLogger is the default unary interceptor for logging purpose.
func (s *GrpcServer) internalUnaryLogger(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	var (
		start    = time.Now()
		res, err = handler(ctx, req)
		duration = time.Since(start)
	)
	if err != nil {
		if s.config.ErrorLogEnabled {
			var (
				code          = gerror.Code(err)
				codeDetail    = code.Detail()
				codeDetailStr string
				grpcCode      codes.Code
				grpcMessage   string
			)
			if grpcStatus, ok := status.FromError(err); ok {
				grpcCode = grpcStatus.Code()
				grpcMessage = grpcStatus.Message()
			}
			if codeDetail != nil {
				codeDetailStr = gstr.Replace(fmt.Sprintf(`%+v`, codeDetail), "\n", " ")
			}
			content := fmt.Sprintf(
				`%s, %.3fms, %d, "%s", %+v, %+v, %d, "%s", "%s"`,
				info.FullMethod, float64(duration)/1e6, grpcCode, grpcMessage,
				req, res, code.Code(), code.Message(), codeDetailStr,
			)
			if s.config.ErrorStack {
				if stack := gerror.Stack(err); stack != "" {
					content += "\nStack:\n" + stack
				} else {
					content += ", " + err.Error()
				}
			} else {
				content += ", " + err.Error()
			}
			s.Logger.Stack(false).Stdout(s.config.LogStdout).File(s.config.ErrorLogPattern).Error(ctx, content)
		}
	} else {
		if s.config.AccessLogEnabled {
			content := fmt.Sprintf(
				"%s, %.3fms, %+v, %+v",
				info.FullMethod, float64(duration)/1e6, req, res,
			)
			s.Logger.Stdout(s.config.LogStdout).File(s.config.AccessLogPattern).Print(ctx, content)
		}
	}
	return res, err
}

// handleAccessLog handles the access logging for server.
func (s *GrpcServer) handleAccessLog(
	ctx context.Context, duration time.Duration, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) {
	if !s.config.AccessLogEnabled {
		return
	}

}

// handleErrorLog handles the error logging for server.
func (s *GrpcServer) handleErrorLog(
	ctx context.Context, err error, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) {
	// It does nothing if error logging is custom disabled.
	if !s.config.ErrorLogEnabled {
		return
	}

}
