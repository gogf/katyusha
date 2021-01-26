package tracing

import (
	"context"
	"github.com/gogf/katyusha"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryClientInterceptor returns a grpc.UnaryClientInterceptor suitable
// for use in a grpc.Dial call.
func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, callOpts ...grpc.CallOption) error {
	requestMetadata, _ := metadata.FromOutgoingContext(ctx)
	metadataCopy := requestMetadata.Copy()
	tracer := newConfig(nil).TracerProvider.Tracer(
		"github.com/gogf/katyusha/krpc.GrpcClient",
		trace.WithInstrumentationVersion(katyusha.VERSION),
	)
	name, attr := spanInfo(method, cc.Target())
	var span trace.Span
	ctx, span = tracer.Start(
		ctx,
		name,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attr...),
	)
	defer span.End()

	Inject(ctx, &metadataCopy)

	ctx = metadata.NewOutgoingContext(ctx, metadataCopy)

	messageSent.Event(ctx, 1, req)

	err := invoker(ctx, method, req, reply, cc, callOpts...)

	messageReceived.Event(ctx, 1, reply)
	if err != nil {
		s, _ := status.FromError(err)
		span.SetStatus(codes.Error, s.Message())
		span.SetAttributes(statusCodeAttr(s.Code()))
	} else {
		span.SetAttributes(statusCodeAttr(grpcCodes.OK))
	}
	return err
}

// StreamClientInterceptor returns a grpc.StreamClientInterceptor suitable
// for use in a grpc.Dial call.
func StreamClientInterceptor(
	ctx context.Context, desc *grpc.StreamDesc,
	cc *grpc.ClientConn, method string, streamer grpc.Streamer,
	callOpts ...grpc.CallOption) (grpc.ClientStream, error) {
	requestMetadata, _ := metadata.FromOutgoingContext(ctx)
	metadataCopy := requestMetadata.Copy()

	tracer := newConfig(nil).TracerProvider.Tracer(
		"github.com/gogf/katyusha/krpc.GrpcClient",
		trace.WithInstrumentationVersion(katyusha.VERSION),
	)

	name, attr := spanInfo(method, cc.Target())
	var span trace.Span
	ctx, span = tracer.Start(
		ctx,
		name,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attr...),
	)

	Inject(ctx, &metadataCopy)
	ctx = metadata.NewOutgoingContext(ctx, metadataCopy)

	s, err := streamer(ctx, desc, cc, method, callOpts...)
	stream := wrapClientStream(s, desc)

	go func() {
		if err == nil {
			err = <-stream.finished
		}

		if err != nil {
			s, _ := status.FromError(err)
			span.SetStatus(codes.Error, s.Message())
			span.SetAttributes(statusCodeAttr(s.Code()))
		} else {
			span.SetAttributes(statusCodeAttr(grpcCodes.OK))
		}

		span.End()
	}()

	return stream, err
}
