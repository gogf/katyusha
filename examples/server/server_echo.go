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
	"log"
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
		discovery.EnvKeyEndpoints: "127.0.0.1:2379",
	})

	register, err := discovery.NewRegister()
	if err != nil {
		log.Panic(err)
		return
	}

	s := krpc.NewGrpcServer(krpc.GrpcServerConfig{
		Addr: "0.0.0.0:" + gcmd.GetOpt("port"),
	})
	proto.RegisterEchoServer(s.Server, new(serviceEcho))
	s.Start()

	err = register.Register(&discovery.Service{
		AppId:    "echo",
		Version:  "v1.0",
		Address:  "127.0.0.1:" + gcmd.GetOpt("port"),
		Metadata: g.Map{"weight": 1},
	})
	if err != nil {
		panic(err)
	}
	s.Wait()

	//signalChan := make(chan os.Signal, 1)
	//signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	//<-signalChan
	//register.Unregister(service)
	//s.Stop()

}
