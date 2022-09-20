// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package service

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"golang.org/x/net/context"

	"github.com/gogf/katyusha/example/basic/protobuf"
)

// Echo is the service for echo.
type Echo struct{}

// Say implements the protobuf.EchoServer interface.
func (s *Echo) Say(ctx context.Context, r *protobuf.SayReq) (*protobuf.SayRes, error) {
	g.Log().Print(ctx, "Received:", r.Content)
	text := fmt.Sprintf(`%s: > %s`, gcmd.GetOpt("node", "default"), r.Content)
	return &protobuf.SayRes{Content: text}, nil
}
