// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/katyusha.

package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/katyusha/.examples/basic/protobuf"
	"golang.org/x/net/context"
	"time"
)

// go run client_echo.go
func main() {
	client, err := protobuf.NewClient()
	if err != nil {
		g.Log().Fatal(err)
	}
	for i := 0; i < 500; i++ {
		res, err := client.EchoClient.Say(context.Background(), &protobuf.SayReq{Content: "Hello"})
		if err != nil {
			g.Log().Error(err)
			time.Sleep(time.Second)
			continue
		}
		time.Sleep(time.Second)
		g.Log().Print("Response:", res.Content)
	}
}
