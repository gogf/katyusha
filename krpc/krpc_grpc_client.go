package krpc

import (
	"github.com/gogf/katyusha/balancer"
	"github.com/gogf/katyusha/discovery"
	"google.golang.org/grpc"
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
	conn, err := grpc.Dial(discovery.DefaultScheme+":///"+appId, grpcClientOptions...)
	if err != nil {
		return nil, err
	}
	return conn, nil
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
