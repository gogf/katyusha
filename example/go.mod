module github.com/gogf/katyusha/example

go 1.15

require (
	github.com/gogf/gf/contrib/registry/etcd/v2 v2.1.4
	github.com/gogf/gf/v2 v2.1.4
	github.com/gogf/katyusha v0.4.0
	github.com/golang/protobuf v1.5.2
	golang.org/x/net v0.0.0-20220919232410-f2f64ebce3c1
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.1
)

replace github.com/gogf/katyusha => ../
