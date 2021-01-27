package krpc_test

import (
	"context"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/katyusha/krpc"
	"google.golang.org/grpc/metadata"
	"testing"
)

func Test_Ctx_Basic(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"k1", "v1",
		"k2", "v2",
	))
	gtest.C(t, func(t *gtest.T) {
		m1 := krpc.Ctx.IncomingMap(ctx)
		t.Assert(m1.Get("k1"), "v1")
		t.Assert(m1.Get("k2"), "v2")
		m2 := krpc.Ctx.OutgoingMap(ctx)
		t.Assert(m2.Size(), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		ctx := krpc.Ctx.IncomingToOutgoing(ctx)
		m1 := krpc.Ctx.IncomingMap(ctx)
		t.Assert(m1.Get("k1"), "v1")
		t.Assert(m1.Get("k2"), "v2")
		m2 := krpc.Ctx.OutgoingMap(ctx)
		t.Assert(m2.Get("k1"), "v1")
		t.Assert(m2.Get("k2"), "v2")
	})
	gtest.C(t, func(t *gtest.T) {
		ctx := krpc.Ctx.IncomingToOutgoing(ctx, "k1")
		m1 := krpc.Ctx.IncomingMap(ctx)
		t.Assert(m1.Get("k1"), "v1")
		t.Assert(m1.Get("k2"), "v2")
		m2 := krpc.Ctx.OutgoingMap(ctx)
		t.Assert(m2.Get("k1"), "v1")
		t.Assert(m2.Get("k2"), "")
	})
	gtest.C(t, func(t *gtest.T) {
		ctx := krpc.Ctx.NewIncoming(ctx)
		ctx = krpc.Ctx.SetIncoming(ctx, g.Map{"k1": "v1"})
		ctx = krpc.Ctx.SetIncoming(ctx, g.Map{"k2": "v2"})
		ctx = krpc.Ctx.SetOutgoing(ctx, g.Map{"k3": "v3"})
		ctx = krpc.Ctx.SetOutgoing(ctx, g.Map{"k4": "v4"})
		m1 := krpc.Ctx.IncomingMap(ctx)
		t.Assert(m1.Get("k1"), "v1")
		t.Assert(m1.Get("k2"), "v2")
		m2 := krpc.Ctx.OutgoingMap(ctx)
		t.Assert(m2.Get("k3"), "v3")
		t.Assert(m2.Get("k4"), "v4")
	})
}
