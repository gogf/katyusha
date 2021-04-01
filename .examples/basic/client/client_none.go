// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/katyusha/.examples/basic/protocol"
	"github.com/gogf/katyusha/krpc"
	"golang.org/x/net/context"
	"time"
)

func main() {
	conn, err := krpc.Client.NewGrpcClientConn("none")
	if err != nil {
		g.Log().Fatal(err)
	}
	defer conn.Close()
	client := protocol.NewTimeClient(conn)

	for i := 0; i < 500; i++ {
		res, err := client.Now(context.Background(), &protocol.NowReq{})
		if err != nil {
			g.Log().Error(err)
			time.Sleep(time.Second)
			continue
		}
		time.Sleep(time.Second)
		g.Log().Print("Time:", res.Time)
	}
}
