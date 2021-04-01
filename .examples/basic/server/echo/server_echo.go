// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package main

import (
	"github.com/gogf/katyusha/.examples/basic/protocol"
	"github.com/gogf/katyusha/.examples/basic/service"
	"github.com/gogf/katyusha/krpc"
)

// go run server_echo.go -node node1
// go run server_echo.go -node node2
// go run server_echo.go -node node3
func main() {
	s := krpc.Server.NewGrpcServer()
	protocol.RegisterEchoServer(s.Server, new(service.Echo))
	s.Run()
}
