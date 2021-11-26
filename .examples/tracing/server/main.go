// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/gogf/gcache-adapter/adapter"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/katyusha/.examples/tracing/protobuf/user"
	"github.com/gogf/katyusha/krpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

type server struct{}

const (
	JaegerEndpoint = "http://localhost:14268/api/traces"
	ServiceName    = "tracing-grpc-server"
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

// Insert is a route handler for inserting user info into dtabase.
func (s *server) Insert(ctx context.Context, req *user.InsertReq) (*user.InsertRes, error) {
	res := user.InsertRes{}
	result, err := g.Table("user").Ctx(ctx).Insert(g.Map{
		"name": req.Name,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	res.Id = int32(id)
	return &res, nil
}

// Query is a route handler for querying user info. It firstly retrieves the info from redis,
// if there's nothing in the redis, it then does db select.
func (s *server) Query(ctx context.Context, req *user.QueryReq) (*user.QueryRes, error) {
	res := user.QueryRes{}
	err := g.Table("user").
		Ctx(ctx).
		Cache(5*time.Second, s.userCacheKey(req.Id)).
		WherePri(req.Id).
		Scan(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// Delete is a route handler for deleting specified user info.
func (s *server) Delete(ctx context.Context, req *user.DeleteReq) (*user.DeleteRes, error) {
	res := user.DeleteRes{}
	_, err := g.Model("user").
		Ctx(ctx).
		Cache(-1, s.userCacheKey(req.Id)).
		WherePri(req.Id).
		Delete()
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *server) userCacheKey(id int32) string {
	return fmt.Sprintf(`userInfo:%d`, id)
}

func main() {
	_, err := initTracer(ServiceName, JaegerEndpoint)
	if err != nil {
		g.Log().Fatal(err)
	}

	g.DB().GetCache().SetAdapter(adapter.NewRedis(g.Redis()))

	address := ":8000"
	listen, err := net.Listen("tcp", address)
	if err != nil {
		g.Log().Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			krpc.Server.UnaryError,
			krpc.Server.UnaryRecover,
			krpc.Server.UnaryTracing,
			krpc.Server.UnaryValidate,
		),
	)
	user.RegisterUserServer(s, &server{})
	g.Log().Printf("grpc server starts listening on %s", address)
	if err := s.Serve(listen); err != nil {
		g.Log().Fatalf("failed to serve: %v", err)
	}
}
