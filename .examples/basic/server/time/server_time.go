// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/katyusha/.examples/basic/protocol"
	"github.com/gogf/katyusha/.examples/basic/service"
	"github.com/gogf/katyusha/krpc"
)

// go run server_time.go -node node1 -port 8100
// go run server_time.go -node node2 -port 8101
// go run server_time.go -node node3 -port 8102
func main() {
	c := krpc.Server.NewGrpcServerConfig()
	c.MustSetWithMap(g.Map{
		"AppId":            "time",
		"Address":          fmt.Sprintf("0.0.0.0:%s", gcmd.GetOpt("port")),
		"AccessLogEnabled": true,
	})
	s := krpc.Server.NewGrpcServer(c)
	protocol.RegisterTimeServer(s.Server, new(service.Time))
	s.Run()
}
