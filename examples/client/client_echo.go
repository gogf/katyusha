package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/katyusha/discovery"
	"github.com/gogf/katyusha/examples/protocol"
	"github.com/gogf/katyusha/krpc"
	"golang.org/x/net/context"
	"time"
)

// go run client_echo.go
func main() {
	genv.SetMap(g.MapStrStr{
		discovery.EnvKeyEndpoints: "127.0.0.1:2379",
	})

	conn, err := krpc.Client.NewGrpcClientConn("echo")
	if err != nil {
		g.Log().Fatal(err)
	}
	defer conn.Close()

	echoClient := protocol.NewEchoClient(conn)
	for i := 0; i < 500; i++ {
		res, err := echoClient.Say(context.Background(), &protocol.SayReq{Content: "Hello"})
		if err != nil {
			g.Log().Error(err)
			time.Sleep(time.Second)
			continue
		}
		time.Sleep(time.Second)
		g.Log().Print("Response:", res.Content)
	}
}
