// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package balancer defines APIs for load balancing in gRPC.
package balancer

import (
	"fmt"

	"github.com/gogf/gf/v2/net/gsel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const (
	rawSvcKeyInSubConnInfo = `RawService`
)

func init() {
	Register(gsel.SelectorRandom, gsel.NewBuilderRandom())
	Register(gsel.SelectorWeight, gsel.NewBuilderWeight())
	Register(gsel.SelectorRoundRobin, gsel.NewBuilderRoundRobin())
	Register(gsel.SelectorLeastConnection, gsel.NewBuilderLeastConnection())
}

// Register registers the given balancer builder with the given name.
func Register(name string, builder gsel.Builder) {
	balancer.Register(base.NewBalancerBuilder(
		name,
		&Builder{builder: builder},
		base.Config{HealthCheck: true},
	))
}

// WithRandom returns a grpc.DialOption which enables random load balancing.
func WithRandom() grpc.DialOption {
	return doWithSelectorName(gsel.SelectorRandom)
}

// WithWeight returns a grpc.DialOption which enables weight load balancing.
func WithWeight() grpc.DialOption {
	return doWithSelectorName(gsel.SelectorWeight)
}

// WithRoundRobin returns a grpc.DialOption which enables round-robin load balancing.
func WithRoundRobin() grpc.DialOption {
	return doWithSelectorName(gsel.SelectorRoundRobin)
}

// WithLeastConnection returns a grpc.DialOption which enables the least connection load balancing.
func WithLeastConnection() grpc.DialOption {
	return doWithSelectorName(gsel.SelectorLeastConnection)
}

func doWithSelectorName(name string) grpc.DialOption {
	return grpc.WithDefaultServiceConfig(fmt.Sprintf(
		`{"loadBalancingPolicy": "%s"}`,
		name,
	))
}
