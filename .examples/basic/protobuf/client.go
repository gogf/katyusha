package protobuf

import (
	"google.golang.org/grpc"

	"github.com/gogf/katyusha/krpc"
)

const (
	AppID = "demo"
)

type Client struct {
	EchoClient
	TimeClient
}

func NewClient(options ...grpc.DialOption) (*Client, error) {
	conn, err := krpc.Client.NewGrpcClientConn(AppID, options...)
	if err != nil {
		return nil, err
	}
	return &Client{
		EchoClient: NewEchoClient(conn),
		TimeClient: NewTimeClient(conn),
	}, nil
}
