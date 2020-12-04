package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/katyusha/discovery"
	"github.com/gogf/katyusha/examples/proto"
	"github.com/gogf/katyusha/krpc"
	"golang.org/x/net/context"
	"log"
)

type helloServer struct{}

func (s *helloServer) Say(ctx context.Context, req *proto.SayReq) (*proto.SayResp, error) {
	text := "Hello " + req.Content + ", I am " + gcmd.GetOpt("node")
	g.Log().Println("Say:", text)
	return &proto.SayResp{Content: text}, nil
}

// go run main.go -node node1 -port 28544
// go run main.go -node node2 -port 18562
// go run main.go -node node3 -port 27772
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
	proto.RegisterTestServer(s.Server, new(helloServer))
	s.Start()

	err = register.Register(&discovery.Service{
		Name:     "test",
		AppId:    "test",
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
