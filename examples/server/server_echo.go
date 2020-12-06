package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/katyusha/discovery"
	"github.com/gogf/katyusha/examples/proto"
	"github.com/gogf/katyusha/krpc"
	"golang.org/x/net/context"
)

type serviceEcho struct{}

func (s *serviceEcho) Say(ctx context.Context, r *proto.SayReq) (*proto.SayRes, error) {
	g.Log().Println("Received:", r.Content)
	text := fmt.Sprintf(`%s: > %s`, gcmd.GetOpt("node"), r.Content)
	return &proto.SayRes{Content: text}, nil
}

// go run server_echo.go -node node1 -port 8000
// go run server_echo.go -node node2 -port 8001
// go run server_echo.go -node node3 -port 8002
func main() {
	genv.SetMap(g.MapStrStr{
		discovery.EnvKeyAppId:     `echo`,
		discovery.EnvKeyMetaData:  `{"weight":100}`,
		discovery.EnvKeyEndpoints: "127.0.0.1:2379",
	})
	s := krpc.NewGrpcServer(krpc.GrpcServerConfig{
		Address: fmt.Sprintf("0.0.0.0:%s", gcmd.GetOpt("port")),
	})
	proto.RegisterEchoServer(s.Server, new(serviceEcho))
	s.Run()
}
