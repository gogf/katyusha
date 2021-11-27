// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package service

import (
	"fmt"

	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/katyusha/.examples/basic/protobuf"
	"golang.org/x/net/context"
)

type Time struct{}

func (s *Time) Now(ctx context.Context, r *protobuf.NowReq) (*protobuf.NowRes, error) {
	text := fmt.Sprintf(`%s: %s`, gcmd.GetOpt("node", "default"), gtime.Now().String())
	return &protobuf.NowRes{Time: text}, nil
}
