// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package main

import (
	"context"
	"strings"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gtrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"

	"github.com/gogf/katyusha/.examples/tracing/protobuf/user"
	"github.com/gogf/katyusha/krpc"
)

const (
	JaegerEndpoint = "http://localhost:14268/api/traces"
	ServiceName    = "tracing-grpc-client"
)

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer(serviceName, endpoint string) (tp *trace.TracerProvider, err error) {
	var endpointOption jaeger.EndpointOption
	if strings.HasPrefix(endpoint, "http") {
		// HTTP.
		endpointOption = jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint))
	} else {
		// UDP.
		endpointOption = jaeger.WithAgentEndpoint(jaeger.WithAgentHost(endpoint))
	}

	// Create the Jaeger exporter
	exp, err := jaeger.New(endpointOption)
	if err != nil {
		return nil, err
	}
	tp = trace.NewTracerProvider(
		// Always be sure to batch in production.
		trace.WithBatcher(exp),
		// Record information about this application in an Resource.
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func StartRequests() {
	ctx, span := gtrace.NewSpan(context.Background(), "StartRequests")
	defer span.End()

	grpcClientOptions := make([]grpc.DialOption, 0)
	grpcClientOptions = append(
		grpcClientOptions,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithChainUnaryInterceptor(
			krpc.Client.UnaryError,
			krpc.Client.UnaryTracing,
		),
	)

	conn, err := grpc.Dial(":8000", grpcClientOptions...)
	if err != nil {
		g.Log().Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := user.NewUserClient(conn)

	// Insert.
	insertRes, err := client.Insert(ctx, &user.InsertReq{
		Name: "john",
	})
	if err != nil {
		g.Log().Ctx(ctx).Fatalf(`%+v`, err)
	}
	g.Log().Ctx(ctx).Println("insert:", insertRes.Id)

	// Query.
	queryRes, err := client.Query(ctx, &user.QueryReq{
		Id: insertRes.Id,
	})
	if err != nil {
		g.Log().Ctx(ctx).Printf(`%+v`, err)
		return
	}
	g.Log().Ctx(ctx).Println("query:", queryRes)

	// Delete.
	_, err = client.Delete(ctx, &user.DeleteReq{
		Id: insertRes.Id,
	})
	if err != nil {
		g.Log().Ctx(ctx).Printf(`%+v`, err)
		return
	}
	g.Log().Ctx(ctx).Println("delete:", insertRes.Id)

	// Delete with error.
	_, err = client.Delete(ctx, &user.DeleteReq{
		Id: -1,
	})
	if err != nil {
		g.Log().Ctx(ctx).Printf(`%+v`, err)
		return
	}
	g.Log().Ctx(ctx).Println("delete:", -1)

}

func main() {
	_, err := initTracer(ServiceName, JaegerEndpoint)
	if err != nil {
		g.Log().Fatal(err)
	}

	StartRequests()
}
