module github.com/gogf/katyusha

go 1.11

require (
	github.com/gogf/gf v1.16.7-0.20210903025403-077a41911bac
	github.com/golang/protobuf v1.5.2
	go.etcd.io/etcd/api/v3 v3.5.1
	go.etcd.io/etcd/client/v3 v3.5.1
	go.opentelemetry.io/otel v1.2.0
	go.opentelemetry.io/otel/trace v1.2.0
	go.opentelemetry.io/otel/exporters/jaeger v1.2.0
    go.opentelemetry.io/otel/sdk v1.2.0
	golang.org/x/net v0.0.0-20211123203042-d83791d6bcd9
	golang.org/x/sys v0.0.0-20211124211545-fe61309f8881 // indirect
	google.golang.org/grpc v1.40.0
)

//replace (
//	go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.0.0-20201103155942-6e800b9b0161
//	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20201103155942-6e800b9b0161
//	google.golang.org/grpc => google.golang.org/grpc v1.29.1
//)
