package protobuf

import (
	"google.golang.org/grpc"

	"github.com/gogf/katyusha/krpc"
)

const (
	AppID = "demo"
)

type Client struct {
	conn *grpc.ClientConn
}

func NewClient(options ...grpc.DialOption) (*Client, error) {
	conn, err := krpc.Client.NewGrpcClientConn(AppID, options...)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Echo() EchoClient {
	return NewEchoClient(c.conn)
}

func (c *Client) Time() TimeClient {
	return NewTimeClient(c.conn)
}
