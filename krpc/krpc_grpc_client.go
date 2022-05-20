// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package krpc

import (
	"fmt"

	"github.com/gogf/gf/v2/net/gsvc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gogf/katyusha/balancer"
)

// DefaultGrpcDialOptions returns the default options for creating grpc client connection.
func (c krpcClient) DefaultGrpcDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		balancer.WithRoundRobin(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}

// NewGrpcClientConn NewGrpcConn creates and returns a client connection for given service `appId`.
func (c krpcClient) NewGrpcClientConn(name string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	autoLoadAndRegisterEtcdRegistry()
	var (
		service           = gsvc.NewServiceWithName(name)
		grpcClientOptions = make([]grpc.DialOption, 0)
	)
	grpcClientOptions = append(grpcClientOptions, c.DefaultGrpcDialOptions()...)
	if len(opts) > 0 {
		grpcClientOptions = append(grpcClientOptions, opts...)
	}
	grpcClientOptions = append(grpcClientOptions, c.ChainUnary(
		c.UnaryTracing,
		c.UnaryError,
	))
	grpcClientOptions = append(grpcClientOptions, c.ChainStream(
		c.StreamTracing,
	))
	conn, err := grpc.Dial(fmt.Sprintf(`%s://%s`, gsvc.Schema, service.GetKey()), grpcClientOptions...)
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
