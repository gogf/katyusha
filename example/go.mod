module github.com/gogf/katyusha/example

go 1.15

require (
	github.com/gogf/gf/contrib/registry/etcd/v2 v2.1.0-rc3.0.20220520082600-c90acf81d6a8
	github.com/gogf/gf/v2 v2.1.0-rc3.0.20220520082600-c90acf81d6a8
	github.com/gogf/katyusha v0.3.0
	github.com/golang/protobuf v1.5.2
	golang.org/x/net v0.0.0-20220520000938-2e3eb7b945c2
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.0
)

replace github.com/gogf/katyusha => ../
