package main

import (
	"flag"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/katyusha/examples/proto"
	"github.com/gogf/katyusha/registry"
	"github.com/gogf/katyusha/server"
	etcd3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"log"
	"time"
)

var nodeID = flag.String("node", "node1", "node ID")
var port = flag.Int("port", 8080, "listening port")

type helloServer struct{}

func (s *helloServer) Say(ctx context.Context, req *proto.SayReq) (*proto.SayResp, error) {
	text := "Hello " + req.Content + ", I am " + *nodeID
	log.Println(text)

	return &proto.SayResp{Content: text}, nil
}

func StartService() {

	register, err := registry.NewRegister(&registry.EtcdConfig{
		RegistryDir: "/backend/services",
		TTL:         10 * time.Second,
		EtcdConfig: &etcd3.Config{
			Endpoints: []string{"127.0.0.1:2379"},
		},
	})
	if err != nil {
		log.Panic(err)
		return
	}
	service := &registry.Service{
		Name:     "test",
		AppId:    "test",
		Version:  "v1.0",
		Address:  fmt.Sprintf("127.0.0.1:%d", *port),
		Metadata: g.Map{"weight": 1},
	}
	s := server.NewGrpcServer(server.GrpcServerConfig{
		Addr: fmt.Sprintf("0.0.0.0:%d", *port),
	})
	proto.RegisterTestServer(s.Server, new(helloServer))
	s.Start()
	if err := register.Register(service); err != nil {
		panic(err)
	}
	s.Wait()

	//signalChan := make(chan os.Signal, 1)
	//signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	//<-signalChan
	//register.Unregister(service)
	//s.Stop()

}

// go run main.go -node node1 -port 28544
// go run main.go -node node2 -port 18562
// go run main.go -node node3 -port 27772
func main() {
	flag.Parse()
	StartService()
}
