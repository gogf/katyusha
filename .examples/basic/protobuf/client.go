package protobuf

import (
	"github.com/gogf/katyusha/krpc"
	"google.golang.org/grpc"
)

const (
	AppId = "demo"
)

type Client struct {
	EchoClient
	TimeClient
}

func NewClient(options ...grpc.DialOption) (*Client, error) {
	conn, err := krpc.Client.NewGrpcClientConn(AppId, options...)
	if err != nil {
		return nil, err
	}
	return &Client{
		EchoClient: NewEchoClient(conn),
		TimeClient: NewTimeClient(conn),
	}, nil
}
