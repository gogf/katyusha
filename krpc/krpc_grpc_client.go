package krpc

import (
	"github.com/gogf/katyusha/balancer"
	"github.com/gogf/katyusha/discovery"
	"google.golang.org/grpc"
)

type GrpcClientConn struct {
	*grpc.ClientConn
}

var (
	DefaultGrpcClientConnOptions = []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBalancerName(balancer.RoundRobin),
	}
)

func NewGrpcClientConn(appId string, opts ...grpc.DialOption) (*GrpcClientConn, error) {
	grpcClientOptions := opts
	if len(grpcClientOptions) == 0 {
		grpcClientOptions = DefaultGrpcClientConnOptions
	}
	conn, err := grpc.Dial(discovery.DefaultScheme+":///"+appId, DefaultGrpcClientConnOptions...)
	if err != nil {
		return nil, err
	}
	return &GrpcClientConn{ClientConn: conn}, nil
}
