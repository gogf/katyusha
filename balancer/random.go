package balancer

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	"math/rand"
	"sync"
	"time"
)

type randomPickerBuilder struct{}

type randomPicker struct {
	subConns []balancer.SubConn
	mu       sync.Mutex
	rand     *rand.Rand
}

const Random = "random_x"

func init() {
	balancer.Register(newRandomBuilder())
}

// newRandomBuilder creates a new random balancer builder.
func newRandomBuilder() balancer.Builder {
	return base.NewBalancerBuilderV2(Random, &randomPickerBuilder{}, base.Config{HealthCheck: true})
}

func (*randomPickerBuilder) Build(buildInfo base.PickerBuildInfo) balancer.V2Picker {
	grpclog.Infof("randomPicker: newPicker called with buildInfo: %v", buildInfo)
	if len(buildInfo.ReadySCs) == 0 {
		return base.NewErrPickerV2(balancer.ErrNoSubConnAvailable)
	}
	var scs []balancer.SubConn

	for subCon, subConnInfo := range buildInfo.ReadySCs {
		weight := GetWeight(subConnInfo.Address)
		for i := 0; i < weight; i++ {
			scs = append(scs, subCon)
		}
	}
	return &randomPicker{
		subConns: scs,
		rand:     rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (p *randomPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	ret := balancer.PickResult{}
	p.mu.Lock()
	ret.SubConn = p.subConns[p.rand.Intn(len(p.subConns))]
	p.mu.Unlock()
	return ret, nil
}
