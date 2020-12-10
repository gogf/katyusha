package krpc

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/katyusha/balancer"
	"github.com/gogf/katyusha/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	// DefaultGrpcClientConnOptions is the default options for creating grpc client connection.
	DefaultGrpcClientConnOptions = []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBalancerName(balancer.RoundRobin),
	}
)

// NewGrpcConn creates and returns a client connection for given service `appId`.
func (c *krpcClient) NewGrpcClientConn(appId string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	grpcClientOptions := make([]grpc.DialOption, 0)
	grpcClientOptions = append(grpcClientOptions, DefaultGrpcClientConnOptions...)
	if len(opts) > 0 {
		grpcClientOptions = append(grpcClientOptions, opts...)
	}
	grpcClientOptions = append(grpcClientOptions, c.ChainUnary(
		c.UnaryError,
	))
	conn, err := grpc.Dial(discovery.DefaultScheme+":///"+appId, grpcClientOptions...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// UnaryError handles the error types converting between grpc and gerror.
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

// ChainUnary creates a single interceptor out of a chain of many interceptors.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainUnaryClient(one, two, three) will execute one before two before three.
func (c *krpcClient) ChainUnary(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(interceptors...)
}

// ChainStream creates a single interceptor out of a chain of many interceptors.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainStreamClient(one, two, three) will execute one before two before three.
func (c *krpcClient) ChainStream(interceptors ...grpc.StreamClientInterceptor) grpc.DialOption {
	return grpc.WithChainStreamInterceptor(interceptors...)
}
