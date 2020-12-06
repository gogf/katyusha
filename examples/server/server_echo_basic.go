package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/katyusha/discovery"
	"github.com/gogf/katyusha/examples/protocol"
	"github.com/gogf/katyusha/examples/service"
	"github.com/gogf/katyusha/krpc"
)

func main() {
	genv.SetMap(g.MapStrStr{
		discovery.EnvKeyAppId:     `echo`,
		discovery.EnvKeyMetaData:  `{"weight":100}`,
		discovery.EnvKeyEndpoints: "127.0.0.1:2379",
	})
	s := krpc.NewGrpcServer()
	protocol.RegisterEchoServer(s.Server, new(service.Echo))
	s.Run()
}
