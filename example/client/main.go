package main

import (
	etcd3 "github.com/coreos/etcd/clientv3"
	"github.com/gogf/katyusha/balancer"
	"github.com/gogf/katyusha/example/proto"
	"github.com/gogf/katyusha/registry"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	etcdConfig := etcd3.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
	}
	registry.RegisterResolver("etcd3", etcdConfig, "/backend/services", "test", "1.0")

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
