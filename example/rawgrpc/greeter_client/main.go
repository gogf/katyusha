package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gogf/katyusha/balancer"
	pb "github.com/gogf/katyusha/example/rawgrpc/helloworld"
	"github.com/gogf/katyusha/resolver"
)

func main() {
	resolver.SetRegistry(etcd.New("127.0.0.1:2379"))

	var (
		ctx     = gctx.New()
		service = gsvc.NewServiceWithName(`hello`)
	)
	// Set up a connection to the server.
	conn, err := grpc.Dial(
		fmt.Sprintf(`%s://%s`, gsvc.Schema, service.GetKey()),
		balancer.WithRandom(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		g.Log().Fatalf(ctx, "did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	for i := 0; i < 10; i++ {
		res, err := client.SayHello(ctx, &pb.HelloRequest{Name: `GoFrame`})
		if err != nil {
			g.Log().Fatalf(ctx, "could not greet: %+v", err)
		}
		g.Log().Printf(ctx, "Greeting: %s", res.GetMessage())
		time.Sleep(time.Second)
	}
}
