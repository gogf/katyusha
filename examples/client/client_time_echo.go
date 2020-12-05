package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/katyusha/discovery"
	"github.com/gogf/katyusha/examples/proto"
	"github.com/gogf/katyusha/krpc"
	"golang.org/x/net/context"
	"time"
)

func main() {
	genv.SetMap(g.MapStrStr{
		discovery.EnvKeyEndpoints: "127.0.0.1:2379",
	})

	go func() {
		conn, err := krpc.NewGrpcClientConn("echo")
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		echoClient := proto.NewEchoClient(conn)
		for i := 0; i < 500; i++ {
			res, err := echoClient.Say(context.Background(), &proto.SayReq{Content: "Hello"})
			if err != nil {
				g.Log().Error(err)
				time.Sleep(time.Second)
				continue
			}
			time.Sleep(time.Second)
			g.Log().Print("Response:", res.Content)
		}
	}()
	time.Sleep(5 * time.Second)

	conn, err := krpc.NewGrpcClientConn("time")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := proto.NewTimeClient(conn)

	for i := 0; i < 500; i++ {
		res, err := client.Now(context.Background(), &proto.NowReq{})
		if err != nil {
			g.Log().Error(err)
			time.Sleep(time.Second)
			continue
		}
		time.Sleep(time.Second)
		g.Log().Print("Time:", res.Time)
	}
}
