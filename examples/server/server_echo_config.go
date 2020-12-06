package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/katyusha/discovery"
	"github.com/gogf/katyusha/examples/protocol"
	"github.com/gogf/katyusha/examples/service"
	"github.com/gogf/katyusha/krpc"
)

// go run server_echo.go -node node1 -port 8000
// go run server_echo.go -node node2 -port 8001
// go run server_echo.go -node node3 -port 8002
func main() {
	genv.SetMap(g.MapStrStr{
		discovery.EnvKeyAppId:     `echo`,
		discovery.EnvKeyMetaData:  `{"weight":100}`,
		discovery.EnvKeyEndpoints: "127.0.0.1:2379",
	})
	c := &krpc.GrpcServerConfig{
		Address:          fmt.Sprintf("0.0.0.0:%s", gcmd.GetOpt("port")),
		LogStdout:        true,
		ErrorLogEnabled:  true,
		AccessLogEnabled: true,
	}
	s := krpc.NewGrpcServer(c)
	protocol.RegisterEchoServer(s.Server, new(service.Echo))
	s.Run()
}
