module github.com/gogf/katyusha

go 1.11

require (
	github.com/gogf/gcache-adapter v0.0.4-0.20210126062229-c84b9cefa528
	github.com/gogf/gf v1.15.2-0.20210126085339-24e2c7926e39
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.3
	github.com/json-iterator/go v1.1.10 // indirect
	go.etcd.io/etcd/api/v3 v3.5.0-pre
	go.etcd.io/etcd/client/v3 v3.0.0-20201118182908-c11ddc65cea1
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.16.0
	go.opentelemetry.io/otel/sdk v0.16.0
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/text v0.3.4 // indirect
	google.golang.org/grpc v1.33.2
	google.golang.org/protobuf v1.25.0
)

replace (
	go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.0.0-20201103155942-6e800b9b0161
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20201103155942-6e800b9b0161
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)
