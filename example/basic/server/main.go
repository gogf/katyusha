// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package main

import (
	"github.com/gogf/katyusha/example/basic/protobuf"
	"github.com/gogf/katyusha/example/basic/service"
	"github.com/gogf/katyusha/krpc"
)

func main() {
	s := krpc.Server.NewGrpcServer()
	protobuf.RegisterEchoServer(s.Server, new(service.Echo))
	protobuf.RegisterTimeServer(s.Server, new(service.Time))
	s.Run()
}
