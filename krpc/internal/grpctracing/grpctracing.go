// opentelemetry-go-contrib/instrumentation/google.golang.org/grpc/otelgrpc/interceptor.go

package grpctracing

import (
	"context"
	"google.golang.org/grpc/metadata"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

const (
	// GRPCStatusCodeKey is convention for numeric status code of a gRPC request.
	GRPCStatusCodeKey = label.Key("rpc.grpc.status_code")
)

const (
	tracingMaxContentLogSize         = 512 * 1024 // Max log size for request and response body.
	tracingInstrumentGrpcClient      = "github.com/gogf/katyusha/krpc.GrpcClient"
	tracingInstrumentGrpcServer      = "github.com/gogf/katyusha/krpc.GrpcServer"
	tracingEventGrpcRequest          = "grpc.request"
	tracingEventGrpcRequestMessage   = "grpc.request.message"
	tracingEventGrpcRequestBaggage   = "grpc.request.baggage"
	tracingEventGrpcMetadataOutgoing = "grpc.metadata.outgoing"
	tracingEventGrpcMetadataIncoming = "grpc.metadata.incoming"
	tracingEventGrpcResponse         = "grpc.response"
	tracingEventGrpcResponseMessage  = "grpc.response.message"
)

type metadataSupplier struct {
	metadata *metadata.MD
}

func (s *metadataSupplier) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (s *metadataSupplier) Set(key string, value string) {
	s.metadata.Set(key, value)
}

// Inject injects correlation context and span context into the gRPC
// metadata object. This function is meant to be used on outgoing
// requests.
func Inject(ctx context.Context, metadata *metadata.MD) {
	otel.GetTextMapPropagator().Inject(ctx, &metadataSupplier{
		metadata: metadata,
	})
}

// Extract returns the correlation context and span context that
// another service encoded in the gRPC metadata object with Inject.
// This function is meant to be used on incoming requests.
func Extract(ctx context.Context, metadata *metadata.MD) ([]label.KeyValue, trace.SpanContext) {
	ctx = otel.GetTextMapPropagator().Extract(ctx, &metadataSupplier{
		metadata: metadata,
	})
	labelSet := baggage.Set(ctx)
	return (&labelSet).ToSlice(), trace.RemoteSpanContextFromContext(ctx)
}
