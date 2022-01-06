module github.com/gogf/example

go 1.16

require (
	github.com/gogf/gf/v2 v2.0.0-rc.0.20220104132444-99455e328b35
	github.com/gogf/katyusha v0.2.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	go.opentelemetry.io/otel v1.2.0
	go.opentelemetry.io/otel/exporters/jaeger v1.2.0
	go.opentelemetry.io/otel/sdk v1.2.0
	golang.org/x/net v0.0.0-20211123203042-d83791d6bcd9
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/gogf/katyusha => ../
