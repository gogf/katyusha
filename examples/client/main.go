package main

import (
	"github.com/gogf/katyusha/balancer"
	"github.com/gogf/katyusha/examples/proto"
	"github.com/gogf/katyusha/registry"
	etcd3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

func main() {
	etcdConfig := &registry.EtcdConfig{
		RegistryDir:  "/backend/services",
		KeepaliveTtl: 10 * time.Second,
		EtcdConfig: &etcd3.Config{
			Endpoints: []string{"127.0.0.1:2379"},
		},
	}
	service := &registry.Service{
		Name:    "test",
		AppId:   "test",
		Version: "v1.0",
	}
	resolver.Register(&registry.EtcdResolver{
		EtcdScheme: "etcd3",
		EtcdConfig: etcdConfig,
		Service:    service,
	})

	c, err := grpc.Dial("etcd3:///", grpc.WithInsecure(), grpc.WithBalancerName(balancer.RoundRobin))
	if err != nil {
		log.Printf("grpc dial: %s", err)
		return
	}
	defer c.Close()
	client := proto.NewTestClient(c)

	for i := 0; i < 500; i++ {
		resp, err := client.Say(context.Background(), &proto.SayReq{Content: "round robin"})
		if err != nil {
			log.Println("aa:", err)
			time.Sleep(time.Second)
			continue
		}
		time.Sleep(time.Second)
		log.Printf(resp.Content)
	}
}
