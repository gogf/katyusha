// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package krpc

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/gogf/katyusha/balancer"
	"github.com/gogf/katyusha/discovery"
)

// DefaultGrpcDialOptions returns the default options for creating grpc client connection.
func (c krpcClient) DefaultGrpcDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(
			`{"loadBalancingPolicy": "%s"}`,
			balancer.RoundRobin,
		)),
	}
}

// NewGrpcClientConn NewGrpcConn creates and returns a client connection for given service `appId`.
func (c krpcClient) NewGrpcClientConn(appID string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if err := discovery.InitDiscoveryFromConfig(); err != nil {
		return nil, err
	}
	grpcClientOptions := make([]grpc.DialOption, 0)
	grpcClientOptions = append(grpcClientOptions, c.DefaultGrpcDialOptions()...)
	if len(opts) > 0 {
		grpcClientOptions = append(grpcClientOptions, opts...)
	}
	grpcClientOptions = append(grpcClientOptions, c.ChainUnary(
		c.UnaryTracing,
		c.UnaryError))
	grpcClientOptions = append(grpcClientOptions, c.ChainStream(
		c.StreamTracing,
	))
	conn, err := grpc.Dial(discovery.DefaultScheme+":///"+appID, grpcClientOptions...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// ChainUnary creates a single interceptor out of a chain of many interceptors.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainUnaryClient(one, two, three) will execute one before two before three.
func (c krpcClient) ChainUnary(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(interceptors...)
}

// ChainStream creates a single interceptor out of a chain of many interceptors.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainStreamClient(one, two, three) will execute one before two before three.
func (c krpcClient) ChainStream(interceptors ...grpc.StreamClientInterceptor) grpc.DialOption {
	return grpc.WithChainStreamInterceptor(interceptors...)
}
