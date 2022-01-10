// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/katyusha/.examples/basic/protobuf"
)

// go run client_echo.go
func main() {
	var (
		ctx         = gctx.New()
		client, err = protobuf.NewClient()
	)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}

	res, err := client.Echo().Say(ctx, &protobuf.SayReq{Content: "Hello"})
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	g.Log().Print(ctx, "Response:", res.Content)
}
