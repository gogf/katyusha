package krpc

import (
	"context"
	"github.com/gogf/katyusha/krpc/internal/tracing"
	"google.golang.org/grpc"
)

// StreamTracing returns a grpc.StreamClientInterceptor suitable for use in a grpc.Dial call.
func (c *krpcClient) StreamTracing(
	ctx context.Context, desc *grpc.StreamDesc,
	cc *grpc.ClientConn, method string, streamer grpc.Streamer,
	callOpts ...grpc.CallOption) (grpc.ClientStream, error) {
	return tracing.StreamClientInterceptor(ctx, desc, cc, method, streamer, callOpts...)
}
