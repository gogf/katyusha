// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package balancer

import (
	"github.com/gogf/gf/util/grand"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync"
)

const BlRandom = "katyusha_balancer_random"

type randomPickerBuilder struct{}

type randomPicker struct {
	subConns []balancer.SubConn
	mu       sync.Mutex
}

func init() {
	balancer.Register(newRandomBuilder())
}

// newRandomBuilder creates a new random balancer builder.
func newRandomBuilder() balancer.Builder {
	return base.NewBalancerBuilderV2(BlRandom, &randomPickerBuilder{}, base.Config{HealthCheck: true})
}

func (*randomPickerBuilder) Build(buildInfo base.PickerBuildInfo) balancer.V2Picker {
	if len(buildInfo.ReadySCs) == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}
	var subConns []balancer.SubConn
	for subCon, _ := range buildInfo.ReadySCs {
		subConns = append(subConns, subCon)
	}
	return &randomPicker{
		subConns: subConns,
	}
}

func (p *randomPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	return balancer.PickResult{
		SubConn: p.subConns[grand.Intn(len(p.subConns))],
	}, nil
}
