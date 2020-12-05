package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/katyusha/discovery"
	"github.com/gogf/katyusha/examples/proto"
	"github.com/gogf/katyusha/krpc"
	"golang.org/x/net/context"
)

type serviceTime struct{}

func (s *serviceTime) Now(ctx context.Context, r *proto.NowReq) (*proto.NowRes, error) {
	text := fmt.Sprintf(`%s: %s`, gcmd.GetOpt("node"), gtime.Now().String())
	return &proto.NowRes{Time: text}, nil
}

// go run server_time.go -node node1 -port 8100
// go run server_time.go -node node2 -port 8101
// go run server_time.go -node node3 -port 8102
func main() {
	genv.SetMap(g.MapStrStr{
		discovery.EnvKeyEndpoints: "127.0.0.1:2379",
	})

	s := krpc.NewGrpcServer(krpc.GrpcServerConfig{
		Addr: "0.0.0.0:" + gcmd.GetOpt("port"),
	})
	proto.RegisterTimeServer(s.Server, new(serviceTime))
	s.Start()

	err := discovery.Register(&discovery.Service{
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
