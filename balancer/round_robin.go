package balancer

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	"math/rand"
	"sync"
)

type roundRobinPickerBuilder struct{}

type roundRobinPicker struct {
	subConns []balancer.SubConn
	mu       sync.Mutex
	next     int
}

const RoundRobin = "round_robin_x"

func init() {
	balancer.Register(newRoundRobinBuilder())
}

// newRoundRobinBuilder creates a new roundrobin balancer builder.
func newRoundRobinBuilder() balancer.Builder {
	return base.NewBalancerBuilderV2(RoundRobin, &roundRobinPickerBuilder{}, base.Config{HealthCheck: true})
}

func (*roundRobinPickerBuilder) Build(buildInfo base.PickerBuildInfo) balancer.V2Picker {
	grpclog.Infof("roundrobinPicker: newPicker called with buildInfo: %v", buildInfo)

	if len(buildInfo.ReadySCs) == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn
	for subConn, subConnInfo := range buildInfo.ReadySCs {
		weight := GetWeight(subConnInfo.Address)
		for i := 0; i < weight; i++ {
			scs = append(scs, subConn)
		}
	}

	return &roundRobinPicker{
		subConns: scs,
		next:     rand.Intn(len(scs)),
	}
}

func (p *roundRobinPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	ret := balancer.PickResult{}
	p.mu.Lock()
	ret.SubConn = p.subConns[p.next]
	p.next = (p.next + 1) % len(p.subConns)
	p.mu.Unlock()
	return ret, nil
}
