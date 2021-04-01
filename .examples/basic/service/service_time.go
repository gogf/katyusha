// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package service

import (
	"fmt"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/katyusha/.examples/basic/protocol"
	"golang.org/x/net/context"
)

type Time struct{}

func (s *Time) Now(ctx context.Context, r *protocol.NowReq) (*protocol.NowRes, error) {
	text := fmt.Sprintf(`%s: %s`, gcmd.GetOpt("node"), gtime.Now().String())
	return &protocol.NowRes{Time: text}, nil
}
