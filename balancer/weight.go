// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package balancer

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

const (
	// BlWeight  .
	BlWeight = "katyusha_balancer_weight"
	// WeightKey .
	WeightKey = "weight"
)

var (
	// defaultWeight is used if no weight configured.
	defaultWeight = 1
)

type weightPickerBuilder struct{}

type weightPicker struct {
	subConns []balancer.SubConn
}

func init() {
	balancer.Register(newWeightBuilder())
}

// newWeightBuilder creates a new weight balancer builder.
func newWeightBuilder() balancer.Builder {
	return base.NewBalancerBuilder(
		BlWeight,
		&weightPickerBuilder{},
		base.Config{HealthCheck: false},
	)
}

func (*weightPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}
	var subConns []balancer.SubConn
	for subConn, addr := range info.ReadySCs {
		for i := 0; i < getWeight(addr.Address); i++ {
			subConns = append(subConns, subConn)
		}
	}
	return &weightPicker{
		subConns: subConns,
	}
}

func (p *weightPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	return balancer.PickResult{
		SubConn: p.subConns[grand.Intn(len(p.subConns))],
	}, nil
}

func getWeight(addr resolver.Address) int {
	if addr.Attributes == nil {
		return defaultWeight
	}
	if v := addr.Attributes.Value(WeightKey); v != nil {
		return gconv.Int(v)
	}
	return defaultWeight
}
