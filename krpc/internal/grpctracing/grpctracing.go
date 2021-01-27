// opentelemetry-go-contrib/instrumentation/google.golang.org/grpc/otelgrpc/interceptor.go

package grpctracing

import (
	"context"

	"google.golang.org/grpc/metadata"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	// GRPCStatusCodeKey is convention for numeric status code of a gRPC request.
	GRPCStatusCodeKey = label.Key("rpc.grpc.status_code")
)

const (
	tracingMaxContentLogSize         = 512 * 1024 // Max log size for request and response body.
	tracingEventGrpcRequest          = "grpc.request"
	tracingEventGrpcRequestMessage   = "grpc.request.message"
	tracingEventGrpcMetadataOutgoing = "grpc.metadata.outgoing"
	tracingEventGrpcMetadataIncoming = "grpc.metadata.incoming"
	tracingEventGrpcResponse         = "grpc.response"
	tracingEventGrpcResponseMessage  = "grpc.response.message"
)

// config is a group of options for this instrumentation.
type config struct {
	Propagators    propagation.TextMapPropagator
	TracerProvider trace.TracerProvider
}

// Option applies an option value for a config.
type Option interface {
	Apply(*config)
}

// newConfig returns a config configured with all the passed Options.
func newConfig(opts []Option) *config {
	c := &config{
		Propagators:    otel.GetTextMapPropagator(),
		TracerProvider: otel.GetTracerProvider(),
	}
	c.Propagators = propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	for _, o := range opts {
		o.Apply(c)
	}
	return c
}

type propagatorsOption struct{ p propagation.TextMapPropagator }

func (o propagatorsOption) Apply(c *config) {
	c.Propagators = o.p
}

// WithPropagators returns an Option to use the Propagators when extracting
// and injecting trace context from requests.
func WithPropagators(p propagation.TextMapPropagator) Option {
	return propagatorsOption{p: p}
}

type tracerProviderOption struct{ tp trace.TracerProvider }

func (o tracerProviderOption) Apply(c *config) {
	c.TracerProvider = o.tp
}

// WithTracerProvider returns an Option to use the TracerProvider when
// creating a Tracer.
func WithTracerProvider(tp trace.TracerProvider) Option {
	return tracerProviderOption{tp: tp}
}

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
func Inject(ctx context.Context, metadata *metadata.MD, opts ...Option) {
	c := newConfig(opts)
	c.Propagators.Inject(ctx, &metadataSupplier{
		metadata: metadata,
	})
}

// Extract returns the correlation context and span context that
// another service encoded in the gRPC metadata object with Inject.
// This function is meant to be used on incoming requests.
func Extract(ctx context.Context, metadata *metadata.MD, opts ...Option) ([]label.KeyValue, trace.SpanContext) {
	c := newConfig(opts)
	ctx = c.Propagators.Extract(ctx, &metadataSupplier{
		metadata: metadata,
	})
	labelSet := baggage.Set(ctx)
	return (&labelSet).ToSlice(), trace.RemoteSpanContextFromContext(ctx)
}
