// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package service

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/katyusha/.examples/basic/protobuf"
	"golang.org/x/net/context"
)

type Echo struct{}

func (s *Echo) Say(ctx context.Context, r *protobuf.SayReq) (*protobuf.SayRes, error) {
	g.Log().Println("Received:", r.Content)
	text := fmt.Sprintf(`%s: > %s`, gcmd.GetOpt("node", "default"), r.Content)
	return &protobuf.SayRes{Content: text}, nil
}
