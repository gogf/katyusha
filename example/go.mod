module github.com/gogf/katyusha/example

go 1.15

require (
	github.com/gogf/gf/contrib/registry/etcd/v2 v2.1.3
	github.com/gogf/gf/v2 v2.1.3
	github.com/gogf/katyusha v0.4.0
	github.com/golang/protobuf v1.5.2
	golang.org/x/net v0.0.0-20220822230855-b0a4917ee28c
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.1
)

replace github.com/gogf/katyusha => ../
