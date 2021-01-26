package tracing

import (
	"context"
	"github.com/gogf/katyusha"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a grpc.UnaryServerInterceptor suitable
// for use in a grpc.NewServer call.
func UnaryServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	requestMetadata, _ := metadata.FromIncomingContext(ctx)
	metadataCopy := requestMetadata.Copy()

	entries, spanCtx := Extract(ctx, &metadataCopy)
	ctx = baggage.ContextWithValues(ctx, entries...)

	tracer := newConfig(nil).TracerProvider.Tracer(
		"github.com/gogf/katyusha/krpc.GrpcServer",
		trace.WithInstrumentationVersion(katyusha.VERSION),
	)

	name, attr := spanInfo(info.FullMethod, peerFromCtx(ctx))
	ctx, span := tracer.Start(
		trace.ContextWithRemoteSpanContext(ctx, spanCtx),
		name,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attr...),
	)
	defer span.End()

	messageReceived.Event(ctx, 1, req)

	resp, err := handler(ctx, req)
	if err != nil {
		s, _ := status.FromError(err)
		span.SetStatus(codes.Error, s.Message())
		span.SetAttributes(statusCodeAttr(s.Code()))
		messageSent.Event(ctx, 1, s.Proto())
	} else {
		span.SetAttributes(statusCodeAttr(grpcCodes.OK))
		messageSent.Event(ctx, 1, resp)
	}

	return resp, err
}

// StreamServerInterceptor returns a grpc.StreamServerInterceptor suitable
// for use in a grpc.NewServer call.
func StreamServerInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	ctx := ss.Context()

	requestMetadata, _ := metadata.FromIncomingContext(ctx)
	metadataCopy := requestMetadata.Copy()

	entries, spanCtx := Extract(ctx, &metadataCopy)
	ctx = baggage.ContextWithValues(ctx, entries...)

	tracer := newConfig(nil).TracerProvider.Tracer(
		"github.com/gogf/katyusha/krpc.GrpcServer",
		trace.WithInstrumentationVersion(katyusha.VERSION),
	)

	name, attr := spanInfo(info.FullMethod, peerFromCtx(ctx))
	ctx, span := tracer.Start(
		trace.ContextWithRemoteSpanContext(ctx, spanCtx),
		name,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attr...),
	)
	defer span.End()

	err := handler(srv, wrapServerStream(ctx, ss))

	if err != nil {
		s, _ := status.FromError(err)
		span.SetStatus(codes.Error, s.Message())
		span.SetAttributes(statusCodeAttr(s.Code()))
	} else {
		span.SetAttributes(statusCodeAttr(grpcCodes.OK))
	}

	return err
}
