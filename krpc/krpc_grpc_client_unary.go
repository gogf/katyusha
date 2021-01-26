package krpc

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/katyusha/krpc/internal/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryError handles the error types converting between grpc and gerror.
// Note that, the minus error code is only used locally which will not be sent to other side.
func (c *krpcClient) UnaryError(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if ok {
			if code := grpcStatus.Code(); code != 0 {
				return gerror.NewCode(int(code), grpcStatus.Message())
			}
			return gerror.New(grpcStatus.Message())
		}
	}
	return err
}

// UnaryTracing returns a grpc.UnaryClientInterceptor suitable for use in a grpc.Dial call.
func (c *krpcClient) UnaryTracing(
	ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return tracing.UnaryClientInterceptor(ctx, method, req, reply, cc, invoker, opts...)
}
