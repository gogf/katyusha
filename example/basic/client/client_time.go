// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package main

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/katyusha/example/basic/protobuf"
)

func main() {
	var (
		ctx         = gctx.New()
		client, err = protobuf.NewClient()
	)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}

	for i := 0; i < 500; i++ {
		res, err := client.Time().Now(ctx, &protobuf.NowReq{})
		if err != nil {
			g.Log().Error(ctx, err)
			time.Sleep(time.Second)
			continue
		}
		time.Sleep(time.Second)
		g.Log().Print(ctx, "Time:", res.Time)
	}
}
