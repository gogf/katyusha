// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package krpc_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"google.golang.org/grpc/metadata"

	"github.com/gogf/katyusha/krpc"
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
