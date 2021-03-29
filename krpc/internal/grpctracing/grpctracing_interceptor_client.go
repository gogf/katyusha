package grpctracing

import (
	"context"
	"github.com/gogf/gf/net/gtrace"
	"github.com/gogf/katyusha"
	"github.com/gogf/katyusha/krpc/internal/grpcctx"
	"github.com/gogf/katyusha/krpc/internal/grpcutils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/attribute"
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
	tracer := otel.GetTracerProvider().Tracer(
		tracingInstrumentGrpcClient,
		trace.WithInstrumentationVersion(katyusha.VERSION),
	)
	requestMetadata, _ := metadata.FromOutgoingContext(ctx)
	metadataCopy := requestMetadata.Copy()
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

	span.SetAttributes(gtrace.CommonLabels()...)

	span.AddEvent(tracingEventGrpcRequest, trace.WithAttributes(
		attribute.Any(tracingEventGrpcRequestBaggage, gtrace.GetBaggageMap(ctx)),
		attribute.Any(tracingEventGrpcMetadataOutgoing, grpcctx.Ctx.OutgoingMap(ctx)),
		attribute.String(
			tracingEventGrpcRequestMessage,
			grpcutils.MarshalMessageToJsonStringForTracing(
				req, "Request", tracingMaxContentLogSize,
			),
		),
	))

	err := invoker(ctx, method, req, reply, cc, callOpts...)

	span.AddEvent(tracingEventGrpcResponse, trace.WithAttributes(
		attribute.String(
			tracingEventGrpcResponseMessage,
			grpcutils.MarshalMessageToJsonStringForTracing(
				reply, "Response", tracingMaxContentLogSize,
			),
		),
	))

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
	tracer := otel.GetTracerProvider().Tracer(
		tracingInstrumentGrpcClient,
		trace.WithInstrumentationVersion(katyusha.VERSION),
	)
	requestMetadata, _ := metadata.FromOutgoingContext(ctx)
	metadataCopy := requestMetadata.Copy()
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

	span.SetAttributes(gtrace.CommonLabels()...)

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
