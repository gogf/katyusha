package krpc

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/katyusha/krpc/internal/grpctracing"
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
	return grpctracing.UnaryClientInterceptor(ctx, method, req, reply, cc, invoker, opts...)
}

// StreamTracing returns a grpc.StreamClientInterceptor suitable for use in a grpc.Dial call.
func (c *krpcClient) StreamTracing(
	ctx context.Context, desc *grpc.StreamDesc,
	cc *grpc.ClientConn, method string, streamer grpc.Streamer,
	callOpts ...grpc.CallOption) (grpc.ClientStream, error) {
	return grpctracing.StreamClientInterceptor(ctx, desc, cc, method, streamer, callOpts...)
}
