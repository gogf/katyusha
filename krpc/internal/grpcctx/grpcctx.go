package grpcctx

import (
	"context"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"google.golang.org/grpc/metadata"
)

type (
	grpcCtx struct{}
)

var (
	Ctx = &grpcCtx{}
)

func (c *grpcCtx) NewIncoming(ctx context.Context, data ...g.Map) context.Context {
	if len(data) > 0 {
		incomingMd := make(metadata.MD)
		for key, value := range data[0] {
			incomingMd.Set(key, gconv.String(value))
		}
		return metadata.NewIncomingContext(ctx, incomingMd)
	}
	return metadata.NewIncomingContext(ctx, nil)
}

func (c *grpcCtx) NewOutgoing(ctx context.Context, data ...g.Map) context.Context {
	if len(data) > 0 {
		outgoingMd := make(metadata.MD)
		for key, value := range data[0] {
			outgoingMd.Set(key, gconv.String(value))
		}
		return metadata.NewOutgoingContext(ctx, outgoingMd)
	}
	return metadata.NewOutgoingContext(ctx, nil)
}

func (c *grpcCtx) IncomingToOutgoing(ctx context.Context, keys ...string) context.Context {
	incomingMd, _ := metadata.FromIncomingContext(ctx)
	if incomingMd == nil {
		return ctx
	}
	outgoingMd, _ := metadata.FromOutgoingContext(ctx)
	if outgoingMd == nil {
		outgoingMd = make(metadata.MD)
	}
	if len(keys) > 0 {
		for _, key := range keys {
			outgoingMd[key] = append(outgoingMd[key], incomingMd.Get(key)...)
		}
	} else {
		for key, values := range incomingMd {
			outgoingMd[key] = append(outgoingMd[key], values...)
		}
	}
	return metadata.NewOutgoingContext(ctx, outgoingMd)
}

func (c *grpcCtx) IncomingMap(ctx context.Context) *gmap.Map {
	var (
		data          = gmap.New()
		incomingMd, _ = metadata.FromIncomingContext(ctx)
	)
	for key, values := range incomingMd {
		if len(values) == 1 {
			data.Set(key, values[0])
		} else {
			data.Set(key, values)
		}
	}
	return data
}

func (c *grpcCtx) OutgoingMap(ctx context.Context) *gmap.Map {
	var (
		data          = gmap.New()
		outgoingMd, _ = metadata.FromOutgoingContext(ctx)
	)
	for key, values := range outgoingMd {
		if len(values) == 1 {
			data.Set(key, values[0])
		} else {
			data.Set(key, values)
		}
	}
	return data
}

func (c *grpcCtx) SetIncoming(ctx context.Context, data g.Map) context.Context {
	incomingMd, _ := metadata.FromIncomingContext(ctx)
	if incomingMd == nil {
		incomingMd = make(metadata.MD)
	}
	for key, value := range data {
		incomingMd.Set(key, gconv.String(value))
	}
	return metadata.NewIncomingContext(ctx, incomingMd)
}

func (c *grpcCtx) SetOutgoing(ctx context.Context, data g.Map) context.Context {
	outgoingMd, _ := metadata.FromOutgoingContext(ctx)
	if outgoingMd == nil {
		outgoingMd = make(metadata.MD)
	}
	for key, value := range data {
		outgoingMd.Set(key, gconv.String(value))
	}
	return metadata.NewOutgoingContext(ctx, outgoingMd)
}
